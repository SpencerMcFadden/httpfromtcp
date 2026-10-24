package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const clrf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	idx := bytes.Index(raw, []byte(clrf))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}
	requestLineText := string(raw[:idx])

	requestParts := strings.Split(requestLineText, clrf)

	requestLine, err := parseRequestLine(requestParts[0])
	if err != nil {
		return nil, err
	}

	finalRequest := Request{
		RequestLine: *requestLine,
	}
	return &finalRequest, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	lineParts := strings.Split(line, " ")
	if len(lineParts) != 3 {
		return nil, fmt.Errorf("incorrect parts length for request-line: %s", line)
	}

	method := lineParts[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := lineParts[1]

	versionParts := strings.Split(lineParts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("incorrectly formatted start-line: %s", line)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	finalLine := &RequestLine{
		HttpVersion:   version,
		RequestTarget: requestTarget,
		Method:        method,
	}
	return finalLine, nil
}
