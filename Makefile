CONTEXT_NAME ?= kind

image:
	docker build . -t ghcr.io/bryopsida/k8s-wireguard-mgr:local

run:
	docker run ghcr.io/bryopsida/k8s-wireguard-mgr:local

format:
	gofmt main.go

cluster:
	kind create cluster

cluster-go-away:
	kind delete cluster

build:
	skaffold config set --global collect-metrics false
	skaffold build -q > build_result.json

deploy: build
	skaffold deploy --load-images=true -a build_result.json

verify:
	skaffold verify

clean-logs:
	rm -fr /tmp/kind-logs

logs: clean-logs
	mkdir -p /tmp/kind-logs
	kind export logs /tmp/kind-logs --name=$(CONTEXT_NAME)

test: deploy verify
