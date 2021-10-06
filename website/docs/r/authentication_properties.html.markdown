---
layout: "aci"
page_title: "ACI: aci_aaa_authentication"
sidebar_current: "docs-aci-resource-aaa_authentication"
description: |-
  Manages ACI AAA Authentication
---

# aci_aaa_authentication #

Manages ACI AAA Authentication

## API Information ##

* `Class` - aaaAuthRealm
* `Distinguished Named` - uni/userext/authrealm

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_aaa_authentication" "example" {

  annotation = "orchestrator:terraform"
  def_role_policy = "no-login"

}
```

## Argument Reference ##



* `annotation` - (Optional) Annotation of object AAA Authentication.

* `def_role_policy` - (Optional) Default Role Policy.The default role policy for the remote user with invalid (or no) CiscoAVPairs returned by the AAA server. CiscoAVPairs provide support for Remote Access Dial-In User Service attribute-value (AV) pairs. Allowed values are "assign-default-role", "no-login", and default value is "no-login". Type: String.


## Importing ##

An existing AAAAuthentication can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_aaa_authentication.example <Dn>
```