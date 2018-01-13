package vaulttest

import (
	"context"
	"testing"

	vault "github.com/hashicorp/vault/api"
)

func TestRunVault(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, rootToken := Run(ctx)
	client, err := vault.NewClient(config)
	if err != nil {
		t.Fatalf("Unable to create vault client: %v", err)
	}

	client.SetToken(rootToken)
	data := make(map[string]interface{})
	data["test"] = "Hello Vault!"
	_, err = client.Logical().Write("secret/test", data)
	if err != nil {
		t.Errorf("Unable to write test value to vault: %v", err)
	}

	secret, err := client.Logical().Read("secret/test")
	if err != nil {
		t.Errorf("Unable to read test value from vault: %v", err)
	}

	if testString, ok := secret.Data["test"].(string); !ok || testString != "Hello Vault!" {
		t.Errorf("Wrong value returned from vault: %v", testString)
	}

}
