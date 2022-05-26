package circleci

import (
	"context"
	"errors"
	"github.com/grezar/go-circleci"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathContextKey() *framework.Path {
	return &framework.Path{
		Pattern: "contexts/" + framework.GenericNameRegex("context") + "/" + framework.GenericNameRegex("env"),

		HelpSynopsis: "Read and write environment variables in CircleCI contexts",
		HelpDescription: "TODO: write description for path",

		Fields: map[string]*framework.FieldSchema{
			"context": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "The name of the CircleCI context you would like to alter.",
			},
			"env": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "The name of the environment variable you want to read or write in the given CircleCI context.",
			},
			"value": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: "The name of the environment variable you want to read or write in the given CircleCI context.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: withFieldValidator(b.pathContextKeyWrite),
			logical.UpdateOperation: withFieldValidator(b.pathContextKeyWrite),
		},
	}
}

// pathContextKeyWrite corresponds to PUT/POST gcpkms/decrypt/:key and is
// used to decrypt the ciphertext string using the named key.
func (b *backend) pathContextKeyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	circleCIContext := d.Get("context").(string)
	envVariable := d.Get("env").(string)
	value := d.Get("value").(string)
	circleCIClient, closer, err := b.CircleCIClient(req.Storage)
	orgID:= "e88e9a5b-2c4a-4dff-a770-d1f340c465d1"
	contextList, err := circleCIClient.Contexts.List(ctx, circleci.ContextListOptions{OwnerID : &orgID })
	if err != nil {
		closer()
		return nil, err
	}

	for _, context := range contextList.Items {
		if context.Name == circleCIContext {
			contextVariable, err := circleCIClient.Contexts.AddOrUpdateVariable(ctx, context.ID, envVariable, circleci.ContextAddOrUpdateVariableOptions{Value: &value})
			closer()
			if err != nil {
				return nil, err
			}
			b.Logger().Debug("Variable in context successfully created or updated", "context", context.Name, "contextID", context.ID, "envVariable", contextVariable.Variable)
			return &logical.Response{
				Data: map[string]interface{}{
					"contextEnvironmentVariable" : contextVariable.Variable,
				},
			}, nil
		}
	}
	closer()
	return nil, errors.New("context with that name was not found")
}
