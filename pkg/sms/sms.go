package sms

import "omni-learn-hub/config"

type SMSClient interface {
	SendSMS(message string, to string) error
	GetTemplates() config.Templates
}

type SMSService struct {
	Client SMSClient
}

func NewSmsService(client SMSClient) *SMSService {
	return &SMSService{
		Client: client,
	}
}

func (s *SMSService) GetTemplates() config.Templates {
	return s.Client.GetTemplates()
}

func (s *SMSService) SendSMS(message string, to string) error {
	return s.Client.SendSMS(message, to)

}
