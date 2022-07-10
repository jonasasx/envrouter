.PHONY: openapi
openapi:
	mkdir -p internal/envrouter/api
	oapi-codegen -generate types -package api api/openapi-spec/openapi.yaml > internal/envrouter/api/types.go
	oapi-codegen -generate gin -package api api/openapi-spec/openapi.yaml > internal/envrouter/api/server.go
	oapi-codegen -generate spec -package api api/openapi-spec/openapi.yaml > internal/envrouter/api/spec.go
	openapi-generator-cli generate -i api/openapi-spec/openapi.yaml -g typescript-axios -o web/src/axios

.PHONY: build
build:
	go build -o build/envrouter cmd/envrouter/main.go

.PHONY: helm-docs
helm-docs:
	scripts/helm-docs.sh