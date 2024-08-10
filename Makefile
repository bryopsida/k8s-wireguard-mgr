READY_STATE ?= true

image:
	docker build . -t ghcr.io/bryopsida/k8s-wireguard-mgr:local

run:
	docker run ghcr.io/bryopsida/k8s-wireguard-mgr:local

format:
	gofmt main.go

test:
    