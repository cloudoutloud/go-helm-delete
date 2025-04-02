package main

import (
	"fmt"
	"log"
	"os"
	"time"

	action "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

func main() {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Fatalf("Failed to initialize Helm configuration: %v", err)
	}

	namespace := os.Getenv("HELM_NAMESPACE")
	if namespace == "" {
		log.Fatalf("HELM_NAMESPACE environment variable is not set")
	}

	settings.SetNamespace(namespace)

	listAction := action.NewList(actionConfig)
	listAction.AllNamespaces = false // Limit to a specific namespace
	listAction.StateMask = action.ListAll

	releases, err := listAction.Run()
	if err != nil {
		log.Fatalf("Failed to list Helm releases: %v", err)
	}

	now := time.Now()
	threshold := now.Add(-6 * time.Hour) // 6 hours threshold
	deleteAction := action.NewUninstall(actionConfig)

	for _, r := range releases {
		if r.Info.LastDeployed.Time.Before(threshold) {
			fmt.Printf("Deleting release %s in namespace %s (Deployed: %v)\n", r.Name, r.Namespace, r.Info.LastDeployed)
			_, err := deleteAction.Run(r.Name)
			if err != nil {
				log.Printf("Failed to delete release %s: %v", r.Name, err)
			} else {
				fmt.Printf("Successfully deleted release %s\n", r.Name)
			}
		} else {
			fmt.Printf("No releases to delete")
		}
	}
	fmt.Printf("Script Done\n")
}
