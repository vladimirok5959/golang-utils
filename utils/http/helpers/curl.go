package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CurlGetStatusError struct {
	expected int
	received int
}

func (c *CurlGetStatusError) Error() string {
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
		return b, &CurlGetStatusError{
			expected: http.StatusOK,
			received: resp.StatusCode,
		}
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return b, nil
}
