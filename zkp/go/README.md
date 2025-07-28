# Zero-knowledge proof of transaction

The ZK proof is computed using GNARK with PLONK proof system. It requires a trusted setup that has to be computed **once**.
The proof takes a few seconds to be generated and can be verified on-chain with our Solidity contract.

## Initialization
In order to initialize the go package:
```bash
make init
```

## Formatting the transaction information
In order to format the transaction information:
```bash
./get_transaction <msgHash> <r> <s> <pubX> <pubY>
```
An example with a working transaction:
```bash
./get_transaction \
98dd76b826575fbb239b926fd1d49aabebd56bf7810a7d886788d7740c28f9de \
364159a43e2ed67f690d16092933f84ceaeaa5bc389a96a076d1b2f78bcd5cad \
4409a0a1d98d739d57b558367a792dadd9de30836adc54edc12ae11f3a80cc6d \
0463c3596ceb913203c8e124371024f4874323e0d6a54928acc56d656d18e93a \
2623e23651c1e5a3c477ad87e978aaa7beed236ba833cb53ca32ae4af6051b09
```
This creates a file `input/transaction_input.json`.

## Trusted setup
:warning: Defining a new trusted setup requires updating the on-chain contracts. This can be done using:
```bash
make setup
```
It creates a file `output/r1cs.bin` containing the setup, and also the corresponding Solidity contract `../solidity/src/Verifier.sol`.

## Witness generation
In order to generate a witness (i.e. a random nonce and the corresponding public key commitment):
```bash
make witness
```

## Generating the proof
After having the transaction in the right format, you can generate the zero-knowledge proof:
```bash
make proof
```

## Export to Solidity
It is possible to export a verifier test file `../solidity/test/Verifier.t.sol`.

It also outputs the proof together with the public inputs. Here is an example:
```
=======================
PROOF and PUBLIC INPUTS
=======================
0x1eb0ce1d9f444f8444ad11ebd60cd61885741dd874029bd2d19a61f23e99bc2e26212c6f89c46673d3d4fd2fa2789fb64b1352d6a5ce0d25c52d9ed859485ea1086f76aef6f277089aeb1ea44f982a730448b6e24d99ec413d91929da10e40bc0a004b2e8e67d0f9d745068a865d4aad92e1a8bd559d913ebf93eb7863f0cb1f2bb9f65c26f1ff017f3b32e6aa4aa3375fe522880624a766c79932a5de49277f1ca9869ebd46235d21fef9186924b10e90891f884ab6c125574cc1ab933f81df2093c2ee160a38d921835d1b052470082a859fa30dd8b6ebce1546562f5242511f6aadf1a953b59f5f3c0c306789f5da628c9ee06eefdbf486fdf5671183aea61ffa5f46c140525fb04667c21e45590fca9b3b11977ae8e57f80d8fcfd6be4e62a48436968d46639df198c83706252b6167db603310009ca195a2fe0f7dc5cf828b9abd461008eca031c5d868782b75cfcd8b027098c34d525bdf9d4ab6622422d74fca2fdfbedee75acf11f73933054984a3bdd513dd317c1ab6365a775cd0b06a8f87700e83d32f8e7effa008dff5cf7e676a94e1090f8405e6f9256e4e88c0c647cb3d9ef7375b13b0d20874963a51c37b44d51960531e52c0b0f5b3aaa921efaf55d80741931b7b3172864912ffaca6f4f204d57f14e997b72b7b65d930116c957c1406b43f07f178776ac5e06acd31fcb92d294593395044832bb5e2b031d8d9d926f2ab82e66fe986632339a01dea62607bb3d408d24a35a84e23f097d0a0a53b4a84b83167f4042a778ed04b95b4afaf22af0700327163d1d2f21fc381da1e7e00bc5fe2e4a08a4de471b7f9d89b4b0f3db53cb8b563fcff4c367a33a2013efa8653d1de97af8bd102cc173956e172e67e8cf5bc843a649febfbeb2b93031becb0fa8633a14eaa8e40e81c0307af4fee46b4c606a9c5f3f006512d8a615bf724f264d5fe78d654b4b24fb1689ac975c65132a80eb27fdc56d27b0612529b49919ad40b5f4a4ee8e13588ea20ec250b79eb88331fb1d669f052f78712a002cc47090b292d82328938d58106ba387d2eec86752bf8921630b10e94c94eb1e5f828ba75f31399f8a7100bc8c6857fcbc7b04951ea1a8b13801404d529db62c7ccd2beca9e3f8602b4f2b8849555e15b521b70125b800300509af86fea3631768b4e9c8498c9005728b3a7596a204832309069fdfb7fd87517e5e1cc1b5c7, 0x98dd76b826575fbb239b926fd1d49aabebd56bf7810a7d886788d7740c28f9de, 0x98dd76b826575fbb239b926fd1d49aabebd56bf7810a7d886788d7740c28f9de
```

## Verification
The proof can be verified using the solidity contract. It can be checked with:
```
cd ../solidty/
forge test -vvv
```


