package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/docker/docker/api/types/registry"
	_ "github.com/joho/godotenv/autoload"
	v13 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	k8sconfig, err := GetKubeRestConfig()
	if err != nil {
		log.Fatal(err)
	}
	kube, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		log.Fatal(err)
	}
	client := ecr.NewFromConfig(cfg)
	token, err := client.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]registry.AuthConfig)
	//types.AuthConfig{}
	for _, v := range token.AuthorizationData {
		u, err := url.Parse(*v.ProxyEndpoint)
		if err != nil {
			log.Fatal(err)
		}
		data[u.Host] = registry.AuthConfig{
			Auth:          *v.AuthorizationToken,
			ServerAddress: *v.ProxyEndpoint,
		}
	}

	secretName := os.Getenv("SECRET_NAME")
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	namespaces, err := kube.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, namespace := range namespaces.Items {
		namespaceName := namespace.Name
		secrets := kube.CoreV1().Secrets(namespaceName)
		log.Printf("Apply secret %s in %s\n", secretName, namespaceName)
		err = secrets.Delete(context.TODO(), secretName, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		secret := v13.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespaceName,
			},

			StringData: map[string]string{
				v13.DockerConfigKey: string(b),
			},
			Type: v13.SecretTypeDockercfg,
		}
		_, err = secrets.Create(context.TODO(), &secret, metav1.CreateOptions{})
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}

func GetKubeRestConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_CONFIG") == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", os.Getenv("KUBERNETES_CONFIG"))
}
