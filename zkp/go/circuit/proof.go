package circuit

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/conversion"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/emulated/emparams"
	"github.com/consensys/gnark/std/signature/ecdsa"
)

// VerifCircuit defines the circuit structure as provided by you.
type VerifCircuit[T, S emulated.FieldParams] struct {
	// PublicKeyHash is 32 bytes, but we split it into two 16-byte variables to fit into BN254 field
	PublicKeyHash frontend.Variable   `gnark:",public"`
	Msg           emulated.Element[S] `gnark:",public"`
	Signature     ecdsa.Signature[S]
	PublicKey     ecdsa.PublicKey[T, S] // actual public key is also a private input
	Nonce         frontend.Variable
}

func (c *VerifCircuit[T, S]) Define(api frontend.API) error {

	// Convert the public key into X and Y bytes
	xbytes, err := conversion.EmulatedToBytes(api, &c.PublicKey.X)
	if err != nil {
		return fmt.Errorf("failed to convert PublicKey.X to bytes: %w", err)
	}
	ybytes, err := conversion.EmulatedToBytes(api, &c.PublicKey.Y)
	if err != nil {
		return fmt.Errorf("failed to convert PublicKey.Y to bytes: %w", err)
	}

	// Compue MiMC hash as a public key commitment
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return fmt.Errorf("failed to create MiMC instance: %w", err)
	}
	x, err := conversion.BytesToNative(api, xbytes)
	if err != nil {
		return fmt.Errorf("failed to convert x into native: %w", err)
	}
	h.Write(x)
	y, err := conversion.BytesToNative(api, ybytes)
	if err != nil {
		return fmt.Errorf("failed to convert y into native: %w", err)
	}
	h.Write(y)
	h.Write(c.Nonce)
	computedHash := h.Sum()

	// Check the hash against the public commitment
	api.AssertIsEqual(c.PublicKeyHash, computedHash)

	// Verifying the ECDSA signature
	c.PublicKey.Verify(api, sw_emulated.GetCurveParams[emparams.Secp256k1Fp](), &c.Msg, &c.Signature)

	return nil
}

// Struct for JSON serialization of the transaction inputs.
type TransactionInput struct {
	MsgHash string `json:"msgHash"` // Hex string of the message hash
	R       string `json:"r"`       // Hex string of signature R
	S       string `json:"s"`       // Hex string of signature S
	PubX    string `json:"pubX"`    // Hex string of public key X
	PubY    string `json:"pubY"`    // Hex string of public key Y
}

// Struct for JSON serialization of the witnesses of the proof.
type ProofWitness struct {
	MsgHash string `json:"msgHash"` // Hex string of the message hash
	R       string `json:"r"`       // Hex string of signature R
	S       string `json:"s"`       // Hex string of signature S
	PubX    string `json:"pubX"`    // Hex string of public key X
	PubY    string `json:"pubY"`    // Hex string of public key Y
	Nonce   string `json:"Nonce"`   // Hex string of Nonce
	Com     string `json:"Com"`     // Hex string of Commitment
}

// readFromFile is a helper to deserialize and read gnark objects or JSON from files.
func ReadTransactionFromFile(filename string, data interface{}) ([]byte, []byte, []byte, []byte, []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error opening file %s: %w", filename, err)
	}
	defer file.Close()

	switch v := data.(type) {
	case *TransactionInput:
		decoder := json.NewDecoder(file)
		err = decoder.Decode(v)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("error decoding JSON from file %s: %w\n", filename, err)
		}
		// Decode hex strings back to big.Int and byte slices
		rBytes, err := hex.DecodeString(v.R)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("Error decoding R hex: %v\n", err)
		}
		sBytes, err := hex.DecodeString(v.S)
		if err != nil {
			return rBytes, nil, nil, nil, nil, fmt.Errorf("Error decoding S hex: %v\n", err)
		}
		msgHashBytes, err := hex.DecodeString(v.MsgHash)
		if err != nil {
			return rBytes, sBytes, nil, nil, nil, fmt.Errorf("Error decoding S hex: %v\n", err)
		}
		pubXBytes, err := hex.DecodeString(v.PubX)
		if err != nil {
			return rBytes, sBytes, msgHashBytes, nil, nil, fmt.Errorf("Error decoding S hex: %v\n", err)
		}
		pubYBytes, err := hex.DecodeString(v.PubY)
		if err != nil {
			return rBytes, sBytes, msgHashBytes, pubXBytes, nil, fmt.Errorf("Error decoding S hex: %v\n", err)
		}
		return rBytes, sBytes, msgHashBytes, pubXBytes, pubYBytes, nil

	default:
		return nil, nil, nil, nil, nil, fmt.Errorf("unsupported type for reading from file: %T", data)
	}

}

// readFromFile is a helper to deserialize and read gnark objects or JSON from files.
func ReadPublicInputFromFile(filename string, data interface{}) ([]byte, []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file %s: %w", filename, err)
	}
	defer file.Close()

	switch v := data.(type) {
	case *PublicInputJSON:
		decoder := json.NewDecoder(file)
		err = decoder.Decode(v)
		if err != nil {
			return nil, nil, fmt.Errorf("error decoding JSON from file %s: %w\n", filename, err)
		}
		// Decode hex strings back to big.Int and byte slices
		MsgHash, err := hex.DecodeString(v.MsgHash)
		if err != nil {
			return nil, nil, fmt.Errorf("Error decoding MsgHash hex: %v\n", err)
		}
		Com, err := hex.DecodeString(v.Com)
		if err != nil {
			return MsgHash, nil, fmt.Errorf("Error decoding Com hex: %v\n", err)
		}
		return MsgHash, Com, nil

	default:
		return nil, nil, fmt.Errorf("unsupported type for reading from file: %T", data)
	}

}
