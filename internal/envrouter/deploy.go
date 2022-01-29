package envrouter

import (
	"strings"
)

type DeployService interface {
	Deploy(applicationName string, ref string) error
}

type deployService struct {
	applicationService ApplicationService
	webhookService     WebhookService
}

func NewDeployService(
	applicationService ApplicationService,
	webhookService WebhookService,
) DeployService {
	return &deployService{
		applicationService,
		webhookService,
	}
}

func (d *deployService) Deploy(applicationName string, ref string) error {
	application, err := d.applicationService.FindByName(applicationName)
	if err != nil {
		return err
	}
	if application.Webhook != nil {
		webhook := *application.Webhook
		webhook = strings.ReplaceAll(webhook, "{ref}", ref)
		return d.webhookService.Invoke(webhook)
	}
	return nil
}
