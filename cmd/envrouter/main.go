package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/jonasasx/envrouter/internal/envrouter"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/utils"
	"io"
	"net/http"
)

func init() {
	log.Infof("Init")
}

func main() {
	var err error
	client := k8s.NewClient("")

	dataStorageFactory := k8s.NewDataStorageFactory(client)

	repositoryService := envrouter.NewRepositoryService(dataStorageFactory.NewRepositoryStorage())

	credentialsSecretService := envrouter.NewCredentialsSecretService(dataStorageFactory.NewCredentialsSecretStorage())

	deploymentService, stop := k8s.NewDeploymentService(context.TODO(), client)
	defer close(stop)

	podServiceFactoryMethod := k8s.NewPodServiceFactoryMethod(context.TODO(), client)

	applicationService := envrouter.NewApplicationService(deploymentService, dataStorageFactory.NewApplicationStorage(), repositoryService)

	environmentService := envrouter.NewEnvironmentService(deploymentService)

	instanceService := envrouter.NewInstanceService(deploymentService)

	instancePodObserver := utils.NewObserver()
	instancePodService, stop := envrouter.NewInstancePodService(podServiceFactoryMethod, instancePodObserver)
	defer close(stop)

	webhookService := envrouter.NewWebhookService()

	deployService := envrouter.NewDeployService(applicationService, webhookService)

	refService := envrouter.NewRefService(dataStorageFactory.NewRefBindingStorage(), environmentService, applicationService, deployService)

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	server := &ServerInterfaceImpl{
		repositoryService,
		credentialsSecretService,
		applicationService,
		environmentService,
		instanceService,
		instancePodService,
		refService,
		instancePodObserver,
	}
	router.GET("/api/v1/subscription", server.streamPods)
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	api.RegisterHandlers(router, server)

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}

type ServerInterfaceImpl struct {
	repositoryService        envrouter.RepositoryService
	credentialsSecretService envrouter.CredentialsSecretService
	applicationService       envrouter.ApplicationService
	environmentService       envrouter.EnvironmentService
	instanceService          envrouter.InstanceService
	instancePodService       envrouter.InstancePodService
	refService               envrouter.RefService
	instancePodObserver      utils.Observer
}

func (s *ServerInterfaceImpl) GetApiV1Repositories(c *gin.Context) {
	result, err := s.repositoryService.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) PostApiV1Repositories(c *gin.Context) {
	var json api.Repository
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := s.repositoryService.Save(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) DeleteApiV1RepositoriesName(c *gin.Context, name string) {
	err := s.repositoryService.DeleteByName(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (s *ServerInterfaceImpl) GetApiV1CredentialsSecrets(c *gin.Context) {
	result, err := s.credentialsSecretService.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) PostApiV1CredentialsSecrets(c *gin.Context) {
	var json api.CredentialsSecretRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := s.credentialsSecretService.Save(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) DeleteApiV1CredentialsSecretsName(c *gin.Context, name string) {
	err := s.credentialsSecretService.DeleteByName(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (s *ServerInterfaceImpl) GetApiV1Applications(c *gin.Context) {
	result, err := s.applicationService.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) PutApiV1ApplicationsName(c *gin.Context, name string) {
	var json api.Application
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.Name != name {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Names in path and body are different"})
		return
	}
	result, err := s.applicationService.Save(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) GetApiV1Environments(c *gin.Context) {
	result, err := s.environmentService.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) GetApiV1Instances(c *gin.Context) {
	result, err := s.instanceService.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) GetApiV1InstancePods(c *gin.Context) {
	result, err := s.instancePodService.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) GetApiV1RefBindings(c *gin.Context, params api.GetApiV1RefBindingsParams) {
	result, err := s.refService.FindAllBindings(params.Environment, params.Application, params.Ref)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) PostApiV1RefBindings(c *gin.Context) {
	var json api.RefBinding
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := s.refService.SaveBinding(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(200, result)
	}
}

func (s *ServerInterfaceImpl) streamPods(c *gin.Context) {
	subscriber := make(chan utils.ObserverEvent)
	s.instancePodObserver.Subscribe(&subscriber)
	defer s.instancePodObserver.Unsubscribe(&subscriber)

	c.Stream(func(w io.Writer) bool {
		event := <-subscriber
		c.SSEvent("", event)
		return true
	})
}
