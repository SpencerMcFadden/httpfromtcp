package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
)

const port = 42069

func main() {
	srv, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		handlerProxy(w, req)
		return
	}
	switch req.RequestLine.RequestTarget {
	case "/video":
		handlerVideo(w, req)
		return
	case "/yourproblem":
		handler400(w, req)
		return
	case "/myproblem":
		handler500(w, req)
		return
	default:
		handler200(w, req)
		return
	}
}

func handlerVideo(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.OK)
	const filepath = "assets/vim.mp4"
	videoBytes, err := os.ReadFile(filepath)
	if err != nil {
		handler500(w, nil)
		return
	}

	h := response.GetDefaultHeaders(len(videoBytes))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(videoBytes)
}

func handlerProxy(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	baseURL := "https://httpbingo.org/"
	res, err := http.Get(baseURL + target)
	if err != nil {
		handler500(w, req)
		return
	}
	defer res.Body.Close()

	w.WriteStatusLine(response.OK)
	defaultHeaders := response.GetDefaultHeaders(0)
	defaultHeaders.Delete("Content-Length")
	defaultHeaders.Set("Transfer-Encoding", "chunked")
	defaultHeaders.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(defaultHeaders)

	buf := make([]byte, 1024)
	var totalBody []byte
	for {
		n, err := res.Body.Read(buf)
		if n > 0 {
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Println("error writing chunked body:", err)
				break
			}
			totalBody = append(totalBody, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error reading response:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Println("error writing chunked body done:", err)
	}
	responseHash := sha256.Sum256(totalBody)
	trailers := headers.Headers{}
	trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", responseHash))
	trailers.Set("X-Content-Length", strconv.Itoa(len(totalBody)))
	err = w.WriteTrailers(trailers)
	if err != nil {
		log.Println("error writing trailers:", err)
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.BadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.InternalError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.OK)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}
