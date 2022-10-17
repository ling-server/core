package orm

import "github.com/beego/beego/orm"

var (
	// Crt is a global instance of ORM creator
	Crt = NewCreator()
)

// NewCreator creates an ORM creator
func NewCreator() Creator {
	return &creator{}
}

// Creator creates ORMer
// Introducing the "Creator" interface to eliminate the dependency on database
type Creator interface {
	Create() orm.Ormer
}

type creator struct{}

func (c *creator) Create() orm.Ormer {
	return orm.NewOrm()
}
