package logger

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	Content []byte
	Size    int
	Status  int
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if RollBarEnabled {
		if !RollBarSkipStatusCodes.Contain(w.Status) {
			w.Content = append(w.Content, b...)
		}
	}
	size, err := w.ResponseWriter.Write(b)
	w.Size += size
	return size, err
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}
