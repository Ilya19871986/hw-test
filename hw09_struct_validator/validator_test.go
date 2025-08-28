package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

//nolint:all
func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid User struct",
			in: User{
				ID:     "123456789012345678901234567890123456", // 36 символов.
				Name:   "John",                                 // Нет валидации.
				Age:    25,                                     // Между 18-50.
				Email:  "test@example.com",                     // Валидный email по regex.
				Role:   "admin",                                // В разрешенных значениях.
				Phones: []string{"79123456789", "79234567890"}, // По 11 символов каждый.
			},
			expectedErr: nil,
		},
		{
			name: "invalid User - multiple errors",
			in: User{
				ID:     "short",                               // Слишком короткий (5 символов).
				Name:   "John",                                // Нет валидации.
				Age:    16,                                    // Слишком молод.
				Email:  "invalid-email",                       // Невалидный email.
				Role:   "guest",                               // Не в разрешенных значениях.
				Phones: []string{"79123456789", "7923456789"}, // Второй телефон имеет 10 символов.
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("%w: expected length 36, got 5", ErrLength)},
				{Field: "Age", Err: fmt.Errorf("%w: value 16 is less than minimum 18", ErrMin)},
				{Field: "Email", Err: fmt.Errorf("%w: value 'invalid-email' doesn't match pattern '^\\w+@\\w+\\.\\w+$'", ErrRegexp)}, //nolint:all
				{Field: "Role", Err: fmt.Errorf("%w: value 'guest' not in allowed values: admin,stuff", ErrIn)},
				{Field: "Phones[1]", Err: fmt.Errorf("%w: expected length 11, got 10", ErrLength)},
			},
		},
		{
			name: "valid App struct",
			in: App{
				Version: "1.2.3", // 5 символов.
			},
			expectedErr: nil,
		},
		{
			name: "invalid App - version too long",
			in: App{
				Version: "1.2.3.4", // 7 символов.
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("%w: expected length 5, got 7", ErrLength)},
			},
		},
		{
			name: "Token struct without validation tags",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			name: "valid Response struct",
			in: Response{
				Code: 200,  // В разрешенных значениях.
				Body: "OK", // Нет валидации.
			},
			expectedErr: nil,
		},
		{
			name: "invalid Response - code not in allowed values",
			in: Response{
				Code: 400, // Не в 200,404,500.
				Body: "Bad Request",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("%w: value 400 not in allowed values: 200,404,500", ErrIn)},
			},
		},
		{
			name: "empty User slice validation",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{}, // Пустой слайс - должен пройти валидацию.
			},
			expectedErr: nil,
		},
		{
			name:        "nil input (should error)",
			in:          nil,
			expectedErr: errors.New("input must be a struct"),
		},
		{
			name:        "non-struct input (should error)",
			in:          "not a struct",
			expectedErr: errors.New("input must be a struct"),
		},
		{
			name:        "pointer to struct (should error)",
			in:          &User{ID: "123456789012345678901234567890123456"},
			expectedErr: errors.New("input must be a struct"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)

			// Ожидается отсутствие ошибки.
			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("Ожидалось отсутствие ошибки, получено: %v", err)
				}
				return
			}

			// Ожидается наличие ошибки.
			if err == nil {
				t.Error("Ожидалась ошибка, получено nil")
				return
			}

			// Проверка для не-ValidationErrors (например, "input must be a struct").
			if expectedErr, ok := tt.expectedErr.(error); ok {
				var validationErrors ValidationErrors
				if !errors.As(err, &validationErrors) {
					// Это простая ошибка, сравниваем сообщения.
					if err.Error() != expectedErr.Error() {
						t.Errorf("Ожидалась ошибка %v, получено %v", expectedErr, err)
					}
					return
				}
			}

			// Проверка для ValidationErrors.
			var expectedErrors ValidationErrors
			if !errors.As(tt.expectedErr, &expectedErrors) {
				t.Errorf("Неверный тестовый случай: expectedErr должен быть ValidationErrors или error")
				return
			}

			var actualErrors ValidationErrors
			if !errors.As(err, &actualErrors) {
				t.Errorf("Ожидались ValidationErrors, получено %T: %v", err, err)
				return
			}

			// Проверка количества ошибок.
			if len(actualErrors) != len(expectedErrors) {
				t.Errorf("Ожидалось %d ошибок, получено %d. Ошибки: %v", len(expectedErrors), len(actualErrors), actualErrors)
				return
			}

			// Проверка каждой ошибки.
			for j, expectedErr := range expectedErrors {
				actualErr := actualErrors[j]

				// Проверка имени поля.
				if actualErr.Field != expectedErr.Field {
					t.Errorf("Ошибка %d: ожидалось поле %s, получено %s", j, expectedErr.Field, actualErr.Field)
				}

				// Проверка сообщения об ошибке.
				if actualErr.Err.Error() != expectedErr.Err.Error() {
					t.Errorf("Ошибка %d: ожидалось сообщение '%s', получено '%s'", j, expectedErr.Err.Error(), actualErr.Err.Error())
				}
			}
		})
	}
}
