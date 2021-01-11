package processing

import (
	"bytes"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func newEntity() (*openpgp.Entity, error) {
	entity, err := openpgp.NewEntity("cp-scoring", "test", "test@example.com", nil)
	if err != nil {
		return nil, err
	}
	for _, id := range entity.Identities {
		if err := id.SelfSignature.SignUserId(id.UserId.Id, entity.PrimaryKey, entity.PrivateKey, nil); err != nil {
			return nil, err
		}
	}

	return entity, nil
}

func NewPubPrivKeys() ([]byte, []byte, error) {
	entity, err := newEntity()
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := GetPubKey(entity)
	if err != nil {
		return nil, nil, err
	}
	privKey, err := GetPrivKey(entity)
	if err != nil {
		return nil, nil, err
	}

	return pubKey, privKey, err
}

func GetPubKey(entity *openpgp.Entity) ([]byte, error) {
	bufPub := bytes.NewBuffer(nil)
	writerPub, err := armor.Encode(bufPub, openpgp.PublicKeyType, nil)
	if err != nil {
		return nil, err
	}
	if err = entity.Serialize(writerPub); err != nil {
		return nil, err
	}
	if err = writerPub.Close(); err != nil {
		return nil, err
	}

	return bufPub.Bytes(), nil
}

func GetPrivKey(entity *openpgp.Entity) ([]byte, error) {
	bufPriv := bytes.NewBuffer(nil)
	writerPriv, err := armor.Encode(bufPriv, openpgp.PrivateKeyType, nil)
	if err != nil {
		return nil, err
	}
	if err = entity.SerializePrivate(writerPriv, nil); err != nil {
		return nil, err
	}
	if err = writerPriv.Close(); err != nil {
		return nil, err
	}

	return bufPriv.Bytes(), nil
}
