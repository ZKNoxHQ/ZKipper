package circuit

import (
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc"

	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/std/math/emulated"

	"github.com/consensys/gnark/test/unsafekzg"

	"zkp/utils"
)

func ComputeTrustedSetup() (constraint.ConstraintSystem, plonk.ProvingKey, plonk.VerifyingKey, error) {
	// This function runs a trusted setup
	// and save the files:
	// - `scs.bin`,
	// - `proving_key.bin`,
	// - `verifying_key.bin`.
	fmt.Println("--- Generating circuit inputs and performing compliance check ---")

	// 1. Compile the circuit
	circuit := VerifCircuit[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	fmt.Printf("Compiling circuit...\n")
	SCS, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &circuit)
	if err != nil {
		fmt.Printf("Error compiling ECDSA circuit: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("BN254 circuit compiled with %d constraints",
		SCS.GetNbConstraints())

	// 2. Perform Groth16 setup
	fmt.Printf("Starting Plonk setup...\n")
	A, B, _ := unsafekzg.NewSRS(SCS)
	PK, VK, err := plonk.Setup(SCS, A, B)
	if err != nil {
		fmt.Printf("Error during Plonk setup for ECDSA: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Setup done.\n")

	// 3. Save to bin files
	utils.WriteToFile("output/scs.bin", SCS)
	utils.WriteToFile("output/proving_key.bin", PK)
	utils.WriteToFile("output/verifying_key.bin", VK)

	fmt.Println("\nAll input files generated successfully for CGO wrapper.")

	// 4. Export the Solidity verifier contract
	fmt.Println("\n--- Exporting Solidity Verifier ---")
	verifierFile, err := os.Create("../solidity/src/Verifier.sol")
	if err != nil {
		fmt.Printf("Error creating solidity/src/Verifier.sol: %v\n", err)
		os.Exit(1)
	}
	defer verifierFile.Close()
	err = VK.ExportSolidity(verifierFile)
	if err != nil {
		fmt.Printf("Error exporting solidity verifier: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully exported solidty/src/Verifier.sol")

	return SCS, PK, VK, err
}

func LoadTrustedSetup() (constraint.ConstraintSystem, plonk.ProvingKey, plonk.VerifyingKey, error) {
	// 1. Read back the SCS
	loadedSCS := plonk.NewCS(ecc.BN254)
	err := utils.ReadFromFile("output/scs.bin", loadedSCS)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read scs.bin: %w", err)
	}
	fmt.Printf("Read scs.bin (Constraints: %d)\n", loadedSCS.GetNbConstraints())

	// 2. Read back the proving key
	loadedPK := plonk.NewProvingKey(ecc.BN254)
	err = utils.ReadFromFile("output/proving_key.bin", loadedPK)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read proving_key.bin: %w", err)
	}
	fmt.Println("Read proving_key.bin")

	// 3. Read back the verifying key
	loadedVK := plonk.NewVerifyingKey(ecc.BN254)
	err = utils.ReadFromFile("output/verifying_key.bin", loadedVK)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read verifying_key.bin: %w", err)
	}
	fmt.Println("Read verifying_key.bin")

	return loadedSCS, loadedPK, loadedVK, nil
}
