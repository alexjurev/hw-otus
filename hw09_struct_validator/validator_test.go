package hw09structvalidator

import (
	"encoding/json"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint
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

type testCases []struct {
	in          interface{}
	expectedErr error
}

func TestValidate(t *testing.T) {
	tests := testCases{
		{
			in: App{
				Version: "v1.00",
			},
			expectedErr: nil,
		},
		{
			in: Status{
				Name: "online",
			},
			expectedErr: nil,
		},
		{
			in: IP{
				Address: "192.168.0.1",
			},
			expectedErr: nil,
		},
		{
			in: Enum{
				Num: 1,
			},
			expectedErr: nil,
		},
		{
			in: Char{
				Num: 127,
			},
			expectedErr: nil,
		},
		{
			in: Bit{
				Num: 0,
			},
			expectedErr: nil,
		},
		{
			in: Byte{
				Num: 0,
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("Host:127.0.0.1"),
				Payload:   []byte("foobar"),
				Signature: []byte("Zm9vYmFyCg=="),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "foobar",
			},
			expectedErr: nil,
		},
		{
			in:          getCorrectUser(),
			expectedErr: nil,
		},
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
			in: getWrongUser(),
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

			err := Validate(tt.in)

			require.Equal(t, err, tt.expectedErr)
		})
	}
}

func getCorrectUser() User {
	return User{
		ID:     "123456789012345678901234567890123456",
		Name:   "test",
		Age:    29,
		Email:  "foo@bar.baz",
		Role:   "stuff",
		Phones: []string{"+1234567890", "+1234567891"},
	}
}

func getWrongUser() User {
	return User{
		ID:     "123456789012345678901234567890123456",
		Name:   "test",
		Age:    13,
		Email:  "foo@bar.baz",
		Role:   "stuff",
		Phones: []string{"123", "234"},
	}
}
