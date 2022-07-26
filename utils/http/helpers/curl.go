package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrCurlGetStatus = errCurlGetStatus(0, 0)
)

func errCurlGetStatus(e, r int) error {
	return &curlGetStatusError{
		expected: e,
		received: r,
	}
}

type curlGetStatusError struct {
	expected int
	received int
}

func (c *curlGetStatusError) Error() string {
	return fmt.Sprintf("CurlGet: expected %d, received %d", c.received, c.expected)
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
		return b, errCurlGetStatus(
			http.StatusOK,
			resp.StatusCode,
		)
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return b, nil
}
