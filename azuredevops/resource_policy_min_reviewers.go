package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	
	// "github.com/microsoft/azure-devops-go-api/azuredevops/policy"

	// "github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)


func resourcePolicyMinReviewers() *schema.Resource {
	return &schema.Resource{
	  Create: resourcePolicyMinReviewersCreate,
	  Read:   resourcePolicyMinReviewersRead,
	  Update: resourcePolicyMinReviewersUpdate,
	  Delete: resourcePolicyMinReviewersDelete,
  
	  Schema: map[string]*schema.Schema{
			"project_id" :&schema.Schema{
			  Type:     schema.TypeString,
			  Required: true,
			},
			"enabled": &schema.Schema{
			  Type:     schema.TypeBool,
			  Optional: true,
			  Default:	true,
			},
			"is_blocking": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:	true,
			},
			"creator_vote_counts": &schema.Schema{
			  Type:     schema.TypeBool,
			  Optional: true,
			  Default:	false,
			},
			"minimum_approver_count": &schema.Schema{
			  Type:     schema.TypeInt,
			  Optional: true,
			  Default:	1,
			},
			"scope": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource {
				  Schema: map[string]*schema.Schema{
						"repository_id": {
					  	Type:     schema.TypeString,
					  	Required: true,
						},		  
						"repository_ref": {
					  	Type:     schema.TypeString,
					  	Optional:		true,
					  	Default:		"refs/heads/master",
						},		  
						"match_type": {
					  	Type:      schema.TypeString,
					  	Optional:		true,
					  	Default:		"Exact",
						},
				  },
				},
			},
		},
	}
 }

	func resourcePolicyMinReviewersCreate(d *schema.ResourceData, m interface{}) error {
		// clients := m.(*aggregatedClient)
		// projectID := d.Get("project_id").(string)
		// policyTypes, err := clients.PolicyClient.GetPolicyTypes(clients.ctx, policy.GetPolicyTypesArgs{
		// 	Project: converter.String(projectID),
		// })
		// if err != nil {
		// 	log.Printf("something bad happened %s\n", err)
		// }
		// if (policyTypes != nil){
		// 	log.Printf("resource_policy_min_reviewers get policy types %+v\n", policyTypes)
		// }
		
		policy_id := "1234" //this should come from api
		d.SetId(policy_id)
		return resourcePolicyMinReviewersRead(d, m)
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