package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type state int

const (
	initialized state = iota
	done
)

type Request struct {
	RequestLine RequestLine
	State       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	bufferSize = 8
	clrf       = "\r\n"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := &Request{
		State: initialized,
	}
	for request.State != done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = done
				break
			}
			return nil, err
		}
		readToIndex += n
		parsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing: %s", err.Error())
		}
		copy(buf, buf[parsed:readToIndex])
		readToIndex -= parsed
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(clrf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case initialized:
		result, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *result
		r.State = done
		return n, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func requestLineFromString(line string) (*RequestLine, error) {
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
