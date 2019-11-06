package azuredevops

import (
	"log"
	"context"
	"testing"

	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/microsoft/azure-devops-go-api/azuredevops/core"
	// "github.com/microsoft/azure-devops-go-api/azuredevops/operations"
	// "github.com/stretchr/testify/require"
)

func TestAzureDevOpsProject_CreatePolicy_DoesStuff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish();
	
	policyClient := azdosdkmocks.NewMockPolicyClient(ctrl)
	
	clients := &aggregatedClient{
		PolicyClient: policyClient,
		ctx:        context.Background(),
	}

	resource := resourcePolicyMinReviewers()
	config := map[string]interface{}{
		"project_id": "foo",
	}
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, config)
	err := resourcePolicyMinReviewersCreate(resourceData, clients)

	log.Printf("logging foo %s\n", err);
}