Terraform Provider
==================

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.11+
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/norainacloud/terraform-provider-noraina`

```sh
$ mkdir -p $GOPATH/src/github.com/norainacloud/terraform-provider-noraina; cd $GOPATH/src/github.com/norainacloud/terraform-provider-noraina
$ git clone git@github.com:norainacloud/terraform-provider-noraina.git .
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/norainacloud/terraform-provider-noraina
$ go install
```

If you check your `$GOPATH/bin` folder you should see a freshly compiled binary of this provider

And now, in order to get terraform auto-discover mechanism to work you need to create a `~/.terraformrc` file with this content

```hcl
providers {
  noraina = "${GOPATH}/bin/terraform-provider-noraina"
}
```

Using The Provider
---------------------

Terraform code example

```hcl
provider "noraina" {
  email = "<< YOUR EMAIL >>"
  password = "<< YOUR PASSWORD >>"
}

resource "noraina_ece" "test" {
  name = "test"

  service {
    name = "s1"
    fqdn = "s1.domain.com"
    origin_hostheader = "noraina-eu-west-1.s3-eu-west-1.amazonaws.com"
    origin_backend = "http://noraina-eu-west-1.s3-eu-west-1.amazonaws.com"
    provider_region = "eu-west-1"
    provider_name = "aws"
    cert_id = "${noraina_certificate.test.id}"
  }

  service {
    name = "s2"
    fqdn = "s2.domain.com"
    origin_hostheader = "noraina-us-east-1.s3.amazonaws.com"
    origin_backend = "http://noraina-us-east-1.s3.amazonaws.com"
    provider_region = "us-east-1"
    provider_name = "aws"
  }
}

resource "noraina_certificate" "test" {
  name  = "ricard_test"
  cert  = "cert_test.crt"
  key   = "cert_test.key"
  chain = "cert_test.chain"
}
```

Please note the plugin supports the following environment variables: NORAINA_EMAIL and NORAINA_PASSWORD
so you can leave these 2 parameters empty if they are defined in your shell session

In order to execute this just run:

```sh
> terraform init

> terraform plan

> terraform apply
```

Every time you build a new version of the provider you will have to run


```sh
> terraform init --upgrade
```

So that the .terraform folder with your plugins state is updated with the new binary
