package envrouter

import (
	"net/http"
	"net/url"
)

type WebhookService interface {
	Invoke(webhook string) error
}

type webhookService struct {
}

func NewWebhookService() WebhookService {
	return &webhookService{}
}

func (w *webhookService) Invoke(webhook string) error {
	_, err := http.PostForm(webhook, url.Values{})
	if err != nil {
		return err
	}
	return nil
}
