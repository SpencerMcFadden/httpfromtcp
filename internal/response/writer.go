package response

import (
	"fmt"
	"io"

	"httpfromtcp/internal/headers"
)

type (
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

const clrf = "\r\n"

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: WritingStatusLine,
		Input:       w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != WritingStatusLine {
		return fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingStatusLine)
	}

	defer func() { w.writerState = WritingHeaders }()
	_, err := w.Input.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != WritingHeaders {
		return fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingHeaders)
	}

	defer func() { w.writerState = WritingBody }()
	for key, value := range headers {
		_, err := w.Input.Write([]byte(key + ": " + value + clrf))
		if err != nil {
			return err
		}
	}
	_, err := w.Input.Write([]byte(clrf))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingBody)
	}

	defer func() { w.writerState = WritingDone }()
	return w.Input.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingBody)
	}

	hexLength := fmt.Sprintf("%x"+clrf, len(p))
	toWrite := []byte(hexLength)
	p = append(p, []byte(clrf)...)
	toWrite = append(toWrite, p...)
	return w.Input.Write(toWrite)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("writing state: %d\nexpected: %d", w.writerState, WritingBody)
	}

	defer func() { w.writerState = WritingDone }()
	return w.Input.Write([]byte("0" + clrf + clrf))
}
