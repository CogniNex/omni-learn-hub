package whatsapp

type WhatsappClient interface {
	SendMessage(message string, to string) error
}

type WhatsappService struct {
	Client WhatsappClient
}

func NewWhatsappService(client WhatsappClient) *WhatsappService {
	return &WhatsappService{
		Client: client,
	}
}

func (s *WhatsappService) SendMessage(message string, to string) error {
	return s.Client.SendMessage(message, to)

}
