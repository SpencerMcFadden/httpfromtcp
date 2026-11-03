package response

import (
	"fmt"
)

type StatusCode int

const (
	OK            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

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
