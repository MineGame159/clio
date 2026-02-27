package rd

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func get[T any](token string, url string) (T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		var empty T
		return empty, err
	}

	return doRequest[T](token, req)
}

func post[T any](token string, url string, values url.Values) (T, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		var empty T
		return empty, err
	}

	return doRequest[T](token, req)
}

func doRequest[T any](token string, req *http.Request) (T, error) {
	req.Header.Set("User-Agent", "clio")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		var empty T
		return empty, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body struct{ Error string }
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			var empty T
			return empty, err
		}

		var empty T
		return empty, errors.New(body.Error)
	}

	var body T
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		var empty T
		return empty, err
	}

	return body, nil
}
