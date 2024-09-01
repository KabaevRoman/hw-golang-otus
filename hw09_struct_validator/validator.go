package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrMin    = errors.New("minimum value error")
	ErrMax    = errors.New("maximum value error")
	ErrIn     = errors.New("contains error")
	ErrLen    = errors.New("string length error")
	ErrRegexp = errors.New("regexp error")
	ErrClient = errors.New("client error")
)

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return v.Err.Error()
}

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, err := range v {
		if err.Err == nil {
			builder.WriteString("")
			continue
		}
		msg := fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error())
		builder.WriteString(msg)
	}
	return builder.String()
}

type Validator interface {
	IsValid(value reflect.Value) error
}

type MinValidator struct {
	Constraint int64
}

func (v *MinValidator) IsValid(value reflect.Value) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if v.Constraint > value.Index(i).Int() {
				return fmt.Errorf("%w minimum value is %d at index %d", ErrMin, v.Constraint, i)
			}
		}
	case reflect.Int:
		if v.Constraint > value.Int() {
			return fmt.Errorf("%w minimum value is %d", ErrMin, v.Constraint)
		}
	default:
		return fmt.Errorf("%w invalid type for minimum value: %T", ErrClient, value)
	}
	return nil
}

type RegexpValidator struct {
	Constraint *regexp.Regexp
}

func (v *RegexpValidator) IsValid(value reflect.Value) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			pattern := value.Index(i).String()
			if !v.Constraint.MatchString(pattern) {
				return fmt.Errorf(
					"%w string doesn't match provided regular expression at index %d",
					ErrRegexp,
					i,
				)
			}
		}
	case reflect.String:
		if !v.Constraint.MatchString(value.String()) {
			return fmt.Errorf("%w string doesn't match provided regular expression", ErrRegexp)
		}
	default:
		return fmt.Errorf("%w invalid type for regexp value: %T", ErrClient, value)
	}
	return nil
}

type LenValidator struct {
	Constraint int
}

func (v *LenValidator) IsValid(value reflect.Value) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			realLen := len([]rune(value.Index(i).String()))
			if realLen != v.Constraint {
				return fmt.Errorf(
					"%w invalid number of symbols in string expected: %d actual: %d at index %d",
					ErrLen,
					v.Constraint,
					realLen,
					i,
				)
			}
		}
	case reflect.String:
		realLen := len([]rune(value.String()))
		if realLen != v.Constraint {
			return fmt.Errorf(
				"%w invalid number of symbols in string expected: %d actual: %d",
				ErrLen,
				v.Constraint,
				realLen,
			)
		}
	default:
		return fmt.Errorf("%w invalid type for string len value: %T", ErrClient, value)
	}
	return nil
}

type MaxValidator struct {
	Constraint int64
}

// IsValid implements Validator.
func (v *MaxValidator) IsValid(value reflect.Value) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if v.Constraint < value.Index(i).Int() {
				return fmt.Errorf("%wminimum value is %d at index %d", ErrMax, v.Constraint, i)
			}
		}
	case reflect.Int:
		if v.Constraint < value.Int() {
			return fmt.Errorf("%wmaximum value is %d", ErrMax, v.Constraint)
		}
	default:
		return fmt.Errorf("%winvalid type for maximum value: %T", ErrClient, value)
	}
	return nil
}

type InValidator struct {
	Constraint map[any]struct{}
	Kind       reflect.Kind
	Keys       string
}

func (v *InValidator) IsValid(value reflect.Value) error {
	var ok bool
	switch v.Kind { //nolint:exhaustive
	case reflect.Int:
		_, ok = v.Constraint[value.Int()]
	case reflect.String:
		_, ok = v.Constraint[value.String()]
	default:
		return fmt.Errorf("%winvalid type for maximum value: %v", ErrIn, value.Kind())
	}
	if !ok {
		return fmt.Errorf("%wvalue not in specified: %s", ErrClient, v.Keys)
	}
	return nil
}

func getIn(values string, kind reflect.Kind) (Validator, error) {
	var err error
	mapping := make(map[any]struct{})
	switch kind { //nolint:exhaustive
	case reflect.Int:
		for _, val := range strings.Split(values, ",") {
			res, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("%w%s", ErrClient, err.Error())
			}
			mapping[int64(res)] = struct{}{}
		}
	case reflect.String:
		for _, val := range strings.Split(values, ",") {
			mapping[val] = struct{}{}
		}
	default:
		return nil, fmt.Errorf("%wunsupported kind %v", ErrClient, kind)
	}
	return &InValidator{Constraint: mapping, Keys: values, Kind: kind}, err
}

func getCondition(strCondition string, kind reflect.Kind) (validator Validator, err error) {
	split := strings.Split(strCondition, ":")
	if len(split) < 2 {
		return nil, fmt.Errorf("%w invalid value provided for tag", ErrClient)
	}
	switch split[0] {
	case "min":
		val, err := strconv.Atoi(split[1])
		return &MinValidator{Constraint: int64(val)}, err
	case "max":
		val, err := strconv.Atoi(split[1])
		return &MaxValidator{Constraint: int64(val)}, err
	case "in":
		return getIn(split[1], kind)
	case "regexp":
		regString := strings.ReplaceAll(split[1], `\\`, `\`)
		re, err := regexp.Compile(regString)
		return &RegexpValidator{Constraint: re}, err
	case "len":
		val, err := strconv.Atoi(split[1])
		return &LenValidator{Constraint: val}, err
	}
	return nil, fmt.Errorf("%w invalid value provided for tag", ErrClient)
}

func GetValidator(tag string, field reflect.StructField) (validators []Validator, err error) {
	kind := field.Type.Kind()
	if kind == reflect.Slice {
		kind = field.Type.Elem().Kind()
	}
	for _, strCond := range strings.Split(tag, "|") {
		validator, err := getCondition(strCond, kind)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validator)
	}
	return validators, err
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors
	validatingStruct := reflect.ValueOf(v)
	if validatingStruct.Kind() != reflect.Struct {
		return fmt.Errorf("%w value is not a struct", ErrClient)
	}
	for i := 0; i < validatingStruct.NumField(); i++ {
		typeField := validatingStruct.Type().Field(i)
		tag := typeField.Tag.Get("validate")
		if tag == "" {
			continue
		}
		validators, err := GetValidator(tag, typeField)
		if err != nil {
			return err
		}
		for _, validator := range validators {
			err := validator.IsValid(validatingStruct.Field(i))
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: typeField.Name,
					Err:   err,
				})
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
