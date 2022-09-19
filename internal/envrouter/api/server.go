// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/v1/applications)
	GetApiV1Applications(c *gin.Context)

	// (PUT /api/v1/applications/{name})
	PutApiV1ApplicationsName(c *gin.Context, name string)

	// (GET /api/v1/credentialsSecrets)
	GetApiV1CredentialsSecrets(c *gin.Context)

	// (POST /api/v1/credentialsSecrets)
	PostApiV1CredentialsSecrets(c *gin.Context)

	// (DELETE /api/v1/credentialsSecrets/{name})
	DeleteApiV1CredentialsSecretsName(c *gin.Context, name string)

	// (GET /api/v1/environments)
	GetApiV1Environments(c *gin.Context)

	// (GET /api/v1/git/refs)
	GetApiV1GitRefs(c *gin.Context)

	// (GET /api/v1/git/repositories/{repositoryName}/commits/{sha})
	GetApiV1GitRepositoriesRepositoryNameCommitsSha(c *gin.Context, repositoryName string, sha string)

	// (GET /api/v1/instancePods)
	GetApiV1InstancePods(c *gin.Context)

	// (GET /api/v1/instances)
	GetApiV1Instances(c *gin.Context)

	// (GET /api/v1/refBindings)
	GetApiV1RefBindings(c *gin.Context, params GetApiV1RefBindingsParams)

	// (POST /api/v1/refBindings)
	PostApiV1RefBindings(c *gin.Context)
	// Get all repositories
	// (GET /api/v1/repositories)
	GetApiV1Repositories(c *gin.Context)

	// (POST /api/v1/repositories)
	PostApiV1Repositories(c *gin.Context)

	// (DELETE /api/v1/repositories/{name})
	DeleteApiV1RepositoriesName(c *gin.Context, name string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(c *gin.Context)

// GetApiV1Applications operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1Applications(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1Applications(c)
}

// PutApiV1ApplicationsName operation middleware
func (siw *ServerInterfaceWrapper) PutApiV1ApplicationsName(c *gin.Context) {

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", c.Param("name"), &name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter name: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutApiV1ApplicationsName(c, name)
}

// GetApiV1CredentialsSecrets operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1CredentialsSecrets(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1CredentialsSecrets(c)
}

// PostApiV1CredentialsSecrets operation middleware
func (siw *ServerInterfaceWrapper) PostApiV1CredentialsSecrets(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostApiV1CredentialsSecrets(c)
}

// DeleteApiV1CredentialsSecretsName operation middleware
func (siw *ServerInterfaceWrapper) DeleteApiV1CredentialsSecretsName(c *gin.Context) {

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", c.Param("name"), &name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter name: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteApiV1CredentialsSecretsName(c, name)
}

// GetApiV1Environments operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1Environments(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1Environments(c)
}

// GetApiV1GitRefs operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1GitRefs(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1GitRefs(c)
}

// GetApiV1GitRepositoriesRepositoryNameCommitsSha operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1GitRepositoriesRepositoryNameCommitsSha(c *gin.Context) {

	var err error

	// ------------- Path parameter "repositoryName" -------------
	var repositoryName string

	err = runtime.BindStyledParameter("simple", false, "repositoryName", c.Param("repositoryName"), &repositoryName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter repositoryName: %s", err)})
		return
	}

	// ------------- Path parameter "sha" -------------
	var sha string

	err = runtime.BindStyledParameter("simple", false, "sha", c.Param("sha"), &sha)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter sha: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1GitRepositoriesRepositoryNameCommitsSha(c, repositoryName, sha)
}

// GetApiV1InstancePods operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1InstancePods(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1InstancePods(c)
}

// GetApiV1Instances operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1Instances(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1Instances(c)
}

// GetApiV1RefBindings operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1RefBindings(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetApiV1RefBindingsParams

	// ------------- Optional query parameter "application" -------------
	if paramValue := c.Query("application"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "application", c.Request.URL.Query(), &params.Application)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter application: %s", err)})
		return
	}

	// ------------- Optional query parameter "environment" -------------
	if paramValue := c.Query("environment"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "environment", c.Request.URL.Query(), &params.Environment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter environment: %s", err)})
		return
	}

	// ------------- Optional query parameter "ref" -------------
	if paramValue := c.Query("ref"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "ref", c.Request.URL.Query(), &params.Ref)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter ref: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1RefBindings(c, params)
}

// PostApiV1RefBindings operation middleware
func (siw *ServerInterfaceWrapper) PostApiV1RefBindings(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostApiV1RefBindings(c)
}

// GetApiV1Repositories operation middleware
func (siw *ServerInterfaceWrapper) GetApiV1Repositories(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiV1Repositories(c)
}

// PostApiV1Repositories operation middleware
func (siw *ServerInterfaceWrapper) PostApiV1Repositories(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostApiV1Repositories(c)
}

// DeleteApiV1RepositoriesName operation middleware
func (siw *ServerInterfaceWrapper) DeleteApiV1RepositoriesName(c *gin.Context) {

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", c.Param("name"), &name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter name: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteApiV1RepositoriesName(c, name)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	router.GET(options.BaseURL+"/api/v1/applications", wrapper.GetApiV1Applications)

	router.PUT(options.BaseURL+"/api/v1/applications/:name", wrapper.PutApiV1ApplicationsName)

	router.GET(options.BaseURL+"/api/v1/credentialsSecrets", wrapper.GetApiV1CredentialsSecrets)

	router.POST(options.BaseURL+"/api/v1/credentialsSecrets", wrapper.PostApiV1CredentialsSecrets)

	router.DELETE(options.BaseURL+"/api/v1/credentialsSecrets/:name", wrapper.DeleteApiV1CredentialsSecretsName)

	router.GET(options.BaseURL+"/api/v1/environments", wrapper.GetApiV1Environments)

	router.GET(options.BaseURL+"/api/v1/git/refs", wrapper.GetApiV1GitRefs)

	router.GET(options.BaseURL+"/api/v1/git/repositories/:repositoryName/commits/:sha", wrapper.GetApiV1GitRepositoriesRepositoryNameCommitsSha)

	router.GET(options.BaseURL+"/api/v1/instancePods", wrapper.GetApiV1InstancePods)

	router.GET(options.BaseURL+"/api/v1/instances", wrapper.GetApiV1Instances)

	router.GET(options.BaseURL+"/api/v1/refBindings", wrapper.GetApiV1RefBindings)

	router.POST(options.BaseURL+"/api/v1/refBindings", wrapper.PostApiV1RefBindings)

	router.GET(options.BaseURL+"/api/v1/repositories", wrapper.GetApiV1Repositories)

	router.POST(options.BaseURL+"/api/v1/repositories", wrapper.PostApiV1Repositories)

	router.DELETE(options.BaseURL+"/api/v1/repositories/:name", wrapper.DeleteApiV1RepositoriesName)

	return router
}
