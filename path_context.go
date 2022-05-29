package circleci

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathContext() *framework.Path {
	return &framework.Path{
		Pattern: "context/?$",

		HelpSynopsis:    "List contexts, create new contexts.",
		HelpDescription: "TODO: write description for path",

		Fields: map[string]*framework.FieldSchema{
			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the CircleCI context you would like to create.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation:   withFieldValidator(b.pathContextsList),
			logical.CreateOperation: withFieldValidator(b.pathContextWrite),
			logical.UpdateOperation: withFieldValidator(b.pathContextWrite),
			logical.DeleteOperation: withFieldValidator(b.pathContextDelete),
		},
	}
}

// pathContextsList corresponds to PUT/POST gcpkms/decrypt/:key and is
// used to decrypt the ciphertext string using the named key.
func (b *backend) pathContextsList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	config, err := b.Config(b.ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	collectedContexts, err := b.collectContexts(ctx, req, config)
	if err != nil {
		return nil, err
	}

	collectedContextNames := make([]string, len(collectedContexts))
	for i := 0; i < len(collectedContexts); i++ {
		collectedContextNames[i] = collectedContexts[i].Name
	}
	return logical.ListResponse(collectedContextNames), nil
}

// pathContextsList corresponds to PUT/POST gcpkms/decrypt/:key and is
// used to decrypt the ciphertext string using the named key.
func (b *backend) pathContextWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	circleCIClient, closer, err := b.CircleCIClient(req.Storage)
	defer closer()

	circleCIContext := d.Get("context").(string)
	if circleCIContext == "" {
		return nil, errors.New("'context' variable is required to create a new CircleCI context")
	}
	config, err := b.Config(b.ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	createdContext, err := circleCIClient.Contexts.Create(ctx, circleci.ContextCreateOptions{
		Name: &circleCIContext,
		Owner: &circleci.OwnerOptions{
			ID: &config.OrgId,
		},
	})

	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"context": createdContext,
		},
	}, nil
}

func (b *backend) pathContextDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	circleCIContext := d.Get("context").(string)
	if circleCIContext == "" {
		return nil, errors.New("'context' variable is required to delete CircleCI context")
	}
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
			err := circleCIClient.Contexts.Delete(ctx, collectedContexts[i].ID)
			if err != nil {
				return nil, err
			}
			return &logical.Response{
				Data: map[string]interface{}{
					"deletionSuccessful": true,
				},
			}, nil
		}
	}

	return nil, fmt.Errorf("context '%v'was not found", circleCIContext)
}

func (b *backend) collectContexts(ctx context.Context, req *logical.Request, config *Config) ([]*circleci.Context, error) {
	circleCIClient, closer, err := b.CircleCIClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	var collectedContexts []*circleci.Context
	var nextPageToken string
	for {
		contextList, err := circleCIClient.Contexts.List(ctx, circleci.ContextListOptions{OwnerID: &config.OrgId, PageToken: &nextPageToken})
		if err != nil {
			return nil, err
		}
		for _, contextListItem := range contextList.Items {
			collectedContexts = append(collectedContexts, contextListItem)
		}
		if contextList.NextPageToken == "" {
			break
		}
		nextPageToken = contextList.NextPageToken
	}
	return collectedContexts, nil
}
