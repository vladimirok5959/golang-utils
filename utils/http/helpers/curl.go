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

type CurlGetOpts struct {
	ExpectStatusCode int
	Headers          map[string][]string
	Timeout          time.Duration
}

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

func CurlGet(ctx context.Context, url string, opts *CurlGetOpts) ([]byte, error) {
	if opts == nil {
		opts = &CurlGetOpts{
			ExpectStatusCode: http.StatusOK,
			Headers:          nil,
			Timeout:          time.Second * 60,
		}
	}

	var b []byte

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return b, err
	}
	req.Header = opts.Headers

	rcl := &http.Client{}
	var resp *http.Response
	resp, err = rcl.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	if opts.ExpectStatusCode > 0 {
		if resp.StatusCode != opts.ExpectStatusCode {
			return b, error(&CurlGetStatusError{
				Expected: opts.ExpectStatusCode,
				Received: resp.StatusCode,
			})
		}
	}

	return b, nil
}
