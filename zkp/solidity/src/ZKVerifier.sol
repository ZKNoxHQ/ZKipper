// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {PlonkVerifier} from "./Verifier.sol";

contract ZKVerifier {
    PlonkVerifier PlonkV;

    constructor() {
        PlonkV = new PlonkVerifier();
    }

    function Verify(bytes calldata proof, bytes calldata public_key_commitment, bytes calldata transaction_hash)
        public
        view
        returns (bool success)
    {
        // public input contains:
        // - the public key commitment (one uint256)
        // - four uint256 containing 64 bits of the transaction hash each (starting from least significant 64-bit word)

        uint256[] memory public_input = new uint256[](5);
        // public key commitment
        public_input[0] = uint256(bytes32(public_key_commitment));
        // transaction hash
        uint256 mask = 0xFFFFFFFFFFFFFFFF; // 64-bit mask
        bytes32 thash = bytes32(transaction_hash);
        public_input[1] = uint256(thash) & mask; // Least significant 64 bits
        public_input[2] = uint256((thash >> 64)) & mask;
        public_input[3] = uint256((thash >> 128)) & mask;
        public_input[4] = uint256((thash >> 192)) & mask; // Most significant 64 bits

        // calling the plonk verifier verification
        return PlonkV.Verify(proof, public_input);
    }
}
