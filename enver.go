package yadsharez

import (
	"fmt"
	"os"
)

const (
	tokenEnv         = "OAUTH_TOKEN"
	tokenNotFoundErr = "OAUTH_TOKEN env missing"
)

// GetOAuthToken looks up into env var and returns error if OAuth token not set.
func GetOAuthToken() (string, error) {
	val, ok := os.LookupEnv(tokenEnv)
	if !ok {
		return "", fmt.Errorf(tokenNotFoundErr)
	}
	return val, nil
}
