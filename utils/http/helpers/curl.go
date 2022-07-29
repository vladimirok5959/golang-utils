package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrCurlGetStatus = error(&CurlGetStatusError{})
)

type CurlGetStatusError struct {
	Expected int
	Received int
}

func (e *CurlGetStatusError) Is(err error) bool {
	if _, ok := err.(*CurlGetStatusError); !ok {
		return false
	}
	return true
}

func (e *CurlGetStatusError) Error() string {
	return fmt.Sprintf("CurlGet: expected %d, received %d", e.Expected, e.Received)
}

func CurlGet(ctx context.Context, url string, timeout time.Duration) ([]byte, error) {
	var b []byte

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return b, err
	}

	rcl := &http.Client{}
	var resp *http.Response
	resp, err = rcl.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b, error(&CurlGetStatusError{
			Expected: http.StatusOK,
			Received: resp.StatusCode,
		})
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return b, nil
}
