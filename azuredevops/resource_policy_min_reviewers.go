package azuredevops

import (
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
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
		clients := m.(*aggregatedClient)
		projectId := d.Get("project_id").(*string)
		policyTypes, error := clients.PolicyClient.GetPolicyTypes(clients.ctx, policy.GetPolicyTypesArgs{
			Project: projectId,
		})
		if error != nil {
			log.Printf("something bad happened")
		}
		log.Printf("resource_policy_min_reviewers get policy types %x\n", policyTypes)
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