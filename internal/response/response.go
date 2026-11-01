package response

import (
	"fmt"
	"io"
	"strconv"

	"httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

const clrf = "\r\n"

func getStatusLine(statusCode StatusCode) []byte {
	statusPhrase := ""
	switch statusCode {
	case OK:
		statusPhrase = "OK"
	case BadRequest:
		statusPhrase = "Bad Request"
	case InternalError:
		statusPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s"+clrf, statusCode, statusPhrase))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(key + ": " + value + clrf))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte(clrf))
	return err
}
