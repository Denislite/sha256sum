package utils

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"os/signal"
	"sha256sum/internal/model"
)

func CheckSignal(signals chan os.Signal) {
	signal.Notify(signals, os.Interrupt)
	go func() {
		for sig := range signals {
			log.Printf("request canceled by signal %d \n", sig)
			os.Exit(0)
		}
	}()
}

func NewK8SConnection() (*kubernetes.Clientset, *model.ContainerInfo, error) {
	log.Println("### 🚀 K8S checksum starting...")

	log.Println("### 🌀 Attempting to use in cluster config")
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, nil, err
	}

	log.Printf("### 💻 Connecting to Kubernetes API, using host: %s", config.Host)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	deploymentList, err := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE")).Get(context.Background(),
		os.Getenv("DEPLOYMENT_NAME"), metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	container := &model.ContainerInfo{
		PodName:      os.Getenv("POD_NAME"),
		ImageName:    deploymentList.Spec.Template.Spec.Containers[0].Name,
		ImageVersion: deploymentList.Spec.Template.Spec.Containers[0].Image,
	}

	return clientset, container, nil
}
