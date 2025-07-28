package main

import (
	// "bytes"
	// "crypto/rand"
	// "encoding/hex"
	// "encoding/json"
	// "fmt"
	// "math/big"
	// "os"
	// "time"

	// "github.com/consensys/gnark-crypto/ecc"
	// cryptomimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	// "github.com/consensys/gnark/backend/plonk"
	// "github.com/consensys/gnark/frontend"
	// "github.com/consensys/gnark/std/math/emulated"
	// "github.com/consensys/gnark/std/signature/ecdsa"

	"fmt"
	"log"
	"os"
	"time"
	"zkp/circuit"
	"zkp/utils"

	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
)

func main() {

	withSetup := false
	withWitness := false
	withExportSolidity := false

	for _, arg := range os.Args[1:] {
		if arg == "setup" {
			withSetup = true
			break
		}
		if arg == "witness" {
			withWitness = true
			break
		}
		if arg == "export-solidity" {
			withExportSolidity = true
			break
		}
	}

	var R1CS constraint.ConstraintSystem
	var PK plonk.ProvingKey
	var VK plonk.VerifyingKey
	var err error

	if withSetup {
		// Trusted setup
		R1CS, PK, VK, err = circuit.ComputeTrustedSetup()
		if err != nil {
			panic(err)
		}

	} else {
		// Setup is already computed
		R1CS, PK, VK, err = circuit.LoadTrustedSetup()
		if err != nil {
			panic(err)
		}

		var witness, publicInput witness.Witness

		if withWitness {
			witness, publicInput, err = circuit.GenerateWitness()
			if err != nil {
				log.Fatalf("Failed to generate the witness: %v", err)
				panic(err)
			}
			fmt.Println("Witness generated.")
		} else {
			// Load the witness and public input and compute the proof
			// Witness
			fmt.Println("Loading witness and public input")
			witness, publicInput, err = circuit.LoadWitnessAndPublicInput()
			if err != nil {
				log.Fatalf("Failed to load the witness and public input: %v", err)
				panic(err)
			}
			fmt.Println("Witness loaded.")

			// Proof
			fmt.Println("\n--- Proving with loaded setup ---")
			startProve := time.Now()
			proof, err := plonk.Prove(R1CS, PK, witness)
			fmt.Printf("Proof GENERATED (%.1fms).\n", float64(time.Since(startProve).Milliseconds()))

			// Verify the proof
			err = plonk.Verify(proof, VK, publicInput)
			if err != nil {
				fmt.Printf("Error verifying the proof: %v\n", err)
				os.Exit(1)
			}

			if withExportSolidity {
				// Read output/public_input.json
				var PublicInput circuit.PublicInputJSON
				MsgHash, Com, err := circuit.ReadPublicInputFromFile("output/public_input.json", &PublicInput)
				if err != nil {
					fmt.Printf("Error while reading output/public_input.json: %v\n", err)
					os.Exit(1)
				}
				// Export solidity test file
				utils.WriteSolidityTestFile(proof, Com, MsgHash)
			}
		}
	}

}
