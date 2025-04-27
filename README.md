# Terraform Provider for the Hetzner Robot Interface

The hrobot provider is used to interact with the resources from the Hetzner Robot Interface (dedicated server). The provider needs to be configured with the proper credentials before it can be used. 

> Please note this is not an official Hetzner product, the author is not in any way affiliated with Hetzner use at own risk!  

If you are looking for the [Hetzner Cloud](https://cloud.hetzner.com) terraform provider you can check out the [official terraform-provider-hcloud](https://github.com/hetznercloud/terraform-provider-hcloud) provider maintained by Hetzner. 

## Authentication 

The terraform provider needs a Hetzner token to authenticate. This can either be provided in the provider configuration in your main.tf file or via a Environment Variable with the name `HROBOT_TOKEN`. 

The "Token" has the following structure: "username:password".
To use the Token as an enviroment variable as in the example above you can export a variable: `export HROBOT_TOKEN="username:password"` in your terminal. To make it persitent on your system you can put the export command in your ~/.profile file. 

## Example Configuration
``` terraform
# Set the variable value in *.tfvars file
# or using the -var="hrobot_token=..." CLI option
variable "hrobot_token" {
  sensitive = true # Requires terraform >= 0.14
}

# Configure the Hetzner Hrobot Provider
provider "hrobot" {
  token = var.hrobot_token
}
```
