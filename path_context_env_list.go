package circleci

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathContextEnvList() *framework.Path {
	return &framework.Path{
		Pattern: "context/" + framework.GenericNameRegex("context") + "/?$",

		HelpSynopsis:    "List a context's environment variables",
		HelpDescription: "TODO: write description for path",

		Fields: map[string]*framework.FieldSchema{
			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the CircleCI context you would like to create.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: withFieldValidator(b.pathContextEnvLister),
		},
	}
}

// pathContextsList corresponds to PUT/POST gcpkms/decrypt/:key and is
// used to decrypt the ciphertext string using the named key.
func (b *backend) pathContextEnvLister(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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
			listOfVariableNames := make([]string, len(contextVariableList.Items))
			for i:=0; i < len(contextVariableList.Items) ;i++ {
				listOfVariableNames[i] = contextVariableList.Items[i].Variable
			}
			if err != nil {
				return nil, err
			}
			return logical.ListResponse(listOfVariableNames), nil
		}
	}
	return nil, fmt.Errorf("no context with name '%v' was found", circleCIContext)
}
