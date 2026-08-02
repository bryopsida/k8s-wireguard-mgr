package wireguard

import (
	"testing"
)

func TestGenerateWireguardKeyProducesUsableKeys(t *testing.T) {
	privateKey := GenerateWireguardKey()
	publicKey := privateKey.PublicKey()

	if privateKey.String() == "" {
		t.Fatal("private key should not be empty")
	}
	if publicKey.String() == "" {
		t.Fatal("public key should not be empty")
	}
	if privateKey.String() == publicKey.String() {
		t.Fatal("public key should not match the private key")
	}
}
