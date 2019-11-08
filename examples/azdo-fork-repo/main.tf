# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
provider "azuredevops" {
  version = ">= 0.0.1"
}


// This section configures an Azure DevOps Git Repository with branch policies
resource "azuredevops_azure_git_repository" "repository" {
  project_id = "668627d8-dda3-4bc2-8fb6-b2059eda7b7d"
  name       = "Sample-Fork-Repo"
  initialization {
    init_type = "Fork"
    source_type = "test"
    source_url = "test"
  }
}