package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents validation error in a single field.
type ValidationError struct {
	Field string
	Err   error
}

// ApplicationError represents error in validator itself.
type ApplicationError error

// ValidationErrors represents all validation errors of a struct.
type ValidationErrors []ValidationError

type tagValidator []func(v interface{}) error

type tagParser struct {
	lenRegexp   *regexp.Regexp
	reRegexp    *regexp.Regexp
	strInRegexp *regexp.Regexp
	minRegexp   *regexp.Regexp
	maxRegexp   *regexp.Regexp
	numInRegexp *regexp.Regexp
}

type (
	ErrorLower          error
	ErrorHigher         error
	ErrorNotInRange     error
	ErrorLen            error
	ErrorNotInList      error
	ErrorNotMatchRegexp error
)

var (
	errorLower          ErrorLower          = fmt.Errorf("value is lower than specified limit")
	errorHigher         ErrorHigher         = fmt.Errorf("value is higher than specified limit")
	errorNotInRange     ErrorNotInRange     = fmt.Errorf("value is not in list")
	errorLen            ErrorLen            = fmt.Errorf("value size does not match the length")
	errorNotInList      ErrorNotInList      = fmt.Errorf("value is not in list")
	errorNotMatchRegexp ErrorNotMatchRegexp = fmt.Errorf("value doesn't match regexp")
)

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

func apply(validator func(interface{}) error, fieldName string, fieldValue interface{}, result *ValidationErrors) {
	if err := validator(fieldValue); err != nil {
		*result = append(*result, ValidationError{
			Field: fieldName,
			Err:   err,
		})
	}
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("{ Field: %s, Error: %v }", v.Field, v.Err)
}

func (v ValidationErrors) Error() string {
	var res strings.Builder
	for idx, err := range v {
		_, e := res.WriteString(
			fmt.Sprintf(
				"%d) %v\n",
				idx+1,
				err,
			),
		)
		if e != nil {
			panic(ApplicationError(e))
		}
	}
	return res.String()
}

func newTagParser() *tagParser {
	return &tagParser{}
}

func (t *tagParser) initialize() {
	t.lenRegexp = regexp.MustCompile(`len:(\d+)`)
	t.reRegexp = regexp.MustCompile(`regexp:(.+)`)
	t.strInRegexp = regexp.MustCompile(`in:(.+)`)
	t.minRegexp = regexp.MustCompile(`min:(\d+)`)
	t.maxRegexp = regexp.MustCompile(`max:(\d+)`)
	t.numInRegexp = regexp.MustCompile(`in:([0-9,]+)`)
}

func (t tagParser) validateIfValueIsHigher(v interface{}, tag string) error {
	val := reflect.ValueOf(v)
	if val.Type().Kind() != reflect.Int {
		return fmt.Errorf("expected int, got %v", val)
	}

	matches := t.minRegexp.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return fmt.Errorf("invalid tag %s for int min validator", tag)
	}

	criteria := matches[1]

	min, err := strconv.Atoi(criteria)
	if err != nil {
		return fmt.Errorf("invalid tag value %s for int min validator: %w", tag, err)
	}

	if int(val.Int()) >= min {
		return nil
	}

	return errorLower
}

func (t tagParser) validateIfValueIsLower(v interface{}, tag string) error {
	val := reflect.ValueOf(v)
	if val.Type().Kind() != reflect.Int {
		return fmt.Errorf("expected int, got %v", val)
	}

	matches := t.maxRegexp.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return fmt.Errorf("invalid tag %s for int max validator", tag)
	}
	criteria := matches[1]

	max, err := strconv.Atoi(criteria)
	if err != nil {
		return fmt.Errorf("invalid tag value %s for int max validator: %w", tag, err)
	}

	if int(val.Int()) <= max {
		return nil
	}

	return errorHigher
}

func (t tagParser) validateIfValueInRange(v interface{}, tag string) error {
	val := reflect.ValueOf(v)
	if val.Type().Kind() != reflect.Int {
		return fmt.Errorf("expected int, got %v", val)
	}

	matches := t.numInRegexp.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return fmt.Errorf("invalid tag %s for int in validator", tag)
	}

	criteria := matches[1]

	splittedCriteria := strings.Split(criteria, ",")
	if len(matches) < 2 {
		return fmt.Errorf("invalid tag value %s for string in validator", tag)
	}

	if contains(splittedCriteria, strconv.Itoa(int(val.Int()))) {
		return nil
	}

	return errorNotInRange
}

func (t tagParser) validateIfValueMatchesLen(v interface{}, tag string) error {
	val := reflect.ValueOf(v)
	if val.Type().Kind() != reflect.String {
		return fmt.Errorf("expected string, got %v", val)
	}

	matches := t.lenRegexp.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return fmt.Errorf("invalid tag %s for string len validator", tag)
	}

	criteria := matches[1]

	expecteedLen, err := strconv.Atoi(criteria)
	if err != nil {
		return fmt.Errorf("invalid tag value %s for string len validator: %w", tag, err)
	}

	if len(val.String()) == expecteedLen {
		return nil
	}

	return errorLen
}

