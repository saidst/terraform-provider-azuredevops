package azuredevops

import (
	"log"
	"os"
	"testing"
	"fmt"

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
	/*var scopes []map[string]interface{}
	scopes = append(scopes,map[string]interface{}{
		"refName": "refs/heads/master",
		"matchKind": "Exact",
		"repositoryId": nil,
	})*/

	config := map[string]interface{}{
		"project_id": project.Id.String(),
		"scope": []interface{}{
			map[string]interface{}{
				"repository_id": "6c3be0fc-2cb8-465e-9484-d0545fddbfd8",
			},			
		},
	}

	resourceData := schema.TestResourceDataRaw(t, resource.Schema, config)
	resourcePolicyMinReviewersCreate(resourceData, clients)

	log.Printf("logging foo %s\n", err)
}

func testAccPolicyMinReviewersResource(projectID string, repositoryID string) string {
	policyMinReviewersResource := fmt.Sprintf(`
	resource "azuredevops_policy_min_reviews" "policy-review" {
		project_id  = %s
		scope {
			repository_id = %s
		}
	}`, projectID, repositoryID)
	return policyMinReviewersResource
}