READY_STATE ?= true

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
	skaffold verify -a build_result.json

test: cluster-go-away cluster deploy verify
	