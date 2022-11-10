package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// User represents user data.
type User struct {
	Email string `json:"email,omitempty"`
}

// DomainStat represents domains data.
type DomainStat map[string]int

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var user User

// GetDomainStat gets domain users count.
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	res := make(DomainStat)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return DomainStat{}, fmt.Errorf("cannot unmarshal data: %w", err)
		}

		match := strings.Split(user.Email, "@")[1]
		if len(match) < 2 {
			return DomainStat{}, fmt.Errorf("unexpected data")
		}

		if strings.Contains(match, "."+domain) {
			res[strings.ToLower(match)]++
		}
	}

	return res, nil
}
