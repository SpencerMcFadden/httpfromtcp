package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const (
	clrf         = "\r\n"
	specialChars = "!#$%&'*+-.^_`|~"
)

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(clrf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	lineParts := bytes.SplitN(data[:idx], []byte(":"), 2)
	if len(lineParts) != 2 {
		return 0, false, fmt.Errorf("incorrect parts length for header: %s", data[:idx])
	}

	fieldName := string(lineParts[0])
	if strings.ContainsAny(fieldName, " \t") || len(fieldName) < 1 {
		return 0, false, fmt.Errorf("invalid header name: %s", fieldName)
	}

	fieldNameBytes := lineParts[0]
	for _, c := range fieldNameBytes {
		if !isValidChar(c) {
			return 0, false, fmt.Errorf("invalid character in header name: %c", c)
		}
	}

	fieldName = strings.TrimSpace(fieldName)
	fieldValue := bytes.TrimSpace(lineParts[1])

	h.Set(fieldName, string(fieldValue))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func isValidChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}
	if strings.IndexByte(specialChars, c) >= 0 {
		return true
	}
	return false
}
