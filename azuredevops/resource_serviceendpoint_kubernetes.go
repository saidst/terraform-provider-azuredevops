package azuredevops

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	crud "github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/serviceendpoint"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func makeSchemaAzureSubscription(r *schema.Resource) {
	r.Schema["azure_subscription"] = &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "'AzureSubscription'-type of configuration",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"azure_environment": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "AzureCloud",
					Description:  "type of azure cloud: AzureCloud",
					ValidateFunc: validation.StringInSlice([]string{"AzureCloud"}, false),
				},
				"cluster_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "name of aks-resource",
				},
				"subscription_id": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "id of azure subscription",
					ValidateFunc: validation.IsUUID,
				},
				"subscription_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "name of azure subscription",
				},
				"tenant_id": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "id of aad-tenant",
					ValidateFunc: validation.IsUUID,
				},
				"resourcegroup_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "id of resourcegroup",
				},
				"namespace": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "default",
					Description: "accessed namespace",
				},
			},
		},
	}
}

func makeSchemaKubeconfig(r *schema.Resource) {
	r.Schema["kubeconfig"] = &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "'Kubeconfig'-type of configuration",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"kube_config": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Content of the kubeconfig file. The configuration information in your kubeconfig file allows Kubernetes clients to talk to your Kubernetes API servers. This file is used by kubectl and all supported Kubernetes clients.",
				},
				"cluster_context": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Context of your cluster",
				},
				"accept_untrusted_certs": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Enable this if your authentication uses untrusted certificates",
				},
			},
		},
	}
}

func makeSchemaServiceAccount(r *schema.Resource) {
	r.Schema["service_account"] = &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "'ServiceAccount'-type of configuration",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Secret token",
				},
				"ca_cert": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Service account certificate",
				},
			},
		},
	}
}

