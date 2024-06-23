# Terraform Provider pfSense

[![Registry](https://img.shields.io/badge/pfsense-Terraform%20Registry-blue)](https://registry.terraform.io/providers/marshallford/pfsense/latest/docs)

Used to configure [pfSense](https://www.pfsense.org/) firewall/router systems with Terraform. Validated with pfSense CE, compatibility with pfSense Plus is not guaranteed.

> [!WARNING]
> All versions released prior to `v1.0.0` are to be considered [breaking changes](https://semver.org/#how-do-i-know-when-to-release-100).

## Support Matrix

| Release  | pfSense          | Terraform      |
| :------: | :--------------: | :------------: |
| < v1.0.0 | 2.6.0 CE         | >= 1.3, <= 1.5 |

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.3
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Install go dependencies by running the go `install` command:
```shell
go install
```
5. Build the provider using the makefile `build` recipe:
```shell
make build
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will install go dependencies. Then you'll be able to build the provider by running `make build`, this will put the provider binary in the `build` directory.

To generate or update documentation, run `make docs`.

In orded to test the provider using a terraform project, add the following lines to your `$HOME/.terraformrc` file and make sure to replace `{{PATH_TO_BIN}}` for the actual built binary path.
```terraform
provider_installation {
  dev_overrides {
    "registry.terraform.io/marshallford/pfsense" = "{{PATH_TO_BIN}}"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

In order to run the full suite of Acceptance tests, run `make test/acc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make test/acc
```
