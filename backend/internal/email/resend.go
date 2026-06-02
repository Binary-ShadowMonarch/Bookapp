package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const resendEndpoint = "https://api.resend.com/emails"

type ResendSender struct {
	apiKey string
	from   string
	client *http.Client
}

func NewResendSender(apiKey, from string) (*ResendSender, error) {
	apiKey = strings.TrimSpace(apiKey)
	from = strings.TrimSpace(from)
	if apiKey == "" {
		return nil, fmt.Errorf("resend api key is required")
	}
	if from == "" {
		return nil, fmt.Errorf("resend from email is required")
	}

	return &ResendSender{
		apiKey: apiKey,
		from:   from,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (s *ResendSender) SendVerificationCode(to, code string) error {
	payload := map[string]interface{}{
		"from":    s.from,
		"to":      []string{to},
		"subject": "Your verification code",
		"text":    fmt.Sprintf("Your code is %s", code),
		"html":    fmt.Sprintf("<p>Your verification code is <b>%s</b></p>", code),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("resend send failed: %s", readResponseBody(resp.Body, resp.Status))
	}
	return nil
}

func readResponseBody(body io.Reader, status string) string {
	const maxSize = 4096
	buf := make([]byte, maxSize)
	n, _ := body.Read(buf)
	payload := strings.TrimSpace(string(buf[:n]))
	if payload == "" {
		return status
	}
	return status + ": " + payload
}
