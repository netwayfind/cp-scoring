package processing

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"testing"

	"golang.org/x/crypto/openpgp"

	"github.com/sumwonyuno/cp-scoring/model"

	_ "golang.org/x/crypto/ripemd160"
)

func getTestEntities(t *testing.T) []*openpgp.Entity {
	entity, err := newEntity()
	if err != nil {
		t.Fatal("Unable to create entity;", err)
	}
	entities := make([]*openpgp.Entity, 0)
	return append(entities, entity)
}

func compareStates(t *testing.T, state1 model.State, state2 model.State) {
	if state1.OS != state2.OS {
		t.Fatal("OS does not match")
	}
	if state1.Hostname != state2.Hostname {
		t.Fatal("Hostname does not match")
	}
	if state1.Timestamp != state2.Timestamp {
		t.Fatal("Timestamp does not match")
	}
	if len(state1.Users) != len(state2.Users) {
		t.Fatal("Users does not match")
	}
	if len(state1.Groups) != len(state2.Groups) {
		t.Fatal("Groups does not match")
	}
	if len(state1.Processes) != len(state2.Processes) {
		t.Fatal("Processes does not match")
	}
	if len(state1.Software) != len(state2.Software) {
		t.Fatal("Software does not match")
	}
	if len(state1.NetworkConnections) != len(state2.NetworkConnections) {
		t.Fatal("Network connections does not match")
	}
	if len(state1.Errors) != len(state2.Errors) {
		t.Fatal("Errors does not match")
	}
}

func TestStateToJSON(t *testing.T) {
	// no state
	var state model.State
	bs, err := stateToJSON(state)
	if err != nil {
		t.Fatal("Error serializing state to JSON;", err)
	}
	if len(bs) == 0 {
		t.Fatal("No bytes for serializing state")
	}

	var state2 model.State

	err = json.Unmarshal(bs, &state2)
	if err != nil {
		t.Fatal("Bytes not JSON;", err)
	}

	// new state
	state = model.GetNewStateTemplate()
	bs, err = stateToJSON(state)
	if err != nil {
		t.Fatal("Error serializing state to JSON;", err)
	}
	if len(bs) == 0 {
		t.Fatal("No bytes for serializing state")
	}

	// check that result is JSON bytes
	err = json.Unmarshal(bs, &state2)
	if err != nil {
		t.Fatal("Bytes not JSON;", err)
	}
}

func TestStateFromJSON(t *testing.T) {
	_, err := stateFromJSON(nil)
	if err == nil {
		t.Fatal("Parsed state from nil bytes;", err)
	}

	// empty string
	bs := []byte("")
	_, err = stateFromJSON(bs)
	if err == nil {
		t.Fatal("Parsed state from empty string;", err)
	}

	// convert from state to bytes and back
	state := model.GetNewStateTemplate()
	bs, err = stateToJSON(state)
	if err != nil {
		t.Fatal("Unable to serialize state to JSON;", err)
	}
	state2, err := stateFromJSON(bs)
	if err != nil {
		t.Fatal("Unable to deserialize state from JSON;", err)
	}
	compareStates(t, state, state2)
}

func TestBytesCompress(t *testing.T) {
	// nil
	_, err := bytesCompress(nil)
	if err == nil {
		t.Fatal("Compressed nil bytes;", err)
	}

	// empty bytes
	bs := []byte("")
	_, err = bytesCompress(bs)
	if err == nil {
		t.Fatal("Compressed empty string;", err)
	}

	// JSON bytes
	orig := []byte("{\"test\": \"test\"}")
	newBytes, err := bytesCompress(orig)
	if err != nil {
		t.Fatal("Compressed empty string;", err)
	}
	if len(newBytes) == 0 {
		t.Fatal("Expected non-zero length for compressed bytes")
	}

	// check that result is gzip bytes
	buf := bytes.NewBuffer(newBytes)
	r, err := gzip.NewReader(buf)
	if err != nil {
		t.Fatal("Unable to create gzip reader;", err)
	}
	outBuf := bytes.NewBuffer(nil)
	_, err = io.Copy(outBuf, r)
	if err != nil {
		t.Fatal("Bytes not gzip;", err)
	}
	if bytes.Compare(orig, outBuf.Bytes()) != 0 {
		t.Fatal("Decompressed bytes does not match original")
	}
}

