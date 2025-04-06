package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	} else {
		fmt.Printf("Using namespace: %s\n", namespace)
	}

	thresholdEnv := os.Getenv("THRESHOLD_HOURS")
	thresholdHours := 24 // Default to 24 hour
	if thresholdEnv != "" {
		parsedThreshold, err := strconv.Atoi(thresholdEnv)
		if err != nil {
			log.Fatalf("Invalid value for THRESHOLD_HOURS: %v", err)
		}
		thresholdHours = parsedThreshold
	}

	exemptReleasesEnv := os.Getenv("EXEMPT_RELEASES")
	var exemptReleases []string
	if exemptReleasesEnv != "" {
		exemptReleases = strings.Split(exemptReleasesEnv, ",")
		fmt.Printf("Exempt releases: %v\n", exemptReleases)
	} else {
		fmt.Println("No exempt releases specified in EXEMPT_RELEASES")
	}

	settings.SetNamespace(namespace)

	listAction := action.NewList(actionConfig)
	listAction.AllNamespaces = false
	listAction.StateMask = action.ListAll

	releases, err := listAction.Run()
	if err != nil {
		log.Fatalf("Failed to list Helm releases: %v", err)
	}

	now := time.Now()
	threshold := now.Add(-time.Duration(thresholdHours) * time.Hour)
	deleteAction := action.NewUninstall(actionConfig)

	for _, r := range releases {
		if isExempted(r.Name, exemptReleases) {
			fmt.Printf("Skipping exempted release %s\n", r.Name)
			continue
		}

		if r.Info.LastDeployed.Time.Before(threshold) {
			fmt.Printf("Deleting release %s in namespace %s (Deployed: %v)\n", r.Name, r.Namespace, r.Info.LastDeployed)
			_, err := deleteAction.Run(r.Name)
			if err != nil {
				log.Printf("Failed to delete release %s: %v", r.Name, err)
			} else {
				fmt.Printf("Successfully deleted release %s\n", r.Name)
			}
		} else {
			fmt.Printf("No releases to delete\n")
		}
	}
	fmt.Printf("Finished..\n")
}

// Helper function to check if a release name is in the exempt list
func isExempted(name string, exemptReleases []string) bool {
	for _, exempt := range exemptReleases {
		if name == exempt {
			return true
		}
	}
	return false
}
