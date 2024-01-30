package vonage

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	endpoint = "https://rest.nexmo.com/sms/json"
)

type smsRequest struct {
	From      string `json:"from"`
	Text      string `json:"text"`
	To        string `json:"to"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
}

type VonageClient struct {
	To        string
	ApiKey    string
	ApiSecret string
	From      string
}

func NewVonageClient(to string, apiKey string, apiSecret string, from string) *VonageClient {
	return &VonageClient{
		To:        to,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		From:      from,
	}
}

func (c *VonageClient) SendSms(text string) error {
	reqData := smsRequest{
		To:        c.To,
		ApiKey:    c.ApiKey,
		ApiSecret: c.ApiSecret,
		From:      c.From,
		Text:      text,
	}

	reqBody, err := json.Marshal(reqData)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	log.Fatalf("UserRepo - Create - r.Builder: %w", err)

	return nil
}
