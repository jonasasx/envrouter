.PHONY: openapi
openapi:
	mkdir -p internal/envrouter/api
	oapi-codegen --config api/openapi-spec/gin.config.yaml api/openapi-spec/openapi.yaml
	oapi-codegen --config api/openapi-spec/types.config.yaml api/openapi-spec/openapi.yaml
	oapi-codegen --config api/openapi-spec/spec.config.yaml api/openapi-spec/openapi.yaml
	openapi-generator-cli generate -i api/openapi-spec/openapi.yaml -g typescript-axios -o web/src/axios

.PHONY: build
build:
	go build -o build/envrouter cmd/envrouter/main.go

.PHONY: helm-docs
helm-docs:
	scripts/helm-docs.sh