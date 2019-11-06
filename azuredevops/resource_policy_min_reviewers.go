package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)


func resourcePolicyMinReviewers() *schema.Resource {
	return &schema.Resource{
	  Create: resourceTodoListCreate,
	  Read:   resourceTodoListRead,
	  Update: resourceTodoListUpdate,
	  Delete: resourceTodoListDelete,
  
	  Schema: map[string]*schema.Schema{
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
		"scope": {
			Type:     schema.TypeSet,
			Elem: &schema.Resource{
			  Schema: map[string]*schema.Schema{
				"repository_id": {
				  Type:     schema.TypeString,
				  Required: true,
				},		  
				"repository_ref": {
				  Type:     schema.TypeString,
				  Optional:		true,
				  Default:		"refs/heads/master"
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