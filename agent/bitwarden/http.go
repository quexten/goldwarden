package bitwarden

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 20 * time.Second,
}

type errStatusCode struct {
	code int
	body []byte
}

func (e *errStatusCode) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.code), e.body)
}

type AuthToken struct{}

func authenticatedHTTPPost(ctx context.Context, urlstr string, recv, send interface{}) error {
	var r io.Reader
	contentType := "application/json"
	authEmail := ""
	if values, ok := send.(url.Values); ok {
		r = strings.NewReader(values.Encode())
		contentType = "application/x-www-form-urlencoded"
		if email := values.Get("username"); email != "" && values.Get("scope") != "" {
			authEmail = email
		}
	} else {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(send); err != nil {
			return err
		}
		r = buf
	}
	req, err := http.NewRequest("POST", urlstr, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	if authEmail != "" {
		req.Header.Set("Auth-Email", base64.URLEncoding.EncodeToString([]byte(authEmail)))
	}
	return makeAuthenticatedHTTPRequest(ctx, req, recv)
}

func authenticatedHTTPGet(ctx context.Context, urlstr string, recv interface{}) error {
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return err
	}
	return makeAuthenticatedHTTPRequest(ctx, req, recv)
}

func authenticatedHTTPDelete(ctx context.Context, urlstr string, recv interface{}) error {
	req, err := http.NewRequest("DELETE", urlstr, nil)
	if err != nil {
		return err
	}
	return makeAuthenticatedHTTPRequest(ctx, req, recv)
}

func authenticatedHTTPPut(ctx context.Context, urlstr string, recv, send interface{}) error {
	var r io.Reader
	contentType := "application/json"
	if values, ok := send.(url.Values); ok {
		r = strings.NewReader(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(send); err != nil {
			return err
		}
		r = buf
	}
	req, err := http.NewRequest("PUT", urlstr, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	return makeAuthenticatedHTTPRequest(ctx, req, recv)
}

func makeAuthenticatedHTTPRequest(ctx context.Context, req *http.Request, recv interface{}) error {
	if token, ok := ctx.Value(AuthToken{}).(string); ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return &errStatusCode{res.StatusCode, body}
	}
	if err := json.Unmarshal(body, recv); err != nil {
		fmt.Fprintln(os.Stderr, string(body))
		return err
	}
	return nil
}
