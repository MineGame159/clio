package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/http"
	"path"
	"strings"
	"unicode"
)

var videoFileExtensions = []string{".webm", ".mkv", ".flv", ".avi", ".mov", ".mp4"}

func IsVideoFile(name string) bool {
	ext := strings.ToLower(path.Ext(name))

	for _, videoExt := range videoFileExtensions {
		if ext == videoExt {
			return true
		}
	}

	return false
}

func Capitalize(input string) string {
	runes := []rune(input)
	var b strings.Builder

	for i, r := range runes {
		if i == 0 {
			b.WriteRune(unicode.ToTitle(r))
			continue
		}

		if unicode.IsUpper(r) {
			prev := runes[i-1]

			if !unicode.IsUpper(prev) {
				b.WriteRune(' ')
			} else {
				if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					b.WriteRune(' ')
				}
			}
		}

		b.WriteRune(r)
	}

	return b.String()
}

func Count[T any](it iter.Seq[T]) uint {
	count := uint(0)

	for range it {
		count++
	}

	return count
}
func GetJson[T any](url string) (T, error) {
	return DoReqJsonCtx[T](context.Background(), "GET", url, "", nil)
}

func GetJsonCtx[T any](ctx context.Context, url string) (T, error) {
	return DoReqJsonCtx[T](ctx, "GET", url, "", nil)
}

func DoReqJsonCtx[T any](ctx context.Context, method, url, bodyContentType string, body io.Reader) (T, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		var empty T
		return empty, err
	}

	req.Header.Set("User-Agent", "clio")
	req.Header.Set("Accept", "application/json")

	if bodyContentType != "" {
		req.Header.Set("Content-Type", bodyContentType)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		var empty T
		return empty, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var empty T
		return empty, errors.New(fmt.Sprintf("request failed with status code: %d '%s'", res.StatusCode, res.Status))
	}

	var resBody T
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		var empty T
		return empty, err
	}

	return resBody, nil
}
