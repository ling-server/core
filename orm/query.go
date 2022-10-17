package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/orm"
)

// KeyWords ...
type KeyWords = map[string]interface{}

// Query parameters
type Query struct {
	// Filter list
	Keywords KeyWords
	// Sort list
	Sorts []*Sort
	// Page number
	PageNumber int64
	// Page size
	PageSize int64
	// Deprecate, use "Sorts" instead
	Sorting string
}

// First make the query only fetch the first one record in the sorting order
func (q *Query) First(sorting ...*Sort) *Query {
	q.PageNumber = 1
	q.PageSize = 1
	if len(sorting) > 0 {
		q.Sorts = append(q.Sorts, sorting...)
	}

	return q
}

// New returns Query with keywords
func New(kw KeyWords) *Query {
	return &Query{Keywords: kw}
}

// MustClone returns the clone of query when it's not nil
// or returns a new Query instance
func MustClone(query *Query) *Query {
	q := &Query{
		Keywords: map[string]interface{}{},
	}
	if query != nil {
		q.PageNumber = query.PageNumber
		q.PageSize = query.PageSize
		q.Sorts = query.Sorts
		for k, v := range query.Keywords {
			q.Keywords[k] = v
		}
		for _, sort := range query.Sorts {
			q.Sorts = append(q.Sorts, &Sort{
				Key:  sort.Key,
				DESC: sort.DESC,
			})
		}
	}
	return q
}

// Sort specifies the order information
type Sort struct {
	Key  string
	DESC bool
}

// Range query
type Range struct {
	Min interface{}
	Max interface{}
}

// AndList query
type AndList struct {
	Values []interface{}
}

// OrList query
type OrList struct {
	Values []interface{}
}

// FuzzyMatchValue query
type FuzzyMatchValue struct {
	Value string
}

// NewSort creates new sort
func NewSort(key string, desc bool) *Sort {
	return &Sort{
		Key:  key,
		DESC: desc,
	}
}

// NewRange creates a new range
func NewRange(min, max interface{}) *Range {
	return &Range{
		Min: min,
		Max: max,
	}
}

// NewAndList creates a new and list
func NewAndList(values []interface{}) *AndList {
	return &AndList{
		Values: values,
	}
}

// NewOrList creates a new or list
func NewOrList(values []interface{}) *OrList {
	return &OrList{
		Values: values,
	}
}

// NewFuzzyMatchValue creates a new fuzzy match
func NewFuzzyMatchValue(value string) *FuzzyMatchValue {
	return &FuzzyMatchValue{
		Value: value,
	}
}

// QuerySetter generates the query setter according to the provided model and query.
// e.g.
//
//	type Foo struct{
//	  Field1 string `orm:"-"`                         // can not filter/sort
//	  Field2 string `orm:"column(customized_field2)"` // support filter by "Field2", "customized_field2"
//	  Field3 string `sort:"false"`                    // cannot be sorted
//	  Field4 string `sort:"default:desc"`             // the default field/order(asc/desc) to sort if no sorting specified in query.
//	  Field5 string `filter:"false"`                  // cannot be filtered
//	}
//
// // support filter by "Field6", "field6"
//
//	func (f *Foo) FilterByField6(ctx context.Context, qs orm.QuerySetter, key string, value interface{}) orm.QuerySetter {
//	  ...
//		 return qs
//	}
//
// Defining the method "GetDefaultSorts() []*q.Sort" for the model whose default sorting contains more than one fields
//
//	type Bar struct{
//	  Field1 string
//	  Field2 string
//	}
//
// // Sort by "Field1" desc, "Field2"
//
//	func (b *Bar) GetDefaultSorts() []*q.Sort {
//		return []*q.Sort{
//			{
//				Key:  "Field1",
//				DESC: true,
//			},
//			{
//				Key:  "Field2",
//				DESC: false,
//			},
//		 }
//	}
func QuerySetter(ctx context.Context, model interface{}, query *Query) (orm.QuerySeter, error) {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("<orm.QuerySetter> cannot use non-ptr model struct `%s`", getFullName(t.Elem()))
	}
	ormer, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}
	qs := ormer.QueryTable(model)
	if query == nil {
		return qs, nil
	}

	metadata := parseModel(model)
	// set filters
	qs = setFilters(ctx, qs, query, metadata)

	// sorting
	qs = setSorts(qs, query, metadata)

	// pagination
	if query.PageSize > 0 {
		qs = qs.Limit(query.PageSize)
		if query.PageNumber > 0 {
			qs = qs.Offset(query.PageSize * (query.PageNumber - 1))
		}
	}

	return qs, nil
}

