// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

// Defines values for InstanceType.
const (
	InstanceTypeDeployment InstanceType = "deployment"
)

// Application defines model for Application.
type Application struct {
	Name           string  `json:"name"`
	RepositoryName *string `json:"repositoryName,omitempty"`
}

// CredentialsSecretListItem defines model for CredentialsSecretListItem.
type CredentialsSecretListItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// CredentialsSecretRequest defines model for CredentialsSecretRequest.
type CredentialsSecretRequest struct {
	Key      string `json:"key"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// Environment defines model for Environment.
type Environment struct {
	Name string `json:"name"`
}

// Instance defines model for Instance.
type Instance struct {
	Application string       `json:"application"`
	CommitSha   *string      `json:"commitSha,omitempty"`
	Environment string       `json:"environment"`
	Name        string       `json:"name"`
	Ref         *string      `json:"ref,omitempty"`
	Type        InstanceType `json:"type"`
}

// InstanceType defines model for Instance.Type.
type InstanceType string

// RefBinding defines model for RefBinding.
type RefBinding struct {
	Application string `json:"application"`
	Environment string `json:"environment"`
	Ref         string `json:"ref"`
}

// Repository defines model for Repository.
type Repository struct {
	CredentialsSecret string `json:"credentialsSecret"`
	Name              string `json:"name"`
	Url               string `json:"url"`
}

// PutApiV1ApplicationsNameJSONBody defines parameters for PutApiV1ApplicationsName.
type PutApiV1ApplicationsNameJSONBody Application

// PostApiV1CredentialsSecretsJSONBody defines parameters for PostApiV1CredentialsSecrets.
type PostApiV1CredentialsSecretsJSONBody CredentialsSecretRequest

// PostApiV1RefBindingsJSONBody defines parameters for PostApiV1RefBindings.
type PostApiV1RefBindingsJSONBody RefBinding

// PutApiV1ApplicationsNameJSONRequestBody defines body for PutApiV1ApplicationsName for application/json ContentType.
type PutApiV1ApplicationsNameJSONRequestBody PutApiV1ApplicationsNameJSONBody

// PostApiV1CredentialsSecretsJSONRequestBody defines body for PostApiV1CredentialsSecrets for application/json ContentType.
type PostApiV1CredentialsSecretsJSONRequestBody PostApiV1CredentialsSecretsJSONBody

// PostApiV1RefBindingsJSONRequestBody defines body for PostApiV1RefBindings for application/json ContentType.
type PostApiV1RefBindingsJSONRequestBody PostApiV1RefBindingsJSONBody