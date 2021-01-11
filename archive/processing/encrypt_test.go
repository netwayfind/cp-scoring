package processing

import (
	"testing"
)

func TestNewEntity(t *testing.T) {
	_, err := newEntity()
	if err != nil {
		t.Fatal("Unable to create entity;", err)
	}
}

func TestGetPubKey(t *testing.T) {
	entity, err := newEntity()
	if err != nil {
		t.Fatal("Unable to create entity;", err)
	}

	pubKey, err := GetPubKey(entity)
	if err != nil {
		t.Fatal("Unable to get public key;", err)
	}
	if len(pubKey) == 0 {
		t.Fatal("No public key found")
	}
}

func TestGetPrivKey(t *testing.T) {
	entity, err := newEntity()
	if err != nil {
		t.Fatal("Unable to create entity;", err)
	}

	privKey, err := GetPrivKey(entity)
	if err != nil {
		t.Fatal("Unable to get private key;", err)
	}
	if len(privKey) == 0 {
		t.Fatal("No private key found")
	}
}

func TestNewPubPrivKeys(t *testing.T) {
	pubKey, privKey, err := NewPubPrivKeys()
	if err != nil {
		t.Fatal("Unable to get new public and private keys;", err)
	}
	if len(pubKey) == 0 {
		t.Fatal("Got empty public key")
	}
	if len(privKey) == 0 {
		t.Fatal("Got empty private key")
	}
}
