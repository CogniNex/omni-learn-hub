package vonage

import (
	"bytes"
	"encoding/json"
	"net/http"
	"omni-learn-hub/config"
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
	ApiKey    string
	ApiSecret string
	From      string
	Templates config.Templates
}

func NewVonageClient(apiKey string, apiSecret string, from string, templates config.Templates) *VonageClient {
	return &VonageClient{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		From:      from,
		Templates: templates,
	}
}

func (c *VonageClient) GetTemplates() config.Templates {
	return c.Templates
}

func (c *VonageClient) SendSMS(text string, to string) error {
	reqData := smsRequest{
		To:        to,
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

	return nil
}
