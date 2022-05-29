package circleci

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathContextList() *framework.Path {
	return &framework.Path{
		Pattern: "context/" + framework.GenericNameRegex("context"),

		HelpSynopsis:    "List contexts env variables",
		HelpDescription: "TODO: write description for path",

		Fields: map[string]*framework.FieldSchema{
			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the CircleCI context you would like to create.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: withFieldValidator(b.pathContextEnvList),
		},
	}
}

// pathContextsList corresponds to PUT/POST gcpkms/decrypt/:key and is
// used to decrypt the ciphertext string using the named key.
func (b *backend) pathContextEnvList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	circleCIContext := d.Get("context").(string)

	config, err := b.Config(b.ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	collectedContexts, err := b.collectContexts(ctx, req, config)
	if err != nil {
		return nil, err
	}

	circleCIClient, closer, err := b.CircleCIClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	for i := 0; i < len(collectedContexts); i++ {
		if collectedContexts[i].Name == circleCIContext {
			contextVariableList, err := circleCIClient.Contexts.ListVariables(ctx, collectedContexts[i].ID)
			if err != nil {
				return nil, err
			}
			b.Logger().Debug("Number of variables: ", "size", len(contextVariableList.Items))
			return &logical.Response{
				Data: map[string]interface{}{
					"variables": contextVariableList.Items,
				},
			}, nil
		}
	}
	return nil, fmt.Errorf("no context with name '%v' was found", circleCIContext)
}
