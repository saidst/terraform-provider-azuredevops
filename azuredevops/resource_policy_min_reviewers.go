package azuredevops

import (
	"strconv"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourcePolicyMinReviewers() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyMinReviewersCreate,
		Read:   resourcePolicyMinReviewersRead,
		Update: resourcePolicyMinReviewersUpdate,
		Delete: resourcePolicyMinReviewersDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"is_blocking": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"creator_vote_counts": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"minimum_approver_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"scope": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repository_ref": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "refs/heads/master",
						},
						"match_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Exact",
						},
					},
				},
			},
		},
	}
}

var policyTypeMinReviewer = "fa4e907d-c16b-4a4c-9dfa-4906e5d171dd";

func resourcePolicyMinReviewersCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	projectID := d.Get("project_id").(string)
	policy_id, _ := uuid.Parse(policyTypeMinReviewer)
	
	var scopes []map[string]interface{}
	tfScopes := d.Get("scope").(*schema.Set).List()
	for _, tfScope := range tfScopes {
		tfScopeMap := tfScope.(map[string]interface{})
		scope := map[string]interface{}{
			"refName": tfScopeMap["repository_ref"].(string),
			"matchKind": tfScopeMap["match_type"].(string),
			"repositoryId": tfScopeMap["repository_id"].(string),
		}
		scopes = append(scopes, scope)
	}	

	createPolicyConfigurationArgs := policy.CreatePolicyConfigurationArgs{
		Project: &projectID,
		Configuration: &policy.PolicyConfiguration{
			IsEnabled: converter.Bool(d.Get("enabled").(bool)),
			IsBlocking: converter.Bool(d.Get("is_blocking").(bool)),
			Type: &policy.PolicyTypeRef{
				Id: &policy_id,
			},
			Settings: map[string]interface{}{
				"minimumApproverCount": converter.Int(d.Get("minimum_approver_count").(int)),
				"creatorVoteCounts": converter.Bool(d.Get("creator_vote_counts").(bool)),
				"scope": scopes, 
			},
		},
	}

	policyConfiguration, err := clients.PolicyClient.CreatePolicyConfiguration(clients.ctx, createPolicyConfigurationArgs)
	resourceId := *policyConfiguration.Id
	d.SetId(strconv.Itoa(resourceId))
	return err
}

func resourcePolicyMinReviewersRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePolicyMinReviewersUpdate(d *schema.ResourceData, m interface{}) error {
	return resourcePolicyMinReviewersRead(d, m)
}

func resourcePolicyMinReviewersDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
