# Terraform Provider Avengers

This repo references the [Sourav Patnaik's tutorial] (https://blog.devgenius.io/custom-terraform-provider-design-c39287c816e9)


## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-avengers
```

## Local release build

```shell
$ go install github.com/goreleaser/goreleaser@latest
```

```shell
$ make release
```


## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```
