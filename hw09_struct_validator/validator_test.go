package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		// meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Status struct {
		Name string `validate:"in:online,offline"`
	}

	IP struct {
		Address string `validate:"regexp:^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$"`
	}

	Enum struct {
		Num int `validate:"min:0"`
	}
	Char struct {
		Num int `validate:"max:255"`
	}
	Bit struct {
		Num int `validate:"in:0,1"`
	}

	Byte struct {
		Num int `validate:"min:0|max:7"`
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

func TestValidatePositive(t *testing.T) {
	tests := []interface{}{
		App{
			Version: "v1.00",
		},
		Status{
			Name: "online",
		},
		Status{
			Name: "offline",
		},
		IP{
			Address: "192.168.0.1",
		},
		Enum{
			Num: 1,
		},
		Char{
			Num: 127,
		},
		Bit{
			Num: 0,
		},
		Bit{
			Num: 1,
		},
		Byte{
			Num: 0,
		},
		Byte{
			Num: 7,
		},
		Token{
			Header:    []byte("Host:127.0.0.1"),
			Payload:   []byte("foobar"),
			Signature: []byte("Zm9vYmFyCg=="),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)
			require.NoError(t, err)
		})
	}
}

func TestValidateRegression(t *testing.T) {
	tests := []interface{}{
		Response{
			Code: 200,
			Body: "foobar",
		},
		User{
			ID:     "123456789012345678901234567890123456",
			Name:   "test",
			Age:    29,
			Email:  "foo@bar.baz",
			Role:   "stuff",
			Phones: []string{"+1234567890", "+1234567891"},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)
			require.NoError(t, err)
		})
	}
}

func TestValidateNegative(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: App{
				Version: "v1.00.00",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errorLen,
				},
			},
		},
		{
			in: App{
				Version: "v1",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errorLen,
				},
			},
		},
		{
			in: Status{
				Name: "AFK",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Name",
					Err:   errorNotInList,
				},
			},
		},
		{
			in: IP{
				Address: "4192.5168.3450.1111",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Address",
					Err:   errorNotMatchRegexp,
				},
			},
		},
		{
			in: Char{
				Num: 300,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Num",
					Err:   errorHigher,
				},
			},
		},
		{
			in: Byte{
				Num: 8,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Num",
					Err:   errorHigher,
				},
			},
		},
		{
			in: Response{
				Code: 499,
				Body: "client close the connection",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errorNotInRange,
				},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "test",
				Age:    13,
				Email:  "foo@bar.baz",
				Role:   "stuff",
				Phones: []string{"123", "234"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   errorLower,
				},
				ValidationError{
					Field: "Phones",
					Err:   errorLen,
				},
				ValidationError{
					Field: "Phones",
					Err:   errorLen,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			actualErr := Validate(tt.in)
			require.True(t, actualErr != nil)

			actualErrList := ValidationErrors{}

			if !errors.As(actualErr, &actualErrList) {
				t.Fail()
			}

			if len(tt.expectedErr) != len(actualErrList) {
				t.Fail()
			}

			for i, e := range tt.expectedErr {
				require.Equal(t, e, actualErrList[i])
			}
		})
	}
}
