package processing

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/sumwonyuno/cp-scoring/model"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func ToBytes(state model.State, entities []*openpgp.Entity) ([]byte, error) {
	// state to JSON bytes
	bs, err := stateToJSON(state)
	if err != nil {
		return nil, err
	}
	// compress bytes
	bs, err = bytesCompress(bs)
	if err != nil {
		return nil, err
	}
	// encrypt bytes
	bs, err = bytesEncrypt(bs, entities)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func FromBytes(bs []byte, entities openpgp.EntityList) (model.State, error) {
	var state model.State

	if len(bs) == 0 {
		return state, errors.New("Empty bytes given")
	}

	// decrypt bytes
	bs, err := bytesDecrypt(bs, entities)
	if err != nil {
		return state, err
	}

	// decompress bytes
	bs, err = bytesDecompress(bs)
	if err != nil {
		return state, err
	}
	// JSON bytes to state
	state, err = stateFromJSON(bs)
	if err != nil {
		return state, err
	}

	return state, nil
}

func stateToJSON(state model.State) ([]byte, error) {
	bs, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func stateFromJSON(bs []byte) (model.State, error) {
	var state model.State
	err := json.Unmarshal(bs, &state)
	if err != nil {
		return state, err
	}

	return state, nil
}

func bytesCompress(bs []byte) ([]byte, error) {
	if bs == nil {
		return nil, errors.New("Cannot compress nil bytes")
	}
	if len(bs) == 0 {
		return nil, errors.New("Cannot compress 0 bytes")
	}

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(bs)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func bytesDecompress(bs []byte) ([]byte, error) {
	buf := bytes.NewBuffer(bs)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	outBuf := bytes.NewBuffer(nil)
	_, err = io.Copy(outBuf, r)
	if err != nil {
		return nil, err
	}

	if err := r.Close(); err != nil {
		return nil, err
	}

	return outBuf.Bytes(), nil
}

func bytesEncrypt(bs []byte, entities []*openpgp.Entity) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w, err := armor.Encode(buf, openpgp.PublicKeyType, nil)
	if err != nil {
		return nil, err
	}

	plaintext, err := openpgp.Encrypt(w, entities, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	plaintext.Write(bs)
	plaintext.Close()
	w.Close()

	return buf.Bytes(), nil
}

func bytesDecrypt(bs []byte, entities openpgp.EntityList) ([]byte, error) {
	buf := bytes.NewBuffer(bs)
	result, err := armor.Decode(buf)
	if err != nil {
		return nil, err
	}
	message, err := openpgp.ReadMessage(result.Body, entities, nil, nil)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(message.UnverifiedBody)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
