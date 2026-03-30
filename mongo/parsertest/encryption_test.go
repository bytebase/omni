package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestEncryptionGetKeyVault(t *testing.T) {
	node := mustParse(t, `db.getMongo().getKeyVault()`)
	stmt := node.(*ast.EncryptionStatement)
	if stmt.Target != "keyVault" {
		t.Errorf("expected target 'keyVault', got %q", stmt.Target)
	}
	if len(stmt.ChainedMethods) != 0 {
		t.Errorf("expected 0 chained methods, got %d", len(stmt.ChainedMethods))
	}
}

func TestEncryptionGetKeyVaultGetKeys(t *testing.T) {
	node := mustParse(t, `db.getMongo().getKeyVault().getKeys()`)
	stmt := node.(*ast.EncryptionStatement)
	if stmt.Target != "keyVault" {
		t.Errorf("expected target 'keyVault', got %q", stmt.Target)
	}
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "getKeys" {
		t.Errorf("expected 'getKeys', got %q", stmt.ChainedMethods[0].Name)
	}
}

func TestEncryptionCreateKey(t *testing.T) {
	node := mustParse(t, `db.getMongo().getKeyVault().createKey("local")`)
	stmt := node.(*ast.EncryptionStatement)
	if stmt.Target != "keyVault" {
		t.Errorf("expected target 'keyVault', got %q", stmt.Target)
	}
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "createKey" {
		t.Errorf("expected 'createKey', got %q", stmt.ChainedMethods[0].Name)
	}
	if len(stmt.ChainedMethods[0].Args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(stmt.ChainedMethods[0].Args))
	}
}

func TestEncryptionGetClientEncryption(t *testing.T) {
	node := mustParse(t, `db.getMongo().getClientEncryption()`)
	stmt := node.(*ast.EncryptionStatement)
	if stmt.Target != "clientEncryption" {
		t.Errorf("expected target 'clientEncryption', got %q", stmt.Target)
	}
}

func TestEncryptionEncrypt(t *testing.T) {
	node := mustParse(t, `db.getMongo().getClientEncryption().encrypt(UUID("key-id"), "secret", "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic")`)
	stmt := node.(*ast.EncryptionStatement)
	if stmt.Target != "clientEncryption" {
		t.Errorf("expected target 'clientEncryption', got %q", stmt.Target)
	}
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "encrypt" {
		t.Errorf("expected 'encrypt', got %q", stmt.ChainedMethods[0].Name)
	}
}

func TestEncryptionDecrypt(t *testing.T) {
	node := mustParse(t, `db.getMongo().getClientEncryption().decrypt(BinData(6, "base64data"))`)
	stmt := node.(*ast.EncryptionStatement)
	if stmt.Target != "clientEncryption" {
		t.Errorf("expected target 'clientEncryption', got %q", stmt.Target)
	}
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "decrypt" {
		t.Errorf("expected 'decrypt', got %q", stmt.ChainedMethods[0].Name)
	}
}

func TestEncryptionCreateKeyWithKMS(t *testing.T) {
	node := mustParse(t, `db.getMongo().getKeyVault().createKey("aws", { region: "us-east-1", key: "arn:aws:kms:..." })`)
	stmt := node.(*ast.EncryptionStatement)
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if len(stmt.ChainedMethods[0].Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(stmt.ChainedMethods[0].Args))
	}
}
