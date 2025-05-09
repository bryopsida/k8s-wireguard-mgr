# Kubernetes Wireguard Manager

## What is this?

This is a CLI tool intended to be run inside of a kubernetes job that can manage wireguard configurations held inside of kubernetes objects.
It's initial use case is to generate a wireguard server private key as a part of a hook for a helm chart install.

## Why not just use bash and kubectl?

1. I prefer to avoid bash once conditionals start coming into play
2. By using the wireguard library directly and the kubernetes API directly, the key can be generated without having to do any gymnastics keeping it out of CLI args that would show up on the node process table.
3. It can be re-used/expanded to manage kubernetes resources that hold wireguard objects to manage a VPN.

## Environment variables

| Name                              | Description                                                                             |
|-----------------------------------|-----------------------------------------------------------------------------------------|
| K8S_WG_MGR_SERVER_SECRET_NAME     | Name of the secret to create                                                            |
| K8S_WG_MGR_SERVER_PUBLIC_KEY_NAME | Name of the public key object to create, default uses secret name with `-public` suffix |
| K8S_WG_MGR_SERVER_PUBLIC_KEY_TYPE | Type of public key object to create, either `secret` or `configmap`, default `secret`   |
