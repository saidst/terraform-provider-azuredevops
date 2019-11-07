package azuredevops

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/microsoft/azure-devops-go-api/azuredevops/core"
	// "github.com/microsoft/azure-devops-go-api/azuredevops/operations"
	// "github.com/stretchr/testify/require"
)

func TestAzureDevOpsProject_CreatePolicy_DoesStuff(t *testing.T) {
	adoPat, adoUrl := os.Getenv("AZDO_PERSONAL_ACCESS_TOKEN"), os.Getenv("AZDO_ORG_SERVICE_URL")
	clients, err := getAzdoClient(adoPat, adoUrl)
	project, err := projectRead(clients, "", "az-infrax")

	resource := resourcePolicyMinReviewers()
	config := map[string]interface{}{
		"project_id": project.Id.String(),
	}

	resourceData := schema.TestResourceDataRaw(t, resource.Schema, config)
	resourcePolicyMinReviewersCreate(resourceData, clients)

	log.Printf("logging foo %s\n", err)
}
