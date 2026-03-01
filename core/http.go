package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

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

func ListenAndServe(handler http.Handler) (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	//goland:noinspection GoUnhandledErrorResult
	go http.Serve(listener, handler)

	address := "http://localhost:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	return address, nil
}

func WriteJson(res http.ResponseWriter, data any) {
	res.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(res).Encode(data)
}

func WriteError(res http.ResponseWriter, msg string, code int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)

	_, _ = fmt.Fprintf(res, "{\"error\":\"%s\"}", msg)
}