func (t tagParser) validateIfValueInList(v interface{}, tag string) error {
	val := reflect.ValueOf(v)
	if val.Type().Kind() != reflect.String {
		return fmt.Errorf("expected string, got %v", val)
	}

	matches := t.strInRegexp.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return fmt.Errorf("invalid tag %s for string in validator", tag)
	}
	criteria := matches[1]

	splittedCriteria := strings.Split(criteria, ",")

	if len(matches) < 2 {
		return fmt.Errorf("invalid tag value %s for string in validator", tag)
	}

	if contains(splittedCriteria, val.String()) {
		return nil
	}

	return errorNotInList
}

func (t tagParser) validateIfValueMatchesRegex(v interface{}, tag string) error {
	val := reflect.ValueOf(v)

	if val.Type().Kind() != reflect.String {
		return fmt.Errorf("expected string, got %v", val)
	}

	matches := t.reRegexp.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return fmt.Errorf("invalid tag %s for string re validator", tag)
	}
	criteria := matches[1]
	re, err := regexp.Compile(criteria)
	if err != nil {
		return fmt.Errorf("invalid tag value %s for string re validator: %w", tag, err)
	}
	if re.Match([]byte(val.String())) {
		return nil
	}
	return errorNotMatchRegexp
}

func (t tagParser) parseTags(tagString string, k reflect.Kind) (tagValidator, error) {
	var validator tagValidator
	tags := strings.Split(tagString, "|")
	for _, tag := range tags {
		closure := tag
		switch k { //nolint
		case reflect.Int:
			switch {
			case t.minRegexp.Match([]byte(tag)):
				validator = append(validator, func(v interface{}) error {
					return t.validateIfValueIsHigher(v, closure)
				})
			case t.maxRegexp.Match([]byte(tag)):
				validator = append(validator, func(v interface{}) error {
					return t.validateIfValueIsLower(v, closure)
				})
			case t.numInRegexp.Match([]byte(tag)):
				validator = append(validator, func(v interface{}) error {
					return t.validateIfValueInRange(v, closure)
				})
			default:
				return nil, fmt.Errorf("unknown tag %s for int validator", tag)
			}
		case reflect.String:
			switch {
			case t.lenRegexp.Match([]byte(tag)):
				validator = append(validator, func(v interface{}) error {
					return t.validateIfValueMatchesLen(v, closure)
				})
			case t.strInRegexp.Match([]byte(tag)):
				validator = append(validator, func(v interface{}) error {
					return t.validateIfValueInList(v, closure)
				})
			case t.reRegexp.Match([]byte(tag)):
				validator = append(validator, func(v interface{}) error {
					return t.validateIfValueMatchesRegex(v, closure)
				})
			default:
				return nil, fmt.Errorf("unknown tag %s for string validator", tag)
			}
		default:
			return nil, fmt.Errorf("validation of %v is not supported", k)
		}
	}

	return validator, nil
}

// Validate performs struct validation.
func Validate(v interface{}) error {
	var validationErrs ValidationErrors

	structValue := reflect.ValueOf(v)
	t := structValue.Type()

	parser := newTagParser()

	// We won`t add recover() here since regexp in parser.initialize() are highly
	// unliked to be changed
	parser.initialize()

	if t.Kind() != reflect.Struct {
		return ApplicationError(fmt.Errorf("validator expects struct, got %v", t))
	}

	for i, f := range reflect.VisibleFields(t) {
		tags, ok := f.Tag.Lookup("validate")
		k := f.Type.Kind()
		if !ok {
			continue
		}

		switch {
		case k == reflect.Slice || k == reflect.Array:
			slise := structValue.Field(i)

			sliseElKind := slise.Index(0).Type().Kind()

			validators, err := parser.parseTags(tags, sliseElKind)
			if err != nil {
				return ApplicationError(err)
			}

			for sliceIdx := 0; sliceIdx < slise.Len(); sliceIdx++ {
				sliceElValue := slise.Index(sliceIdx).Interface()
				for _, validator := range validators {
					apply(validator, t.Field(i).Name, sliceElValue, &validationErrs)
				}
			}

		case k == reflect.Int || k == reflect.String:
			validators, err := parser.parseTags(tags, k)
			if err != nil {
				return ApplicationError(err)
			}

			fieldValue := structValue.Field(i).Interface()

			for _, validator := range validators {
				apply(validator, t.Field(i).Name, fieldValue, &validationErrs)
			}

		default:
			return ApplicationError(fmt.Errorf("cannot validate %v", k))
		}
	}

	switch len(validationErrs) {
	case 0:
		return nil
	default:
		return validationErrs
	}
}
