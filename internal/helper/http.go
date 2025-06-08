package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mercury/internal/pkg/logger"
	"net/http"
	"strconv"
)

type HttpMethod string

const (
	POST HttpMethod = "POST"
	GET  HttpMethod = "GET"
)

func HttpRequest[Rp interface{}](url string, body interface{}, headers map[string]string, method HttpMethod) (Rp, error) {

	var response Rp

	client := &http.Client{}

	// Convert body to JSON
	bodyBytes, err := json.Marshal(body)

	var req *http.Request
	if method == GET {
		req, err = http.NewRequest(string(method), url, nil)
		if err != nil {
			return response, err
		}
	} else {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return response, err
		}
		req, err = http.NewRequest(string(method), url, io.NopCloser(bytes.NewReader(bodyBytes)))
		if err != nil {
			return response, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	logger.Debug("HTTP", method, "request to", url, "with body:", string(bodyBytes))
	responseByte, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("error closing response body on ", url, " :", err)
		}
	}(responseByte.Body)

	if responseByte.StatusCode != http.StatusOK {
		err = errors.New("HTTP " + string(method) + " request failed with status code:" + strconv.Itoa(responseByte.StatusCode))
		logger.Error(err.Error())
		return response, err
	}

	responseBody, err := io.ReadAll(responseByte.Body)
	if err != nil {
		logger.Error("error reading response body on ", url, " :", err)
		return response, err
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		logger.Error("error unmarshalling response body on ", url, " :", err, " response body:", string(responseBody))
		return response, errors.New("error unmarshalling response body: " + err.Error())
	}

	logger.Debug("HTTP ", method, " request to", url, "succeeded with response:", response)
	return response, nil
}
