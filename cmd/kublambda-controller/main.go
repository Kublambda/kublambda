package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ensureThirdPartyResource(clientset)
	ensureFunctionStore(clientset)
	ensureRunners(clientset)
	getUpToSpeed(clientset)
	watchForChanges(clientset)

}

func ensureFunctionStore(clientset *kubernetes.Clientset) {
	fmt.Printf("Starting the function store...")
}

func getUpToSpeed(clientset *kubernetes.Clientset) {
	fmt.Printf("SYNCHRONIZING LAMBDAS...\n")
	// get the list of current lambdas from api server
	// make sure their functions are registered in object store
	// expose any functions with http triggers defined
}

func watchForChanges(clientset *kubernetes.Clientset) {
	fmt.Printf("BEGINNING EVENT LOOP...\n")
	// on lambda added
	//  - upload code to object store
	//  - if http triggered, expose uri
	// on lambda updated
	//  - replace code in object store
	//  - add/remove exposures
	// on lambda deleted
	//  - remove from object store
	//  - remove exposures
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Tick...")
	}
}

func ensureRunners(clientset *kubernetes.Clientset) {
	runner, err := clientset.ExtensionsV1beta1().Deployments("default").Get("kublambda-runner")
	if err != nil {
		if errors.IsNotFound(err) {

			runner := &v1beta1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name: "kublambda-runner",
				},
				Spec: v1beta1.DeploymentSpec{
					Template: v1.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: map[string]string{"app": "kublambda-runner"},
						},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  "kublambda-runner",
									Image: "forjared/lambda-runner",
									Ports: []v1.ContainerPort{
										v1.ContainerPort{ContainerPort: 8080},
									},
								},
							},
						},
					},
				},
			}

			result, err := clientset.ExtensionsV1beta1().Deployments("default").Create(runner)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("CREATED: %#v\nFROM: %#v\n", result, runner)
		} else {
			panic(err.Error())
		}
	} else {
		fmt.Printf("SKIPPING: already exists %#v\n", runner)
	}
}

// initialize third party resource if it does not exist
func ensureThirdPartyResource(clientset *kubernetes.Clientset) {
	tpr, err := clientset.ExtensionsV1beta1().ThirdPartyResources().Get("kublambda.vmware.com")
	if err != nil {
		if errors.IsNotFound(err) {
			tpr := &v1beta1.ThirdPartyResource{
				ObjectMeta: v1.ObjectMeta{
					Name: "kublambda.vmware.com",
				},
				Versions: []v1beta1.APIVersion{
					{Name: "v1"},
				},
				Description: "A lambda function",
			}

			result, err := clientset.ExtensionsV1beta1().ThirdPartyResources().Create(tpr)
			if err != nil {
				panic(err)
			}
			fmt.Printf("CREATED: %#v\nFROM: %#v\n", result, tpr)
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("SKIPPING: already exists %#v\n", tpr)
	}
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
