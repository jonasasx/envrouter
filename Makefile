openapi:
	mkdir -p internal/envrouter/api
	oapi-codegen --config api/openapi-spec/gin.config.yaml api/openapi-spec/openapi.yaml
	oapi-codegen --config api/openapi-spec/types.config.yaml api/openapi-spec/openapi.yaml
	oapi-codegen --config api/openapi-spec/spec.config.yaml api/openapi-spec/openapi.yaml
	openapi-generator-cli generate -i api/openapi-spec/openapi.yaml -g typescript-axios -o web/src/axios

build:
	GOOS=linux GOARCH=386 go build -o build/envrouter cmd/envrouter/main.go

run:
	go run cmd/envrouter/main.go