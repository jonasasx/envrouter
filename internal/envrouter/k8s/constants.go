package k8s

const (
	RepositoryConfigMapName         = "envrouter-repositories"
	ApplicationConfigMapName        = "envrouter-applications"
	RefBindingConfigMapName         = "envrouter-ref-binding"
	SecretTypeLabelKey              = "envrouter.io/secret-type"
	CredentialsSecretTypeLabelValue = "credentials-secret"
	ApplicationLabelKey             = "envrouter.io/app"
	EnvironmentLabelKey             = "envrouter.io/env"
	RefAnnotationKey                = "envrouter.io/ref"
	ShaAnnotationKey                = "envrouter.io/sha"
)
