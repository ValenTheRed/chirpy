package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	apiKey, found := strings.CutPrefix(header.Get("Authorization"), "ApiKey ")
	if !found {
		return "", fmt.Errorf("auth: API key is missing in Authorization")
	}
	return apiKey, nil
}