func resourceServiceEndpointKubernetes() *schema.Resource {
	r := crud.GenBaseServiceEndpointResource(flattenServiceEndpointKubernetes, expandServiceEndpointKubernetes)
	r.Schema["apiserver_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "URL to Kubernete's API-Server",
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
	}
	r.Schema["authorization_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Type of credentials to use",
		ValidateFunc: validation.StringInSlice([]string{"AzureSubscription", "Kubeconfig", "ServiceAccount"}, false),
	}
	makeSchemaAzureSubscription(r)
	makeSchemaKubeconfig(r)
	makeSchemaServiceAccount(r)

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointKubernetes(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := crud.DoBaseExpansion(d)
	serviceEndpoint.Type = converter.String("kubernetes")
	serviceEndpoint.Url = converter.String(d.Get("apiserver_url").(string))

	switch d.Get("authorization_type").(string) {
	case "AzureSubscription":
		configurationRaw := d.Get("azure_subscription").(*schema.Set).List()
		configuration := configurationRaw[0].(map[string]interface{})
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"azureEnvironment": configuration["azure_environment"].(string),
				"azureTenantId":    configuration["tenant_id"].(string),
			},
			Scheme: converter.String("Kubernetes"),
		}

		clusterID := fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.ContainerService/managedClusters/%s", configuration["subscription_id"].(string), configuration["resourcegroup_id"].(string), configuration["cluster_name"].(string))
		serviceEndpoint.Data = &map[string]string{
			"authorizationType":     "AzureSubscription",
			"azureSubscriptionId":   configuration["subscription_id"].(string),
			"azureSubscriptionName": configuration["subscription_name"].(string),
			"clusterId":             clusterID,
			"namespace":             configuration["namespace"].(string),
		}
	case "Kubeconfig":
		configurationRaw := d.Get("kubeconfig").(*schema.Set).List()
		configuration := configurationRaw[0].(map[string]interface{})

		clusterContextInput := configuration["cluster_context"].(string)
		if clusterContextInput == "" {
			kubeConfigYAML := configuration["kube_config"].(string)
			var kubeConfigYAMLUnmarshalled map[string]interface{}
			err := yaml.Unmarshal([]byte(kubeConfigYAML), &kubeConfigYAMLUnmarshalled)
			if err != nil {
				errResult := fmt.Errorf("kube_config contains an invalid YAML: %s", err)
				return nil, nil, errResult
			}
			clusterContextInputList := kubeConfigYAMLUnmarshalled["contexts"].([]interface{})[0].(map[interface{}]interface{})
			clusterContextInput = clusterContextInputList["name"].(string)
		}

		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"clusterContext": clusterContextInput,
				"kubeconfig":     configuration["kube_config"].(string),
			},
			Scheme: converter.String("Kubernetes"),
		}

		serviceEndpoint.Data = &map[string]string{
			"authorizationType":    "Kubeconfig",
			"acceptUntrustedCerts": fmt.Sprintf("%v", configuration["accept_untrusted_certs"].(bool)),
		}
	case "ServiceAccount":
		configurationRaw := d.Get("service_account").(*schema.Set).List()
		configuration := configurationRaw[0].(map[string]interface{})

		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"apiToken":                  configuration["token"].(string),
				"serviceAccountCertificate": configuration["ca_cert"].(string),
			},
			Scheme: converter.String("Token"),
		}

		serviceEndpoint.Data = &map[string]string{
			"authorizationType": "ServiceAccount",
		}

	}

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointKubernetes(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	crud.DoBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("authorization_type", (*serviceEndpoint.Data)["authorizationType"])
	d.Set("apiserver_url", (*serviceEndpoint.Url))

	switch (*serviceEndpoint.Data)["authorizationType"] {
	case "AzureSubscription":
		azureSubscriptionResource := &schema.Resource{
			Schema: map[string]*schema.Schema{},
		}
		makeSchemaAzureSubscription(azureSubscriptionResource)

		clusterIDSplit := strings.Split((*serviceEndpoint.Data)["clusterId"], "/")
		var clusterNameIndex int
		var resourceGroupIDIndex int
		for k, v := range clusterIDSplit {
			if v == "resourcegroups" {
				resourceGroupIDIndex = k + 1
			}
			if v == "managedClusters" {
				clusterNameIndex = k + 1
			}
		}
		configItems := []interface{}{
			map[string]interface{}{
				"azure_environment": (*serviceEndpoint.Authorization.Parameters)["azureEnvironment"],
				"tenant_id":         (*serviceEndpoint.Authorization.Parameters)["azureTenantId"],
				"subscription_id":   (*serviceEndpoint.Data)["azureSubscriptionId"],
				"subscription_name": (*serviceEndpoint.Data)["azureSubscriptionName"],
				"cluster_name":      clusterIDSplit[clusterNameIndex],
				"resourcegroup_id":  clusterIDSplit[resourceGroupIDIndex],
				"namespace":         (*serviceEndpoint.Data)["namespace"],
			},
		}

		azureSubscriptionSchemaSet := schema.NewSet(schema.HashResource(azureSubscriptionResource), configItems)
		d.Set("azure_subscription", azureSubscriptionSchemaSet)
	case "Kubeconfig":
		kubeconfigResource := &schema.Resource{
			Schema: map[string]*schema.Schema{},
		}
		makeSchemaKubeconfig(kubeconfigResource)

		acceptUntrustedCerts, _ := strconv.ParseBool((*serviceEndpoint.Data)["acceptUntrustedCerts"])
		configItems := []interface{}{
			map[string]interface{}{
				"kube_config":            (*serviceEndpoint.Authorization.Parameters)["kubeconfig"],
				"cluster_context":        (*serviceEndpoint.Authorization.Parameters)["clusterContext"],
				"accept_untrusted_certs": acceptUntrustedCerts,
			},
		}

		kubeConfigSchemaSet := schema.NewSet(schema.HashResource(kubeconfigResource), configItems)
		d.Set("kubeconfig", kubeConfigSchemaSet)
	case "ServiceAccount":
		serviceAccountResource := &schema.Resource{
			Schema: map[string]*schema.Schema{},
		}
		makeSchemaServiceAccount(serviceAccountResource)

		configItems := []interface{}{
			map[string]interface{}{
				"token":   (*serviceEndpoint.Authorization.Parameters)["apiToken"],
				"ca_cert": (*serviceEndpoint.Authorization.Parameters)["serviceAccountCertificate"],
			},
		}

		serviceAccountSchemaSet := schema.NewSet(schema.HashResource(serviceAccountResource), configItems)
		d.Set("service_account", serviceAccountSchemaSet)
	}
}
