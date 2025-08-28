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

type ValidationErrors []ValidationError

// Ошибки валидаторов.
var (
	ErrValidation      = errors.New("validation error")
	ErrLength          = errors.New("length validation failed")
	ErrRegexp          = errors.New("regexp validation failed")
	ErrIn              = errors.New("in validation failed")
	ErrMin             = errors.New("min validation failed")
	ErrMax             = errors.New("max validation failed")
	ErrUnsupportedType = errors.New("unsupported field type")
)

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "no validation errors"
	}

	var sb strings.Builder
	sb.WriteString("validation errors:\n")
	for i, err := range v {
		sb.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, err.Field, err.Err.Error()))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input must be a struct")
	}
	var validationErrors ValidationErrors
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		// Пропускаем непубличные поля.
		if !field.IsExported() {
			continue
		}
		// Получаем тег validate.
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		fieldName := field.Name

		// Обрабатываем слайсы.
		if fieldValue.Kind() == reflect.Slice {
			for j := 0; j < fieldValue.Len(); j++ {
				element := fieldValue.Index(j)
				elementName := fmt.Sprintf("%s[%d]", fieldName, j)

				if errs := validateField(elementName, element, validateTag); errs != nil {
					validationErrors = append(validationErrors, errs...)
				}
			}
			continue
		}

		// Обрабатываем обычные поля.
		if errs := validateField(fieldName, fieldValue, validateTag); errs != nil {
			validationErrors = append(validationErrors, errs...)
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(fieldName string, fieldValue reflect.Value, validateTag string) ValidationErrors {
	var errors ValidationErrors
	rules := strings.Split(validateTag, "|")

	for _, rule := range rules {
		ruleParts := strings.SplitN(rule, ":", 2)
		if len(ruleParts) != 2 {
			continue
		}

		validator := ruleParts[0]
		param := ruleParts[1]

		var err error
		switch fieldValue.Kind() {
		case reflect.String:
			err = validateString(fieldValue.String(), validator, param)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = validateInt(fieldValue.Int(), validator, param)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = validateUint(fieldValue.Uint(), validator, param)
		case reflect.Float32, reflect.Float64:
			err = validateFloat(fieldValue.Float(), validator, param)
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Ptr, reflect.Struct, reflect.UnsafePointer,
			reflect.Bool, reflect.Complex64, reflect.Complex128, reflect.Uintptr,
			reflect.Invalid:
			// Все неподдерживаемые типы обрабатываются здесь.
			err = fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}

		if err != nil {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	}

	return errors
}

func validateString(value, validator, param string) error {
	switch validator {
	case "len":
		expectedLen, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid length parameter: %w", err)
		}
		if len(value) != expectedLen {
			return fmt.Errorf("%w: expected length %d, got %d", ErrLength, expectedLen, len(value))
		}
	case "regexp":
		regex, err := regexp.Compile(param)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
		if !regex.MatchString(value) {
			return fmt.Errorf("%w: value '%s' doesn't match pattern '%s'", ErrRegexp, value, param)
		}
	case "in":
		allowedValues := strings.Split(param, ",")
		found := false
		for _, allowed := range allowedValues {
			if value == allowed {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%w: value '%s' not in allowed values: %s", ErrIn, value, param)
		}
	default:
		return fmt.Errorf("unknown string validator: %s", validator)
	}
	return nil
}

func validateInt(value int64, validator, param string) error {
	switch validator {
	case "min":
		min, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid min parameter: %w", err)
		}
		if value < min {
			return fmt.Errorf("%w: value %d is less than minimum %d", ErrMin, value, min)
		}
	case "max":
		max, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid max parameter: %w", err)
		}
		if value > max {
			return fmt.Errorf("%w: value %d is greater than maximum %d", ErrMax, value, max)
		}
	case "in":
		return validateInInt(value, param)
	default:
		return fmt.Errorf("unknown int validator: %s", validator)
	}
	return nil
}

func validateUint(value uint64, validator, param string) error {
	switch validator {
	case "min":
		min, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid min parameter: %w", err)
		}
		if value < min {
			return fmt.Errorf("%w: value %d is less than minimum %d", ErrMin, value, min)
		}
	case "max":
		max, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid max parameter: %w", err)
		}
		if value > max {
			return fmt.Errorf("%w: value %d is greater than maximum %d", ErrMax, value, max)
		}
	case "in":
		return validateInUint(value, param)
	default:
		return fmt.Errorf("unknown uint validator: %s", validator)
	}
	return nil
}

// Общая функция для валидации in для int64.
func validateInInt(value int64, param string) error {
	allowedValues := strings.Split(param, ",")
	found := false
	for _, allowedStr := range allowedValues {
		allowed, err := strconv.ParseInt(allowedStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid in parameter: %w", err)
		}
		if value == allowed {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%w: value %d not in allowed values: %s", ErrIn, value, param)
	}
	return nil
}

// Общая функция для валидации in для uint64.
func validateInUint(value uint64, param string) error {
	allowedValues := strings.Split(param, ",")
	found := false
	for _, allowedStr := range allowedValues {
		allowed, err := strconv.ParseUint(allowedStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid in parameter: %w", err)
		}
		if value == allowed {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%w: value %d not in allowed values: %s", ErrIn, value, param)
	}
	return nil
}

func validateFloat(value float64, validator, param string) error {
	switch validator {
	case "min":
		min, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return fmt.Errorf("invalid min parameter: %w", err)
		}
		if value < min {
			return fmt.Errorf("%w: value %f is less than minimum %f", ErrMin, value, min)
		}
	case "max":
		max, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return fmt.Errorf("invalid max parameter: %w", err)
		}
		if value > max {
			return fmt.Errorf("%w: value %f is greater than maximum %f", ErrMax, value, max)
		}
	case "in":
		allowedValues := strings.Split(param, ",")
		found := false
		for _, allowedStr := range allowedValues {
			allowed, err := strconv.ParseFloat(allowedStr, 64)
			if err != nil {
				return fmt.Errorf("invalid in parameter: %w", err)
			}
			if value == allowed {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%w: value %f not in allowed values: %s", ErrIn, value, param)
		}
	default:
		return fmt.Errorf("unknown float validator: %s", validator)
	}
	return nil
}
