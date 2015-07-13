package signer

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	time "github.com/ssoroka/ttime"
	"net/url"
	"strconv"
)

type Signer struct {
	AppID  string
	Secret string
}

func New(appId string, secret string) *Signer {
	return &Signer{
		AppID:  appId,
		Secret: secret,
	}
}

func (signer *Signer) Sign(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}

	values := parsedURL.Query()
	currentUnixTime := strconv.FormatInt(time.Now().Unix(), 10)
	values.Set("t", currentUnixTime)
	values.Set("app", signer.AppID)
	parsedURL.RawQuery = values.Encode()

	sig := signer.GetSignature(parsedURL.Path + "?" + parsedURL.RawQuery)
	values.Set("sig", sig)
	parsedURL.RawQuery = values.Encode()

	return buildURLString(*parsedURL), nil
}

func (signer *Signer) GetSignature(pathWithQuery string) string {
	secretBytes := []byte(signer.Secret)
	pathWithQueryBytes := []byte(pathWithQuery)
	hasher := hmac.New(sha1.New, secretBytes)
	hasher.Write(pathWithQueryBytes)
	sig := hex.EncodeToString(hasher.Sum(nil))
	return sig
}

func buildURLString(u url.URL) string {
	if u.Scheme != "" && u.Host != "" {
		return u.Scheme + "://" + u.Host + u.Path + "?" + u.RawQuery
	} else {
		return u.Path + "?" + u.RawQuery
	}
}
