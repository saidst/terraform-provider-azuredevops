// +build all resource_serviceendpoint_kubernetes

package azuredevops

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

var kubernetesTestServiceEndpointID = uuid.New()
var kubernetesRandomServiceEndpointProjectID = uuid.New().String()
var kubernetesTestServiceEndpointProjectID = &kubernetesRandomServiceEndpointProjectID

var kubernetesTestServiceEndpointForAzureSubscription = serviceendpoint.ServiceEndpoint{ //todo change
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"azureEnvironment": "AzureCloud",
			"azureTenantId":    "kubernetes_TEST_tenant_id",
		},
		Scheme: converter.String("Kubernetes"),
	},
	Data: &map[string]string{
		"authorizationType":     "AzureSubscription",
		"azureSubscriptionId":   "kubernetes_TEST_subscription_id",
		"azureSubscriptionName": "kubernetes_TEST_subscription_name",
		"clusterId":             "/subscriptions/kubernetes_TEST_subscription_id/resourcegroups/kubernetes_TEST_resource_group_id/providers/Microsoft.ContainerService/managedClusters/kubernetes_TEST_cluster_name",
		"namespace":             "default",
	},
	Id:          &kubernetesTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("kubernetes"),
	Url:         converter.String("https://kubernetes.apiserver.com/"),
	Description: converter.String("description"),
}

var kubernetesTestServiceEndpointForKubeconfig = serviceendpoint.ServiceEndpoint{ //todo change
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"clusterContext": "kubernetes_TEST_cluster_context",
			"kubeconfig":     "kubernetes_TEST_tenant_id",
		},
		Scheme: converter.String("Kubernetes"),
	},
	Data: &map[string]string{
		"authorizationType":    "Kubeconfig",
		"acceptUntrustedCerts": "true",
	},
	Id:          &kubernetesTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("kubernetes"),
	Url:         converter.String("https://kubernetes.apiserver.com/"),
	Description: converter.String("description"),
}

var kubernetesTestServiceEndpointForServiceAccount = serviceendpoint.ServiceEndpoint{ //todo change
	Authorization: &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apiToken":                  "kubernetes_TEST_api_token",
			"serviceAccountCertificate": "kubernetes_TEST_ca_cert",
		},
		Scheme: converter.String("Token"),
	},
	Data: &map[string]string{
		"authorizationType": "ServiceAccount",
	},
	Id:          &kubernetesTestServiceEndpointID,
	Name:        converter.String("UNIT_TEST_CONN_NAME"),
	Owner:       converter.String("library"), // Supported values are "library", "agentcloud"
	Type:        converter.String("kubernetes"),
	Url:         converter.String("https://kubernetes.apiserver.com/"),
	Description: converter.String("description"),
}

/**
 * Begin unit tests
 */

// verifies that the flatten/expand round trip yields the same service endpoint for autorization type "AzureSubscription"
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscription_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpointKubernetes().Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointKubernetes(resourceData)

	require.Nil(t, err)
	require.Equal(t, kubernetesTestServiceEndpointForAzureSubscription, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscription_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &kubernetesTestServiceEndpointForAzureSubscription, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscription_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscription_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscription_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &kubernetesTestServiceEndpointForAzureSubscription,
		EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id,
		Project:    kubernetesTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

// verifies that the flatten/expand round trip yields the same service endpoint for autorization type "Kubeconfig"
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfig_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpointKubernetes().Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointKubernetes(resourceData)
	require.Nil(t, err)
	require.Equal(t, kubernetesTestServiceEndpointForKubeconfig, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfig_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &kubernetesTestServiceEndpointForKubeconfig, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfig_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfig_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfig_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &kubernetesTestServiceEndpointForKubeconfig,
		EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id,
		Project:    kubernetesTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

// verifies that the flatten/expand round trip yields the same service endpoint for autorization type "ServiceAccount"
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccount_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpointKubernetes().Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointKubernetes(resourceData)

	require.Nil(t, err)
	require.Equal(t, kubernetesTestServiceEndpointForServiceAccount, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccount_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: &kubernetesTestServiceEndpointForServiceAccount, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateServiceEndpoint() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateServiceEndpoint() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccount_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccount_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccount_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenServiceEndpointKubernetes(resourceData, &kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   &kubernetesTestServiceEndpointForServiceAccount,
		EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id,
		Project:    kubernetesTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}

/**
 * Begin acceptance tests
 */

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsServiceEndpointKubernetesForAzureSubscription_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_kubernetes.serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointKubernetesCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "authorizationType"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureEnvironment"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureTenantId"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureSubscriptionId"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureSubscriptionName"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "clusterId"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "namespace"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameFirst),
				),
			}, {
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "authorizationType"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureEnvironment"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureTenantId"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureSubscriptionId"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azureSubscriptionName"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "clusterId"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "namespace"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameSecond),
				),
			},
		},
	})
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsServiceEndpointKubernetesForServiceAccount_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_kubernetes.serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointKubernetesCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "authorizationType"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "clusterContext"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "kubeconfig"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "acceptUntrustedCerts"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameFirst),
				),
			}, {
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "authorizationType"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "clusterContext"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "kubeconfig"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "acceptUntrustedCerts"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameSecond),
				),
			},
		},
	})
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsServiceEndpointKubernetesForKubeconfig_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := "azuredevops_serviceendpoint_kubernetes.serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointKubernetesCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "authorizationType"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "apiToken"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "serviceAccountCertificate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameFirst),
				),
			}, {
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "authorizationType"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "clusterContext"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "apiToken"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "serviceAccountCertificate"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameSecond),
				),
			},
		},
	})
}

// Given the name of an AzDO service endpoint, this will return a function that will check whether
// or not the resource (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func testAccCheckServiceEndpointKubernetesResourceExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceEndpointDef, ok := s.RootModule().Resources["azuredevops_serviceendpoint_kubernetes.serviceendpoint"]
		if !ok {
			return fmt.Errorf("Did not find a service endpoint in the TF state")
		}

		serviceEndpoint, err := getServiceEndpointKubernetesFromResource(serviceEndpointDef)
		if err != nil {
			return err
		}

		if *serviceEndpoint.Name != expectedName {
			return fmt.Errorf("Service Endpoint has Name=%s, but expected Name=%s", *serviceEndpoint.Name, expectedName)
		}

		return nil
	}
}

// verifies that all service endpoints referenced in the state are destroyed. This will be invoked
// *after* terrafform destroys the resource but *before* the state is wiped clean.
func testAccServiceEndpointKubernetesCheckDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_serviceendpoint_kubernetes" {
			continue
		}

		// indicates the service endpoint still exists - this should fail the test
		if _, err := getServiceEndpointKubernetesFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a service endpoint that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a service endpoint (and error)
func getServiceEndpointKubernetesFromResource(resource *terraform.ResourceState) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpointDefID, err := uuid.Parse(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testAccProvider.Meta().(*config.AggregatedClient)
	return clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{
		Project:    &projectID,
		EndpointId: &serviceEndpointDefID,
	})
}

func init() {
	InitProvider()
}
