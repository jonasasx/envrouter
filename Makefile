openapi:
	mkdir -p internal/envrouter/api
	oapi-codegen --config api/openapi-spec/gin.config.yaml api/openapi-spec/openapi.yaml
	oapi-codegen --config api/openapi-spec/types.config.yaml api/openapi-spec/openapi.yaml
	oapi-codegen --config api/openapi-spec/spec.config.yaml api/openapi-spec/openapi.yaml
	openapi-generator-cli generate -i api/openapi-spec/openapi.yaml -g typescript-axios -o web/src/axios

deps:
	go mod download

build:
	GOOS=linux GOARCH=386 go build -o build/envrouter cmd/envrouter/main.go

web:
	unset CI
	cd web
	npm i
	npm run build

run:
	go run cmd/envrouter/main.go