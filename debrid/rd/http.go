package rd

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unsafe"
)

func get[T any](c *Client, ctx context.Context, url string) (T, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		var empty T
		return empty, err
	}

	return doRequest[T](c, ctx, req)
}

func post[T any](c *Client, ctx context.Context, url string, values url.Values) (T, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		var empty T
		return empty, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return doRequest[T](c, ctx, req)
}

func doRequest[T any](c *Client, ctx context.Context, req *http.Request) (T, error) {
	_ = c.sem.Acquire(ctx, 1)
	defer c.sem.Release(1)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	waitTime := time.Millisecond * 256

	for i := 0; ; i++ {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			var empty T
			return empty, err
		}

		if res.StatusCode == http.StatusTooManyRequests && i < 7 {
			_ = res.Body.Close()

			time.Sleep(waitTime)
			waitTime += time.Millisecond * 100

			if req.GetBody != nil {
				if body, err := req.GetBody(); err == nil {
					req.Body = body
				}
			}

			continue
		}

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			var body struct{ Error string }
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				_ = res.Body.Close()
				var empty T
				return empty, err
			}

			_ = res.Body.Close()
			var empty T
			return empty, errors.New(body.Error)
		}

		var empty T
		if unsafe.Sizeof(empty) == 0 {
			return empty, nil
		}

		var body T
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			_ = res.Body.Close()
			var empty T
			return empty, err
		}

		_ = res.Body.Close()
		return body, nil
	}
}
