package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var getPods = &cobra.Command{
	Use:   "list [namespace_name]",
	Short: "Faz o parse de um arquivo de configuração YAML ou JSON",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespaceName := args[0]
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("error getting user home dir: %v\n", err)
			os.Exit(1)
		}

		kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			fmt.Printf("Error getting kubernetes config: %v\n", err)
			os.Exit(1)
		}

		clientset, err := kubernetes.NewForConfig(kubeConfig)

		if err != nil {
			fmt.Printf("error getting kubernetes config: %v\n", err)
			os.Exit(1)
		}

		namespace := namespaceName
		pods, err := ListPods(namespace, clientset)
		if err != nil {
			fmt.Println(err.Error)
			os.Exit(1)
		}
		for _, pod := range pods.Items {
			fmt.Printf("Pod name: %v\n", pod.Name)
		}
		var message string
		if namespace == "" {
			message = "Total Pods in all namespaces"
		} else {
			message = fmt.Sprintf("Total Pods in namespace `%s`", namespace)
		}
		fmt.Printf("%s %d\n", message, len(pods.Items))
	},
}

func ListPods(namespace string, client kubernetes.Interface) (*v1.PodList, error) {
	fmt.Println("Get Kubernetes Pods")
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v\n", err)
		return nil, err
	}
	return pods, nil
}

func init() {
	containerCmd.AddCommand(getPods)
}
