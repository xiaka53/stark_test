package request

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Get(funcName string, parametersData map[string]any) ([]byte, error) {
	baseURL := fmt.Sprintf("%s%s", BASE_URL, funcName)

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	params := url.Values{}
	for key, value := range parametersData {
		params.Add(key, fmt.Sprintf("%v", value))
	}
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("quest err")
	}
	return body, nil
}
