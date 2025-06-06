package helper

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"mercury/internal/model/response/responsecode"
	"mercury/internal/pkg"
	"net"
	"net/http"
	"strings"
)

func DecodeAndValidateRequest(r *http.Request, dst interface{}) error {

	v := pkg.GetValidator()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errors.New("invalid request payload: " + err.Error())
	}

	if err := v.Struct(dst); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	return nil
}

func ValidatePayload(p []byte, dst interface{}) error {

	v := pkg.GetValidator()

	if err := json.Unmarshal(p, &dst); err != nil {
		return errors.New("invalid payload: " + err.Error())
	}

	if err := v.Struct(dst); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to decode request: "+err.Error(), http.StatusBadRequest)
	}
}

func GetResponseCodeName(code responsecode.ResponseCode) string {
	if name, exists := responsecode.ResponseCodeNames[code]; exists {
		return name
	}
	return "Unknown Code"
}

func GetClientIP(r *http.Request) string {
	// If the request is behind a reverse proxy, the IP address might be forwarded in the X-Forwarded-For header.
	// First, check for the X-Forwarded-For header.
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		// The X-Forwarded-For header contains a comma-separated list of IPs
		// The first IP in the list is the original client IP.
		return strings.Split(ips, ",")[0]
	}

	// Otherwise, fallback to the remote address.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func StructToMapStringArray(data interface{}) ([]map[string]string, error) {
	var result []map[string]string

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func GenerateRandomAlphaNumericString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if length <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