func TestBytesDecompress(t *testing.T) {
	// nil
	_, err := bytesDecompress(nil)
	if err == nil {
		t.Fatal("Decompressed nil bytes;", err)
	}

	// empty bytes
	bs := []byte("")
	_, err = bytesDecompress(bs)
	if err == nil {
		t.Fatal("Decompressed empty bytes;", err)
	}

	// convert from compressed to decompressed
	orig := []byte("000000000000000")
	bs, err = bytesCompress(orig)
	if err != nil {
		t.Fatal("Unable to compress bytes;", err)
	}
	bs, err = bytesDecompress(bs)
	if err != nil {
		t.Fatal("Unable to decompress bytes;", err)
	}
	if bytes.Compare(orig, bs) != 0 {
		t.Fatal("Decompressed bytes does not match original")
	}
}

func TestBytesEncrypt(t *testing.T) {
	entities := getTestEntities(t)

	// nil
	bs, err := bytesEncrypt(nil, entities)
	if err != nil {
		t.Fatal("Unable to encrypt nil bytes")
	}
	if len(bs) == 0 {
		t.Fatal("No bytes from encrypting nil bytes")
	}

	// empty bytes
	bs = []byte("")
	bs, err = bytesEncrypt(bs, entities)
	if err != nil {
		t.Fatal("Unable to encrypt empty string")
	}
	if len(bs) == 0 {
		t.Fatal("No bytes from encrypting empty bytes")
	}

	// example
	bs = []byte("example")
	bs, err = bytesEncrypt(bs, entities)
	if err != nil {
		t.Fatal("Unable to encrypt bytes;", err)
	}
	if len(bs) == 0 {
		t.Fatal("No bytes from encrypting bytes")
	}
}

func TestBytesDecrypt(t *testing.T) {
	entities := getTestEntities(t)

	// nil
	bs, err := bytesEncrypt(nil, entities)
	if err != nil {
		t.Fatal("Unable to encrypt;", err)
	}
	bs, err = bytesDecrypt(bs, entities)
	if err != nil {
		t.Fatal("Unable to decrypt;", err)
	}
	if len(bs) != 0 {
		t.Fatal("Unexpected decrypted bytes;", bs)
	}

	// empty bytes
	orig := []byte("")
	bs, err = bytesEncrypt(orig, entities)
	if err != nil {
		t.Fatal("Unable to encrypt;", err)
	}
	bs, err = bytesDecrypt(bs, entities)
	if err != nil {
		t.Fatal("Unable to decrypt;", err)
	}
	if bytes.Compare(orig, bs) != 0 {
		t.Fatal("Unexpected decrypted bytes;", bs)
	}

	// example
	orig = []byte("example")
	bs, err = bytesEncrypt(orig, entities)
	if err != nil {
		t.Fatal("Unable to encrypt;", err)
	}
	bs, err = bytesDecrypt(bs, entities)
	if err != nil {
		t.Fatal("Unable to decrypt;", err)
	}
	if bytes.Compare(orig, bs) != 0 {
		t.Fatal("Unexpected decrypted bytes;", bs)
	}
}

func TestToBytes(t *testing.T) {
	entities := getTestEntities(t)

	state := model.GetNewStateTemplate()
	bs, err := ToBytes(state, entities)
	if err != nil {
		t.Fatal("Unable to convert state to bytes;", err)
	}
	if len(bs) == 0 {
		t.Fatal("No bytes created")
	}
}

func TestFromBytes(t *testing.T) {
	entities := getTestEntities(t)

	state := model.GetNewStateTemplate()
	bs, err := ToBytes(state, entities)
	if err != nil {
		t.Fatal("Unable to convert state to bytes;", err)
	}
	state2, err := FromBytes(bs, entities)
	if err != nil {
		t.Fatal("Unable to convert bytes to state;", err)
	}
	compareStates(t, state, state2)
}
