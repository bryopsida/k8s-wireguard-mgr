package kubernetes

import (
	"context"
	"errors"
	"log"
	"os"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getNamespace() string {
	//namespace will be in /var/run/secrets/kubernetes.io/serviceaccount/namespace
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		panic(err.Error())
	}
	return string(data)
}

type apiStatus interface {
	Status() metav1.Status
}

func reasonForError(err error) metav1.StatusReason {
	if status, ok := err.(apiStatus); ok || errors.As(err, &status) {
		return status.Status().Reason
	}
	return metav1.StatusReasonUnknown
}

func kubernetesErrorIsAlreadyExists(err error) bool {
	return reasonForError(err) == metav1.StatusReasonAlreadyExists
}

func GetClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func CreateWireguardServerSecret(clientset *kubernetes.Clientset, secretName string, privateKey wgtypes.Key) {
	createSecret(clientset, secretName, map[string]string{
		"privatekey": privateKey.String(),
	})
}

func CreateWireguardServerPublicKey(clientset *kubernetes.Clientset, name string, publicKey wgtypes.Key, objectType string) {
	data := map[string]string{
		"publickey": publicKey.String(),
	}
	if objectType == "secret" {
		createSecret(clientset, name, data)
	} else if objectType == "configmap" {
		createConfigMap(clientset, name, data)
	} else {
		panic(errors.New("object type must be secret or configmap"))
	}
}

func createSecret(clientset *kubernetes.Clientset, secretName string, data map[string]string) {
	// try and create the secret, if it already exists we will get an error
	// check the error result, if it's the error about it already existing
	// exit with status 0, this is done because we do not want to request to read access to the secret in our role
	// binding but it is ok to request the ability to create one as this can't be used to read existing secret values
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: getNamespace(),
		},
		StringData: data,
		Type:       corev1.SecretTypeOpaque,
	}
	_, err := clientset.CoreV1().Secrets(getNamespace()).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		if !kubernetesErrorIsAlreadyExists(err) {
			panic(err.Error())
		} else {
			log.Printf("Secret %s already exists. Skipping creation.", secretName)
		}
	} else {
		log.Printf("Secret %s created successfully.", secretName)
	}
}

func createConfigMap(clientset *kubernetes.Clientset, configMapName string, data map[string]string) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: getNamespace(),
		},
		Data: data,
	}
	_, err := clientset.CoreV1().ConfigMaps(getNamespace()).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		if !kubernetesErrorIsAlreadyExists(err) {
			panic(err.Error())
		} else {
			log.Printf("ConfigMap %s already exists. Skipping creation.", configMapName)
		}
	} else {
		log.Printf("ConfigMap %s created successfully.", configMapName)
	}
}
