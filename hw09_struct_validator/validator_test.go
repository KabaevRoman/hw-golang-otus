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
	InvalidUser struct {
		Age string `validate:"min:18"`
	}
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

func TestValidate(t *testing.T) {
	testsInvalidInput := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: InvalidUser{
				Age: "18",
			},
			expectedErr: ErrClient,
		},
	}
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "11111111-1111-1111-1111-123456789abc",
				Name:   "Cat Dog",
				Age:    25,                      // Valid Age
				Email:  "v@example.com",         // Valid Email
				Role:   "admin",                 // Valid Role
				Phones: []string{"12345678901"}, // Valid Phone
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "invalid-id-not-36-chars",
				Name:   "Invalid User",
				Age:    17,
				Email:  "invalid-email",
				Role:   "nonexistent-role",
				Phones: []string{"123456789"},
			},
			expectedErr: ValidationErrors{
				{"ID", ErrLen},
				{"Age", ErrMin},
				{"Email", ErrRegexp},
				{"Role", ErrIn},
				{"Phones", ErrLen},
			},
		},
		{
			in: App{
				Version: "v1.0",
			},
			expectedErr: ValidationErrors{
				{"Version", fmt.Errorf("length must be 5")},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 800,
				Body: "Error",
			},
			expectedErr: ValidationErrors{
				{"Code", fmt.Errorf("must be one of: 200, 404, 500")},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in:          struct{}{},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "",
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			var validationErrs, expectedValidationErrs ValidationErrors
			errors.As(err, &validationErrs)
			errors.As(err, &expectedValidationErrs)
			for j := 0; j < len(expectedValidationErrs); j++ {
				errors.Is(validationErrs[j], expectedValidationErrs[j])
			}
		})
	}
	for i, tt := range testsInvalidInput {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			errors.Is(err, tt.expectedErr)
		})
	}
}
