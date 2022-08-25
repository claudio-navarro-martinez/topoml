package kubernetestool

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func mainkubernetestool() {
	// create 
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
    // get the pod list
    podList, _ := kubeClient.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})

       // List() returns a pointer to slice, derefernce it, before iterating
	for _, podInfo := range (*podList).Items {
		fmt.Printf("pods-name=%v\n", podInfo.Name)
		fmt.Printf("pods-status=%v\n", podInfo.Status.Phase)
		fmt.Printf("pods-condition=%v\n", podInfo.Status.Conditions)
	}
	// deploy a Deployment
	deploymentsClient := kubeClient.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var creado bool = false
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
		if d.Name == "demo-deployment" {
			creado = true
		}
	}
	if !creado {
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo-deployment",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "demo",
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "demo",
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "web",
								Image: "nginx:1.12",
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		}
		result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

		if err != nil {
			panic(err)
		}
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
		// sacar output del pod
	}
}

func int32Ptr(i int32) *int32 { return &i }
