package logger

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	Content []byte
	Size    int
	Status  int
}

func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("hijack not supported")
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if RollBarEnabled && !RollBarSkipStatusCodes.contain(w.Status) {
		w.Content = append(w.Content, b...)
	}
	size, err := w.ResponseWriter.Write(b)
	w.Size += size
	return size, err
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}
