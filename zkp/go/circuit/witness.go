package circuit

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"zkp/utils"

	"github.com/consensys/gnark-crypto/ecc"
	cryptomimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/signature/ecdsa"
)

type WitnessJSON struct {
	R     string `json:"r"`     // Hex string of r
	S     string `json:"s"`     // Hex string of s
	PubX  string `json:"pubX"`  // Hex string of public key X
	PubY  string `json:"pubY"`  // Hex string of public key Y
	Nonce string `json:"nonce"` // Hex string of nonce
}

type PublicInputJSON struct {
	MsgHash string `json:"msgHash"` // Hex string of msgHash
	Com     string `json:"com"`     // Hex string of Com
}

func GenerateWitness() (witness.Witness, witness.Witness, error) {
	// Load the signature, public key and message hash from the json file
	// Generate a nonce and the corresponding public key commitment
	// Save it to a file

	EmptyCircuit := VerifCircuit[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	EmptyWitness, err := frontend.NewWitness(&EmptyCircuit, ecc.BN254.ScalarField())

	// Read the JSON file
	var TransactionInput TransactionInput
	rBytes, sBytes, msgHashBytes, pubXBytes, pubYBytes, err := ReadTransactionFromFile("input/transaction_input.json", &TransactionInput)
	if err != nil {
		return EmptyWitness, EmptyWitness, fmt.Errorf("Error while reading input/transaction_input.json")
	}

	// Nonce (160 bits)
	nonceBytes := make([]byte, 20)
	_, err = rand.Read(nonceBytes)
	if err != nil {
		return EmptyWitness, EmptyWitness, fmt.Errorf("Error while creating the nonce: %w\n", err)
	}

	// compute the hash of the public key
	hnew := cryptomimc.NewMiMC()
	_, err = hnew.Write(pubXBytes)
	if err != nil {
		panic(err)
	}
	_, err = hnew.Write(pubYBytes)
	if err != nil {
		panic(err)
	}
	_, err = hnew.Write(nonceBytes)
	if err != nil {
		panic(err)
	}
	pubHash := hnew.Sum(nil)

	// now we prepare the witness for the circuit
	witnessCircuit := VerifCircuit[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		// we splitted the public key hash into two 16-byte variables to fit into BN254 field
		PublicKeyHash: frontend.Variable(pubHash),
		// we construct the public key as non-native element. NB! this means that both X and Y coordinates are 4 limbs of 64 bytes each, so 8 limbs total
		PublicKey: ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](pubXBytes),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](pubYBytes),
		},
		Signature: ecdsa.Signature[emulated.Secp256k1Fr]{
			R: emulated.ValueOf[emulated.Secp256k1Fr](rBytes),
			S: emulated.ValueOf[emulated.Secp256k1Fr](sBytes),
		},
		Msg:   emulated.ValueOf[emulated.Secp256k1Fr](msgHashBytes),
		Nonce: frontend.Variable(nonceBytes),
	}

	witnessFull, err := frontend.NewWitness(&witnessCircuit, ecc.BN254.ScalarField())
	if err != nil {
		log.Fatalf("failed to create witnessFull: %v", err)
	}
	publicInput, err := witnessFull.Public()
	if err != nil {
		log.Fatalf("failed to get public witness: %v", err)
	}

	// Save PublicInput to a file
	PublicInput := PublicInputJSON{
		MsgHash: hex.EncodeToString(msgHashBytes),
		Com:     hex.EncodeToString(pubHash),
	}
	OutputJSON, err := json.MarshalIndent(PublicInput, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling prove public input JSON: %v\n", err)
		os.Exit(1)
	}
	utils.WriteToFile("output/public_input.json", bytes.NewReader(OutputJSON))
	utils.WriteToFile("output/public_input.bin", publicInput)

	// Save Witness to a file
	Witness := WitnessJSON{
		R:     hex.EncodeToString(rBytes),
		S:     hex.EncodeToString(sBytes),
		PubX:  hex.EncodeToString(pubXBytes),
		PubY:  hex.EncodeToString(pubYBytes),
		Nonce: hex.EncodeToString(nonceBytes),
	}
	OutputJSON, err = json.MarshalIndent(Witness, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling prove witness JSON: %v\n", err)
		os.Exit(1)
	}
	utils.WriteToFile("output/witness.json", bytes.NewReader(OutputJSON))
	utils.WriteToFile("output/witness.bin", witnessFull)

	return witnessFull, publicInput, nil
}

func LoadWitnessAndPublicInput() (witness.Witness, witness.Witness, error) {
	// Witness
	Witness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to allocate witness: %w", err)
	}
	err = utils.ReadFromFile("output/witness.bin", Witness)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read output/witness.bin: %w", err)
	}
	fmt.Printf("Read output/witness.bin\n")

	// Public Input
	PublicInput, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to allocate public input: %w", err)
	}
	err = utils.ReadFromFile("output/public_input.bin", PublicInput)
	if err != nil {
		return Witness, nil, fmt.Errorf("failed to read output/public_input.bin: %w", err)
	}
	fmt.Printf("Read output/public_input.bin\n")

	return Witness, PublicInput, nil
}
