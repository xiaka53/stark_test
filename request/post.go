package request

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

const (
	BASE_URL = "https://starknet-mainnet.public.blastapi.io/rpc"
)

func Post(jsonData []byte) ([]byte, error) {
	url := BASE_URL

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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
