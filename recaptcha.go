package gorecaptcha

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const verifyURL = "https://www.google.com/recaptcha/api/siteverify"

type Recaptcha struct {
	key        string
	httpClient *http.Client
}

type Response struct {
	Success     bool        `json:"success"`
	Score       float64     `json:"score"`
	Action      string      `json:"action"`
	ChallengeTS time.Time   `json:"challenge_ts"`
	Hostname    string      `json:"hostname"`
	ErrorCodes  []ErrorCode `json:"error-codes"`
}

// ErrorCode is a Recaptcha error code. Full list of codes can be found here:
// https://developers.google.com/recaptcha/docs/verify#error_code_reference
type ErrorCode string

const (
	ErrorCodeMissingInputSecret   ErrorCode = "missing-input-secret"
	ErrorCodeMissingInputResponse ErrorCode = "missing-input-response"
	ErrorCodeInvalidInputResponse ErrorCode = "invalid-input-response"
	ErrorCodeBadRequest           ErrorCode = "bad-request"
	ErrorCodeTimeoutOrDuplicate   ErrorCode = "timeout-or-duplicate"
)

// New creates a new Recaptcha instance.
func New(key string) Recaptcha {
	return Recaptcha{key: key, httpClient: http.DefaultClient}
}

// WithHTTPClient returns a copy of Recaptcha instance with custom HTTP client.
func (r Recaptcha) WithHTTPClient(client *http.Client) Recaptcha {
	return Recaptcha{key: r.key, httpClient: client}
}

// Verify verifies the token on recaptcha server.
func (r Recaptcha) Verify(ip string, token string) (Response, error) {
	return r.VerifyContext(context.Background(), ip, token)
}

// VerifyContext context version of Verify.
func (r Recaptcha) VerifyContext(ctx context.Context, ip string, token string) (Response, error) {
	var ret Response

	bodyData := url.Values{
		"secret":   {r.key},
		"remoteip": {ip},
		"response": {token},
	}

	body := strings.NewReader(bodyData.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, body)
	if err != nil {
		return ret, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return ret, err
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(respBytes, &ret)
	if err != nil {
		return ret, err
	}

	return ret, err
}
