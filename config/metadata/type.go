package metadata

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Type - Use this interface to define and encapsulate the behavior of validation and transformation
type Type interface {
	// validate the configure value
	validate(str string) error
	// get the real type of current value, if it is int, return int, if it is string return string etc.
	get(str string) (interface{}, error)
}

// StringType ...
type StringType struct {
}

func (t *StringType) validate(str string) error {
	return nil
}

func (t *StringType) get(str string) (interface{}, error) {
	return str, nil
}

// NonEmptyStringType ...
type NonEmptyStringType struct {
	StringType
}

func (t *NonEmptyStringType) validate(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return ErrorStringValueIsEmpty
	}
	return nil
}

// IntType ..
type IntType struct {
}

func (t *IntType) validate(str string) error {
	_, err := parseInt(str)
	return err
}

func (t *IntType) get(str string) (interface{}, error) {
	return parseInt(str)
}

// PortType ...
type PortType struct {
	IntType
}

func (t *PortType) validate(str string) error {
	val, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	if val < 0 {
		return fmt.Errorf("network port should be greater than 0")
	}

	if val > 65535 {
		return fmt.Errorf("network port should be less than 65535")
	}

	return err
}

// Int64Type ...
type Int64Type struct {
}

func (t *Int64Type) validate(str string) error {
	_, err := parseInt64(str)
	return err
}

func (t *Int64Type) get(str string) (interface{}, error) {
	return parseInt64(str)
}

type Float64Type struct{}

func (f *Float64Type) validate(str string) error {
	_, err := parseFloat64(str)
	return err
}

func (f *Float64Type) get(str string) (interface{}, error) {
	return parseFloat64(str)
}

// BoolType ...
type BoolType struct {
}

func (t *BoolType) validate(str string) error {
	_, err := strconv.ParseBool(str)
	return err
}

func (t *BoolType) get(str string) (interface{}, error) {
	return strconv.ParseBool(str)
}

// PasswordType ...
type PasswordType struct {
}

func (t *PasswordType) validate(str string) error {
	return nil
}

func (t *PasswordType) get(str string) (interface{}, error) {
	return str, nil
}

// MapType ...
type MapType struct {
}

func (t *MapType) validate(str string) error {
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &result)
	return err
}

func (t *MapType) get(str string) (interface{}, error) {
	result := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &result)
	return result, err
}

// StringToStringMapType ...
type StringToStringMapType struct {
}

func (t *StringToStringMapType) validate(str string) error {
	result := map[string]string{}
	err := json.Unmarshal([]byte(str), &result)
	return err
}

func (t *StringToStringMapType) get(str string) (interface{}, error) {
	result := map[string]string{}
	err := json.Unmarshal([]byte(str), &result)
	return result, err
}

// QuotaType ...
type QuotaType struct {
	Int64Type
}

func (t *QuotaType) validate(str string) error {
	val, err := parseInt64(str)
	if err != nil {
		return err
	}

	if val <= 0 && val != -1 {
		return fmt.Errorf("quota value should be -1 or great than zero")
	}

	return nil
}

// parseInt64 returns int64 from string which support scientific notation
func parseInt64(str string) (int64, error) {
	val, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return val, nil
	}

	fval, err := strconv.ParseFloat(str, 64)
	if err == nil && fval == math.Trunc(fval) {
		return int64(fval), nil
	}

	return 0, fmt.Errorf("invalid int64 string: %s", str)
}

func parseInt(str string) (int, error) {
	val, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		return int(val), nil
	}

	fval, err := strconv.ParseFloat(str, 32)
	if err == nil && fval == math.Trunc(fval) {
		return int(fval), nil
	}

	return 0, fmt.Errorf("invalid int string: %s", str)
}

func parseFloat64(str string) (float64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return val, nil
	}

	return 0, fmt.Errorf("invalid float64 string: %s", str)
}
