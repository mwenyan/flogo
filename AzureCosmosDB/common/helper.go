package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

//GetHost returns host url
func GetHost(account string) string {
	return fmt.Sprintf("https://%s.documents.azure.com:443", account)
}

//GetAuthToken returns Authorizaiton access token
func GetAuthToken(key, date, verb, resType, resID string) string {
	msg := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		verb,
		resType,
		resID,
		strings.ToLower(date),
		"")

	hmacKey, _ := base64.StdEncoding.DecodeString(key)

	// handle error
	hasher := hmac.New(sha256.New, hmacKey)
	hasher.Write([]byte(msg))
	signature := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	authHeader := fmt.Sprintf("type=master&ver=1.0&sig=%s", signature)
	return url.QueryEscape(authHeader)
}
