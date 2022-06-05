package circleci

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_PathConfigRead(t *testing.T) {
	t.Parallel()

	t.Run("field_validation", func(t *testing.T) {
		t.Parallel()
		testFieldValidation(t, logical.ReadOperation, "config")
	})

	t.Run("not_exist", func(t *testing.T) {
		t.Parallel()

		b, storage := testBackend(t)
		ctx := context.Background()
		resp, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.ReadOperation,
			Path:      "config",
		})
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := resp.Data["APIToken"]; !ok {
			t.Errorf("expected %q to include %q", resp.Data, "api-token")
		}
		if _, ok := resp.Data["OrgId"]; !ok {
			t.Errorf("expected %q to include %q", resp.Data, "org-id")
		}
	})

	t.Run("exist", func(t *testing.T) {
		t.Parallel()

		b, storage := testBackend(t)

		entry, err := logical.StorageEntryJSON("config", &Config{
			APIToken: "my-token",
			OrgId:    "my-org-id",
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.Put(context.Background(), entry); err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()
		resp, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.ReadOperation,
			Path:      "config",
		})
		if err != nil {
			t.Fatal(err)
		}

		if resp.Data["APIToken"].(string) != "my-token" {
			t.Errorf("expected api-token to be 'my-token'")
		}
		if resp.Data["OrgId"].(string) != "my-org-id" {
			t.Errorf("expected org-id to be 'my-org-id'")
		}
	})
}

func TestBackend_PathConfigUpdate(t *testing.T) {
	t.Parallel()

	t.Run("field_validation", func(t *testing.T) {
		t.Parallel()
		testFieldValidation(t, logical.UpdateOperation, "config")
	})

	t.Run("not_exist", func(t *testing.T) {
		t.Parallel()

		b, storage := testBackend(t)
		if _, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.UpdateOperation,
			Path:      "config",
			Data: map[string]interface{}{
				"api-token": "my-token",
				"org-id":    "my-org-id",
			},
		}); err != nil {
			t.Fatal(err)
		}

		config, err := b.Config(context.Background(), storage)
		if err != nil {
			t.Fatal(err)
		}

		if v, exp := config.APIToken, "my-token"; v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}

		if v, exp := config.OrgId, "my-org-id"; v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}
	})

	t.Run("exist", func(t *testing.T) {
		t.Parallel()

		b, storage := testBackend(t)

		entry, err := logical.StorageEntryJSON("config", &Config{
			APIToken: "my-token",
			OrgId:    "my-org-id",
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.Put(context.Background(), entry); err != nil {
			t.Fatal(err)
		}

		if _, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.UpdateOperation,
			Path:      "config",
			Data: map[string]interface{}{
				"api-token": "my-new-token",
				"org-id":    "my-new-org-id",
			},
		}); err != nil {
			t.Fatal(err)
		}

		config, err := b.Config(context.Background(), storage)
		if err != nil {
			t.Fatal(err)
		}

		if v, exp := config.APIToken, "my-new-token"; v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}

		if v, exp := config.OrgId, "my-new-org-id"; v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}
	})
}

func TestBackend_PathConfigDelete(t *testing.T) {
	t.Parallel()

	t.Run("field_validation", func(t *testing.T) {
		t.Parallel()
		testFieldValidation(t, logical.DeleteOperation, "config")
	})

	t.Run("not_exist", func(t *testing.T) {
		t.Parallel()

		b, storage := testBackend(t)
		if _, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.DeleteOperation,
			Path:      "config",
		}); err != nil {
			t.Fatal(err)
		}

		config, err := b.Config(context.Background(), storage)
		if err != nil {
			t.Fatal(err)
		}

		if def := DefaultConfig(); !reflect.DeepEqual(config, def) {
			t.Errorf("expected %v to be %v", config, def)
		}
	})

	t.Run("exist", func(t *testing.T) {
		t.Parallel()

		b, storage := testBackend(t)

		entry, err := logical.StorageEntryJSON("config", &Config{
			APIToken: "my-token",
			OrgId:    "my-org-id",
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.Put(context.Background(), entry); err != nil {
			t.Fatal(err)
		}

		if _, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.DeleteOperation,
			Path:      "config",
		}); err != nil {
			t.Fatal(err)
		}

		config, err := b.Config(context.Background(), storage)
		if err != nil {
			t.Fatal(err)
		}

		if def := DefaultConfig(); !reflect.DeepEqual(config, def) {
			t.Errorf("expected %v to be %v", config, def)
		}
	})
}
