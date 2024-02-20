package wapico

import (
	"bytes"
	"net/http"
	"net/url"
)

const endpoint = "https://biz.wapico.ru/api/"

type messageRequest struct {
}

type WapicoClient struct {
	InstanceId  string
	AccessToken string
}

func NewWapicoClient(instanceId string, accessToken string) *WapicoClient {
	return &WapicoClient{
		InstanceId:  instanceId,
		AccessToken: accessToken,
	}
}

func ConstructURLWithParams(baseURL string, params map[string]string) string {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (w *WapicoClient) SendMessage(text string, to string) error {

	// Parameters
	params := map[string]string{
		"number":       to,
		"type":         "text",
		"message":      text,
		"instance_id":  w.InstanceId,
		"access_token": w.AccessToken,
		// Add more parameters if needed
	}

	// Construct URL with parameters
	fullURL := ConstructURLWithParams(endpoint+"send.php", params)

	req, err := http.NewRequest(http.MethodPost, fullURL, nil)
	if err != nil {
		return err
	}

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