// PaginationOnRawSQL append page information to the raw sql
// It should be called after the order by
// e.g.
// select a, b, c from mytable order by a limit ? offset ?
// it appends the " limit ? offset ? " to sql,
// and appends the limit value and offset value to the params of this query
func PaginationOnRawSQL(query *Query, sql string, params []interface{}) (string, []interface{}) {
	if query != nil && query.PageSize > 0 {
		sql += ` limit ?`
		params = append(params, query.PageSize)

		if query.PageNumber > 0 {
			sql += ` offset ?`
			params = append(params, (query.PageNumber-1)*query.PageSize)
		}
	}
	return sql, params
}

// QuerySetterForCount creates the query setter used for count with the sort and pagination information ignored
func QuerySetterForCount(ctx context.Context, model interface{}, query *Query, ignoredCols ...string) (orm.QuerySeter, error) {
	query = MustClone(query)
	query.Sorts = nil
	query.PageSize = 0
	query.PageNumber = 0
	return QuerySetter(ctx, model, query)
}

// set filters according to the query
func setFilters(ctx context.Context, qs orm.QuerySeter, query *Query, meta *metadata) orm.QuerySeter {
	for key, value := range query.Keywords {
		// The "strings.SplitN()" here is a workaround for the incorrect usage of query which should be avoided
		// e.g. use the query with the knowledge of underlying ORM implementation, the "OrList" should be used instead:
		// https://github.com/goharbor/harbor/blob/v2.2.0/src/controller/project/controller.go#L348
		k := strings.SplitN(key, orm.ExprSep, 2)[0]
		mk, filterable := meta.Filterable(k)
		if !filterable {
			// This is a workaround for the unsuitable usage of query, the keyword format for field and method should be consistent
			// e.g. "ArtifactDigest" or the snake case format "artifact_digest" should be used instead:
			// https://github.com/goharbor/harbor/blob/v2.2.0/src/controller/blob/controller.go#L233
			mk, filterable = meta.Filterable(snakeCase(k))
			if !filterable {
				continue
			}
		}
		// filter function defined, use it directly
		if mk.FilterFunc != nil {
			qs = mk.FilterFunc(ctx, qs, key, value)
			continue
		}
		// fuzzy match
		if f, ok := value.(*FuzzyMatchValue); ok {
			qs = qs.Filter(key+"__icontains", Escape(f.Value))
			continue
		}
		// range
		if r, ok := value.(*Range); ok {
			if r.Min != nil {
				qs = qs.Filter(key+"__gte", r.Min)
			}
			if r.Max != nil {
				qs = qs.Filter(key+"__lte", r.Max)
			}
			continue
		}
		// or list
		if ol, ok := value.(*OrList); ok {
			if ol == nil || len(ol.Values) == 0 {
				qs = qs.Filter(key+"__in", nil)
			} else {
				qs = qs.Filter(key+"__in", ol.Values...)
			}
			continue
		}
		// and list
		if _, ok := value.(*AndList); ok {
			// do nothing as and list needs to be handled by the logic of DAO
			continue
		}
		// exact match
		qs = qs.Filter(key, value)
	}
	return qs
}

// set sorts according to the query
func setSorts(qs orm.QuerySeter, query *Query, meta *metadata) orm.QuerySeter {
	var sortings []string
	for _, sort := range query.Sorts {
		if !meta.Sortable(sort.Key) {
			continue
		}
		sorting := sort.Key
		if sort.DESC {
			sorting = fmt.Sprintf("-%s", sorting)
		}
		sortings = append(sortings, sorting)
	}
	// if no sorts are specified, apply the default sort setting if exists
	if len(sortings) == 0 {
		for _, ds := range meta.DefaultSorts {
			sorting := ds.Key
			if ds.DESC {
				sorting = fmt.Sprintf("-%s", sorting)
			}
			sortings = append(sortings, sorting)
		}
	}
	if len(sortings) > 0 {
		qs = qs.OrderBy(sortings...)
	}
	return qs
}

// get reflect.Type name with package path.
func getFullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}