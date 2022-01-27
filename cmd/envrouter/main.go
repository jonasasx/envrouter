package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/jonasasx/envrouter/internal/envrouter"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
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

	podService, stop := k8s.NewPodService(context.TODO(), client)
	defer close(stop)

	applicationService := envrouter.NewApplicationService(deploymentService, dataStorageFactory.NewApplicationStorage(), repositoryService)

	environmentService := envrouter.NewEnvironmentService(deploymentService)

	instanceService := envrouter.NewInstanceService(deploymentService)

	instancePodService := envrouter.NewInstancePodService(podService)

	refService := envrouter.NewRefService(dataStorageFactory.NewRefBindingStorage(), environmentService, applicationService)

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	api.RegisterHandlers(router, &ServerInterfaceImpl{
		repositoryService,
		credentialsSecretService,
		applicationService,
		environmentService,
		instanceService,
		instancePodService,
		refService,
	})

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

func (s *ServerInterfaceImpl) GetApiV1RefBindings(c *gin.Context) {
	result, err := s.refService.FindAllBindings()
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
