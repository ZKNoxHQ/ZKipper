package utils

import (
	"fmt"
	"os"

	"github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func WriteSolidityTestFile(proof plonk.Proof, PublicKeyCommitment []byte, MsgHash []byte) error {

	// 9. Export the Solidity verifier test
	fmt.Println("\n--- Exporting Solidity Verifier Test ---")
	verifierTestFile, err := os.Create("../solidity/test/Verifier.t.sol")
	if err != nil {
		return fmt.Errorf("Error while creating the file Verifier.t.sol: %w\n", err)
	}
	defer verifierTestFile.Close()

	// header
	verifierTestFile.Write([]byte(`// SPDX-License-Identifier: UNLICENSED
	pragma solidity ^0.8.25;

	import {Test, console} from "forge-std/Test.sol";
	import {ZKVerifier} from "../src/ZKVerifier.sol";

	contract VerifierTest is Test {
	    ZKVerifier ZKV;

	    function setUp() public {
	        ZKV = new ZKVerifier();
	    }

	    function test_verify() public view {
	`))

	Proof := proof.(*plonk_bn254.Proof)
	verifierTestFile.Write([]byte(`bytes memory proof = hex"` + hexutil.Encode(Proof.MarshalSolidity())[2:] + `";`))
	verifierTestFile.Write([]byte("\n"))
	verifierTestFile.Write([]byte(`bytes memory public_key_commitment = hex"` + hexutil.Encode(PublicKeyCommitment)[2:] + `";`))
	verifierTestFile.Write([]byte("\n"))
	verifierTestFile.Write([]byte(`bytes memory transaction_hash = hex"` + hexutil.Encode(MsgHash)[2:] + `";`))
	verifierTestFile.Write([]byte("\n"))

	// footer
	verifierTestFile.Write([]byte(`
	        bool res = ZKV.Verify(proof, public_key_commitment, transaction_hash);
	        assertTrue(res);
	        console.log(res);
	    }
	}
	`))
	fmt.Println("Successfully exported solidty/test/Verifier.t.sol")

	fmt.Print("\n\n\n=======================\nPROOF and PUBLIC INPUTS\n=======================\n0x", hexutil.Encode(Proof.MarshalSolidity())[2:], ", \n", hexutil.Encode(PublicKeyCommitment), ", \n", hexutil.Encode(MsgHash), "\n")

	return nil
}
