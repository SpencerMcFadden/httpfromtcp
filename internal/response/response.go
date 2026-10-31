package response

import (
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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case OK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK" + clrf))
		if err != nil {
			return err
		}
	case BadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request" + clrf))
		if err != nil {
			return err
		}
	case InternalError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error" + clrf))
		if err != nil {
			return err
		}
	default:
		_, err := w.Write([]byte(""))
		if err != nil {
			return err
		}
	}
	return nil
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
	if err != nil {
		return err
	}
	return nil
}
