package processing

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"testing"

	"github.com/sumwonyuno/cp-scoring/model"
)

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
	// compare states
	if state.OS != state2.OS {
		t.Fatal("OS does not match")
	}
	if state.Hostname != state2.Hostname {
		t.Fatal("OS does not match")
	}
	if state.Timestamp != state2.Timestamp {
		t.Fatal("OS does not match")
	}
	if len(state.Users) != len(state2.Users) {
		t.Fatal("OS does not match")
	}
	if len(state.Groups) != len(state2.Groups) {
		t.Fatal("OS does not match")
	}
	if len(state.Processes) != len(state2.Processes) {
		t.Fatal("OS does not match")
	}
	if len(state.Software) != len(state2.Software) {
		t.Fatal("OS does not match")
	}
	if len(state.NetworkConnections) != len(state2.NetworkConnections) {
		t.Fatal("OS does not match")
	}
	if len(state.Errors) != len(state2.Errors) {
		t.Fatal("OS does not match")
	}
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
