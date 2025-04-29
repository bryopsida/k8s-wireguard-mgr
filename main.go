package main

import (
	"log"
	"os"

	"github.com/bryopsida/k8s-wireguard-mgr/kubernetes"
	"github.com/bryopsida/k8s-wireguard-mgr/wireguard"
)

func main() {
	log.Println("Starting Kubernetes Wireguard Manager")
	key := wireguard.GenerateWireguardKey()
	secretName := os.Getenv("K8S_WG_MGR_SERVER_SECRET_NAME")
	clientset, err := kubernetes.GetClientSet()
	if err != nil {
		panic(err.Error())
	}
	kubernetes.CreateWireguardServerSecret(clientset, secretName, key)
	publicKeyName := os.Getenv("K8S_WG_MGR_SERVER_PUBLIC_KEY_NAME")
	if publicKeyName == "" {
		publicKeyName = secretName + "-public"
	}
	publicKeyType := os.Getenv("K8S_WG_MGR_SERVER_PUBLIC_KEY_TYPE")
	if publicKeyType == "" {
		publicKeyType = "secret"
	}
	kubernetes.CreateWireguardServerPublicKey(clientset, publicKeyName, key.PublicKey(), publicKeyType)
	log.Println("Finished")
}
