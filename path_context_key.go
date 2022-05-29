package circleci

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathContextKey() *framework.Path {
	return &framework.Path{
		Pattern: "contexts/" + framework.GenericNameRegex("context") + "/" + framework.GenericNameRegex("env"),

		HelpSynopsis:    "Read and write environment variables in CircleCI contexts",
		HelpDescription: "TODO: write description for path",

		Fields: map[string]*framework.FieldSchema{
			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the CircleCI context you would like to alter.",
				Required:    true,
			},
			"env": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the environment variable you want to read or write in the given CircleCI context.",
				Required:    true,
			},
			"value": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the environment variable you want to read or write in the given CircleCI context.",
				Required:    true,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: withFieldValidator(b.pathContextKeyWrite),
			logical.UpdateOperation: withFieldValidator(b.pathContextKeyWrite),
		},
	}
}

// pathContextsList corresponds to PUT/POST gcpkms/decrypt/:key and is
// used to decrypt the ciphertext string using the named key.
func (b *backend) pathContextKeyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	circleCIContext := d.Get("context").(string)
	envVariable := d.Get("env").(string)
	value := d.Get("value").(string)
	config, err := b.Config(b.ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	contextList, err := b.collectContexts(ctx, req, config)
	if err != nil {
		return nil, err
	}

	circleCIClient, closer, err := b.CircleCIClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	for _, context := range contextList {
		if context.Name == circleCIContext {
			contextVariable, err := circleCIClient.Contexts.AddOrUpdateVariable(ctx, context.ID, envVariable, circleci.ContextAddOrUpdateVariableOptions{Value: &value})
			if err != nil {
				return nil, err
			}
			b.Logger().Debug("Variable in context successfully created or updated", "context", context.Name, "contextID", context.ID, "envVariable", contextVariable.Variable)
			return &logical.Response{
				Data: map[string]interface{}{
					"contextEnvironmentVariable": contextVariable.Variable,
				},
			}, nil
		}
	}
	return nil, errors.New("context with that name was not found")
}
