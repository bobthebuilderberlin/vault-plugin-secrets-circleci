package circleci

import (
	"context"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
	"testing"
)

// testBackend creates a new isolated instance of the backend for testing.
func testBackend(tb testing.TB) (*backend, logical.Storage) {
	tb.Helper()

	config := logical.TestBackendConfig()
	config.StorageView = new(logical.InmemStorage)
	config.Logger = hclog.NewNullLogger()

	b, err := Factory(context.Background(), config)
	if err != nil {
		tb.Fatal(err)
	}
	return b.(*backend), config.StorageView
}

// testFieldValidation verifies the given path has field validation.
func testFieldValidation(tb testing.TB, op logical.Operation, pth string) {
	tb.Helper()

	b, storage := testBackend(tb)
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Storage:   storage,
		Operation: op,
		Path:      pth,
		Data: map[string]interface{}{
			"literally-never-a-key": true,
		},
	})
	if err == nil {
		tb.Error("expected error")
	}
	if !strings.Contains(err.Error(), "unknown field") {
		tb.Error(err)
	}
}
//
//// testKMSClient creates a new KMS client with the default scopes and user
//// agent.
//func testKMSClient(tb testing.TB) *kmsapi.KeyManagementClient {
//	tb.Helper()
//
//	ctx := context.Background()
//	kmsClient, err := kmsapi.NewKeyManagementClient(ctx,
//		option.WithScopes(defaultScope),
//		option.WithUserAgent(useragent.String()),
//	)
//	if err != nil {
//		tb.Fatalf("failed to create kms client: %s", err)
//	}
//
//	return kmsClient
//}
//
//// testKMSKeyRingName creates a keyring name. If the given "name" is
//// blank, a UUID name is generated.
//func testKMSKeyRingName(tb testing.TB, name string) string {
//	tb.Helper()
//
//	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
//	if project == "" {
//		tb.Fatal("missing GOOGLE_CLOUD_PROJECT")
//	}
//
//	if name == "" {
//		name = fmt.Sprintf("vault-test-%s", uuid.NewV4())
//	}
//
//	return fmt.Sprintf("projects/%s/locations/us-east1/keyRings/%s", project, name)
//}
//
//// testCreateKMSKeyRing creates a keyring with the given name.
//func testCreateKMSKeyRing(tb testing.TB, name string) (string, func()) {
//	tb.Helper()
//
//	keyRing := testKMSKeyRingName(tb, name)
//
//	kmsClient := testKMSClient(tb)
//
//	// Check if the key ring exists
//	ctx := context.Background()
//	kr, err := kmsClient.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{
//		Name: keyRing,
//	})
//	if err != nil {
//		if terr, ok := grpcstatus.FromError(err); ok && terr.Code() == grpccodes.NotFound {
//			// Key ring does not exist, try to create it
//			kr, err = kmsClient.CreateKeyRing(ctx, &kmspb.CreateKeyRingRequest{
//				Parent:    path.Dir(path.Dir(keyRing)),
//				KeyRingId: path.Base(keyRing),
//			})
//			if err != nil {
//				tb.Fatalf("failed to create keyring: %s", err)
//			}
//		} else {
//			tb.Fatalf("failed to get keyring: %s", err)
//		}
//	}
//
//	return kr.Name, func() { testCleanupKeyRing(tb, kr.Name) }
//}
//
//// testCreateKMSCryptoKeySymmetric creates a new crypto key under the
//// vault-gcpkms-plugin-test key ring in the given google project.
//func testCreateKMSCryptoKeySymmetric(tb testing.TB) (string, func()) {
//	return testCreateKMSCryptoKeyPurpose(tb,
//		kmspb.CryptoKey_ENCRYPT_DECRYPT,
//		kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
//	)
//}
//
//// testCreateKMSCryptoKeyAsymmetricDecrypt creates a new KMS crypto key that is
//// used for asymmetric decryption.
//func testCreateKMSCryptoKeyAsymmetricDecrypt(tb testing.TB, algo kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm) (string, func()) {
//	return testCreateKMSCryptoKeyPurpose(tb,
//		kmspb.CryptoKey_ASYMMETRIC_DECRYPT,
//		algo,
//	)
//}
//
//// testCreateKMSCryptoKeyAsymmetricSign creates a new KMS crypto key that is
//// used for asymmetric signing.
//func testCreateKMSCryptoKeyAsymmetricSign(tb testing.TB, algo kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm) (string, func()) {
//	return testCreateKMSCryptoKeyPurpose(tb,
//		kmspb.CryptoKey_ASYMMETRIC_SIGN,
//		algo,
//	)
//}
//
//
//
//func TestBackend_KMSClient(t *testing.T) {
//	t.Parallel()
//
//	t.Run("allows_concurrent_reads", func(t *testing.T) {
//		t.Parallel()
//
//		b, storage := testBackend(t)
//
//		_, closer1, err := b.CircleCIClient(storage)
//		if err != nil {
//			t.Fatal(err)
//		}
//		defer closer1()
//
//		doneCh := make(chan struct{})
//		go func() {
//			_, closer2, err := b.CircleCIClient(storage)
//			if err != nil {
//				t.Fatal(err)
//			}
//			defer closer2()
//			close(doneCh)
//		}()
//
//		select {
//		case <-doneCh:
//		case <-time.After(1 * time.Second):
//			t.Errorf("client was not available")
//		}
//	})
//
//	t.Run("caches", func(t *testing.T) {
//		t.Parallel()
//
//		b, storage := testBackend(t)
//
//		client1, closer1, err := b.CircleCIClient(storage)
//		if err != nil {
//			t.Fatal(err)
//		}
//		defer closer1()
//
//		client2, closer2, err := b.CircleCIClient(storage)
//		if err != nil {
//			t.Fatal(err)
//		}
//		defer closer2()
//
//		// Note: not a bug; literally checking object equality
//		if client1 != client2 {
//			t.Errorf("expected %#v to be %#v", client1, client2)
//		}
//	})
//
//	t.Run("expires", func(t *testing.T) {
//		t.Parallel()
//
//		b, storage := testBackend(t)
//		b.kmsClientLifetime = 50 * time.Millisecond
//
//		client1, closer1, err := b.CircleCIClient(storage)
//		if err != nil {
//			t.Fatal(err)
//		}
//		closer1()
//
//		time.Sleep(100 * time.Millisecond)
//
//		client2, closer2, err := b.CircleCIClient(storage)
//		if err != nil {
//			t.Fatal(err)
//		}
//		closer2()
//
//		if client1 == client2 {
//			t.Errorf("expected %#v to not be %#v", client1, client2)
//		}
//	})
//}
//
//func TestBackend_ResetClient(t *testing.T) {
//	t.Parallel()
//
//	t.Run("closes_client", func(t *testing.T) {
//		t.Parallel()
//
//		b, storage := testBackend(t)
//
//		client, closer, err := b.CircleCIClient(storage)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		// Verify the client is "open"
//		if client.Connection().GetState() == connectivity.Shutdown {
//			t.Fatalf("connection is already stopped")
//		}
//
//		// Stop read lock
//		closer()
//
//		// Reset the clients
//		b.ResetClient()
//
//		// Verify the client closed
//		if state := client.Connection().GetState(); state != connectivity.Shutdown {
//			t.Errorf("expected client to be closed, was: %v", state)
//		}
//	})
//}
//
//func TestBackend_Config(t *testing.T) {
//	t.Parallel()
//
//	cases := []struct {
//		name string
//		c    []byte
//		e    *Config
//		err  bool
//	}{
//		{
//			"default",
//			nil,
//			DefaultConfig(),
//			false,
//		},
//		{
//			"saved",
//			[]byte(`{"credentials":"foo", "scopes":["bar"]}`),
//			&Config{
//				Credentials: "foo",
//				Scopes:      []string{"bar"},
//			},
//			false,
//		},
//		{
//			"invalid",
//			[]byte(`{x`),
//			nil,
//			true,
//		},
//	}
//
//	for _, tc := range cases {
//		tc := tc
//
//		t.Run(tc.name, func(t *testing.T) {
//			t.Parallel()
//
//			b, storage := testBackend(t)
//
//			if tc.c != nil {
//				if err := storage.Put(context.Background(), &logical.StorageEntry{
//					Key:   "config",
//					Value: tc.c,
//				}); err != nil {
//					t.Fatal(err)
//				}
//			}
//
//			c, err := b.Config(context.Background(), storage)
//			if (err != nil) != tc.err {
//				t.Fatal(err)
//			}
//
//			if !reflect.DeepEqual(c, tc.e) {
//				t.Errorf("expected %#v to be %#v", c, tc.e)
//			}
//		})
//	}
//}
