package response

import (
	"fmt"
	"io"
	"strconv"

	"httpfromtcp/internal/headers"
)

type (
	StatusCode  int
	writerState int
)

type Writer struct {
	Input       io.Writer
	writerState writerState
}

const (
	WritingStatusLine writerState = iota
	WritingHeaders
	WritingBody
	WritingDone
)

const (
	OK            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

const clrf = "\r\n"

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Input: w,
	}
}

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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != WritingStatusLine {
		return fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingStatusLine)
	}
	_, err := w.Input.Write(getStatusLine(statusCode))
	if err != nil {
		return err
	}
	w.writerState = WritingHeaders
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != WritingHeaders {
		return fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingHeaders)
	}
	for key, value := range headers {
		_, err := w.Input.Write([]byte(key + ": " + value + clrf))
		if err != nil {
			return err
		}
	}
	_, err := w.Input.Write([]byte(clrf))
	if err != nil {
		return err
	}
	w.writerState = WritingBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingBody)
	}
	numBytes, err := w.Input.Write(p)
	if err != nil {
		return 0, err
	}
	w.writerState = WritingDone
	return numBytes, nil
}
