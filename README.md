# Env Router

Env Router - Continuous Delivery Orchestrator for Kubernetes

## Docker

    docker run -it -v ~/.kube/config:/home/envrouter/.kube/config jonasasx/envrouter:latest

## Helm install

    helm repo add envrouter https://jonasasx.gitlab.io/envrouter/charts
    helm install my-release envrouter/envrouter

## Usage

Set label `envrouter.io/app=<your application name>` to `deployment.metadata.labels` and
`deployment.spec.template.metadata.labels`. Env Router watches for such deployments and pods to display
at the dashboard.