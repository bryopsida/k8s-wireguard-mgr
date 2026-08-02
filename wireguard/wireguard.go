package wireguard

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type Key [32]byte

func GenerateWireguardKey() Key {
	key, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	var out Key
	copy(out[:], key.Bytes())
	return out
}

func (k Key) PublicKey() Key {
	privateKey, err := ecdh.X25519().NewPrivateKey(k[:])
	if err != nil {
		panic(err)
	}
	pub := privateKey.PublicKey()
	var out Key
	copy(out[:], pub.Bytes())
	return out
}

func (k Key) String() string {
	return base64.StdEncoding.EncodeToString(k[:])
}

func ParseKey(s string) (Key, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return Key{}, fmt.Errorf("wireguard: parse key: %w", err)
	}
	if len(decoded) != len(Key{}) {
		return Key{}, fmt.Errorf("wireguard: invalid key length: %d", len(decoded))
	}
	var out Key
	copy(out[:], decoded)
	return out, nil
}
