# terraform-provider-dbt

This custom terraform provider is used to manage resources for DBT cloud

The provider is published to [registry.terraform.io/providers/3lvia/dbt](https://registry.terraform.io/providers/3lvia/dbt/latest).

It can be used to manage groups and group permissions in DBT cloud

# General information about creating a custom terraform provider

Se [here](https://learn.hashicorp.com/collections/terraform/providers) for the general information about creating custom providers from Terraform.

# Local Setup
## Install go
[Golang installation guide](https://golang.org/doc/install)

## Install terraform
[Install terraform](https://learn.hashicorp.com/terraform/getting-started/install.html). For windows you can add terraform.exe to {user}/bin. Make sure %USERPROFILE%\go\bin is in path, and above the go-spesific paths.

## Checkout code
Checkout the code-repo to {GOPATH}\src\github.com\3lvia\terraform-provider-dbt

# Project structure
* repo-root
  * examples: folder with terraform files for manually testing the provider.
  * main.go: Standard file, sets up serving of the provider by calling the Provider()-function.
  * provider.go: Defines the provider schema (inputs to the provider), the mapping to resorces, and the interface that is passed to resrouces
  * resource_usergroup.go: Defines the resource schema and methods for usergroups.


## Adding terraform.tfvars to terraform-tester
Create ../terraform-tester/terraform.tfvars and add these variables.

```
service_token = "replaceme"
account_id = "replaceme
```

Note that terraform.tfvars is added to .gitignore. Make sure to newer publish these secrets. This is a public repository.

# Running locally

## Build the provider for a local run

```console
# from repo-root
make
```

## Running terraform locally

```console
# from repo-root/examples
terraform init
terraform apply

```

# Debugging
Debugging the go-code when running from terraform is not suported. It is possible to print debug info as warnings in diag.Diagnostics. Debugging can also be done by writing a file with debug messages:

```
# for a string 
ioutil.WriteFile("custom-log.text", []byte(someString), 0644)

# for an object
serialized, _ := json.Marshal(someObject)
ioutil.WriteFile("custom-log.text", []byte(serialized), 0644)
```
# Publish a new release
## Publish to terraform registry
To publish to [registry.terraform.io/providers/3lvia/dbt](https://registry.terraform.io/providers/3lvia/dbt/latest) create a new github-release in this repo. 
Github-actions is setup to atomatically build and publish new releases. 

Github-actions uses our private signing key to sign the build. The public variant of this key is added in terraform registry.
Backup of the key is found in vault (prod) edna/kv/manual/dbt-provider-build-signing-key
