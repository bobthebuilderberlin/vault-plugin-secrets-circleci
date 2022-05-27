package circleci

import (
	"context"
	"errors"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"sync"

	circleci "github.com/grezar/go-circleci"
)

type backend struct {
	*framework.Backend

	// circleciClient is the actual client for connecting to CircleCI. It is cached on
	// the backend for efficiency.
	circleciClient      *circleci.Client

	// ctx and ctxCancel are used to control overall plugin shutdown. These
	// contexts are given to any client libraries or requests that should be
	// terminated during plugin termination.
	ctx       context.Context
	ctxCancel context.CancelFunc
	ctxLock   sync.Mutex
}

// Factory returns a configured instance of the backend.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns a configured instance of the backend.
func Backend() *backend {
	var b backend

	b.ctx, b.ctxCancel = context.WithCancel(context.Background())

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help: "CircleCI secrets engine.",

		Paths: []*framework.Path{
			b.pathConfig(),
			b.pathContexts(),
			b.pathContextKey(),
		},

		Invalidate: b.invalidate,
		Clean:      b.clean,
	}

	return &b
}

// clean cancels the shared contexts. This is called just before unmounting
// the plugin.
func (b *backend) clean(_ context.Context) {
	b.ctxLock.Lock()
	b.ctxCancel()
	b.ctxLock.Unlock()
}

// invalidate resets the plugin. This is called when a key is updated via
// replication.
func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.ResetClient()
	}
}

// ResetClient closes any connected clients.
func (b *backend) ResetClient() {
	b.circleciClient = nil
}


// CircleCIClient creates a new client for talking to the GCP KMS service.
func (b *backend) CircleCIClient(s logical.Storage) (*circleci.Client, func(), error) {
	// If the client already exists and is valid, return it
	b.ctxLock.Lock()
	if b.circleciClient != nil {
		closer := func() { b.ctxLock.Unlock() }
		return b.circleciClient, closer, nil
	}

	b.Logger().Debug("Creating new CircleCI Client...")

	// Attempt to close an existing client if we have one.
	b.ResetClient()

	// Get the config
	config, err := b.Config(b.ctx, s)
	b.Logger().Debug("CircleCI configuration:", "config", config)

	if err != nil {
		b.ctxLock.Unlock()
		return nil, nil, err
	}

	// If credentials were provided, use those. Otherwise fall back to the
	// default application credentials.
	if len(config.APIToken) == 0 {
		b.ctxLock.Unlock()
		return nil, nil, errors.New("APIToken must not be empty or nil")
	}

	circleCIConfig:= circleci.DefaultConfig()
	circleCIConfig.Token = config.APIToken

	// Create and return the CircleCI client
	client, err := circleci.NewClient(circleCIConfig)

	if err != nil {
		b.ctxLock.Unlock()
		return nil, nil, errwrap.Wrapf("Failed to create CircleCI client: {{err}}", err)
	}

	b.Logger().Debug("CircleCI client created successfully.")

	// Cache the client
	b.circleciClient = client
	b.ctxLock.Unlock()
	closer := func() {
		b.ctxLock.TryLock()
	    b.ctxLock.Unlock()
	}
	return client, closer, nil
}

// Config parses and returns the configuration data from the storage backend.
// Even when no user-defined data exists in storage, a Config is returned with
// the default values.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*Config, error) {
	c := DefaultConfig()

	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, errwrap.Wrapf("failed to get configuration from storage: {{err}}", err)
	}
	if entry == nil || len(entry.Value) == 0 {
		return c, nil
	}

	if err := entry.DecodeJSON(&c); err != nil {
		return nil, errwrap.Wrapf("failed to decode configuration: {{err}}", err)
	}
	return c, nil
}
