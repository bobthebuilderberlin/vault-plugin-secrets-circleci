module github.com/bobthebuilderberlin/vault-plugin-secrets-circleci

go 1.15

require (
	cloud.google.com/go/kms v1.4.0
	github.com/gammazero/workerpool v1.1.2
	github.com/golang/protobuf v1.5.2
	github.com/grezar/go-circleci v0.6.1 // indirect
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-hclog v0.16.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2
	github.com/hashicorp/vault/api v1.5.0
	github.com/hashicorp/vault/sdk v0.4.1
	github.com/jeffchao/backoff v0.0.0-20140404060208-9d7fd7aa17f2
	github.com/satori/go.uuid v1.2.0
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	google.golang.org/api v0.80.0
	google.golang.org/genproto v0.0.0-20220505152158-f39f71e6c8f3
	google.golang.org/grpc v1.46.2
)
