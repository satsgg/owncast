package l402

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"

	"github.com/lightningnetwork/lnd/lntypes"
)

func CreateUniqueIdentifier(paymentHash [32]byte) ([]byte, error) {
	var tokenID [32]byte
	_, err := rand.Read(tokenID[:])
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	var version uint16 = 0
	if err := binary.Write(&buf, binary.BigEndian, version); err != nil {
		return nil, err
	}
	if _, err := buf.Write(paymentHash[:]); err != nil {
		return nil, err
	}
	if _, err := buf.Write(tokenID[:]); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodePaymentHash(r io.Reader) (*lntypes.Hash, error) {
	var version uint16
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return nil, err
	}
	var paymentHash lntypes.Hash
	if _, err := r.Read(paymentHash[:]); err != nil {
		return nil, err
	}
	return &paymentHash, nil
}