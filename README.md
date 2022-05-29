# Vault Secrets Engine for CircleCI Contexts

<!-- [![Build Status](https://travis-ci.com/hashicorp/vault-plugin-secrets-gcpkms.svg?token=xjv5yxmcgdD1zvpeR4me&branch=master)](https://travis-ci.com/hashicorp/vault-plugin-secrets-gcpkms) -->

This is a plugin backend for [HashiCorp Vault][vault] that manages [CircleCI Contexts][contexts] 

**Please note:** Security is taken seriously. If you believe you have found a
security issue, **do not open an issue**. Responsibly disclose by contacting
security@hashicorp.com.


## Usage

The CircleCI Vault secrets engine is not bundled and included
in [Vault][vault] distributions. To use the plugin, run:

```shell script
go install github.com/bobthebuilderberlin/vault-plugin-secrets-circleci
cp "$GOPATH/bin/vault-plugin-secrets-circleci" bin/
vault server \
  -dev \
  -dev-plugin-dir="$(pwd)/bin" 
vault secrets enable -path=circleci -plugin=vault-plugin-secrets-circleci plugin
```

To configure the plugin use the /config endpoint:

```shell script
vault write circleci/config api-token="<api-token>" org-id="<org-id>"
```
where the `api-token` is an API token create [here](https://app.circleci.com/settings/user/tokens) and the org-id is the Organization ID that can be found in the Overview of the Settings for your CircleCI Organization. 

To list all you  CircleCI contexts:
```shell script
vault list circleci/context
```

To create a new CircleCI context:
```shell script
vault write circleci/context context=my-context
```

To list environment variables in a context
```shell script
vault list circleci/context/test-robert-1
```

To write a new environment variable:
```shell script
vault write circleci/context/my-context/foo value=bar
```


## Development

Prerequisites:

- Modern [Go](https://golang.org) (1.11+)
- Git

1. Clone the repo:

    ```shell script
    git clone https://github.com/bobthebuilderberlin/vault-plugin-secrets-circleci
    cd vault-plugin-secrets-circleci
    ```

1. Build the binary:

    ```shell script
    $ make dev
    ```

1. Copy the compiled binary into a scratch dir:

    ```shell script
    $ cp $(which vault-plugin-secrets-circleci) ./bin/
    ```

1. Run Vault plugins from that directory:

    ```shell script
    $ vault server -dev -dev-plugin-dir=./bin
    $ vault secrets enable -path=circleci -plugin=vault-plugin-secrets-circleci plugin
    ```

### Tests


[contexts]: https://circleci.com/docs/2.0/contexts/
[vault]: https://www.vaultproject.io
