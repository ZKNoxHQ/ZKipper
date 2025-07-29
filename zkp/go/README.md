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
It creates a file `output/scs.bin` containing the setup, and also the corresponding Solidity contract `../solidity/src/Verifier.sol`.

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
0x2bd9a350cb8951a39dfeb9008378bab56241eeb1cae3b95bcbbf1ace6e9b997c298b1ae449cf87b2da2eeccbeac65038d8df68d68b25feb9a65bd78ce61033d42cec749262648918b18f86acbfa2e0327e479802a4e426421058794791ee64b50b397ed6b96de794e82dddc3b32fe64be200d69db5161bcfcae967acdd121a72285f55d051f4e3fe69b32eeac3c2d23cb525ef4914517b6584efebf357897cfa19f1e87625e32453fd402e6df911abe704bf665d25d15fb5dc4a960803be61ba0b0e1d39f8bb0476dc9abe3259d8f58649433c777f22a1ce0909c9dc84c1b7f0288a23b313b1869f0967ccbe1b6ae6c093ecd84761b57ca0805a39462ffa1e0f1b129f9b7ea84fe9ea81a35f88c15addfc79240d980dd3237a7f103ee4821f4203fca269dc7b47ee695f79308983c27329a92063c740a27e2ada549d7b566e361a589a701c62650df95ec4af8921d47cc67ee5bb88ad5f8c3c8fb5a4e69091791938bfbfffe1cde611f0bf0a97e5ce856fd2b4222159fcb66dd005c8b9dbc51c1aacd50da9c271138301d1192b0621381cfe0b179111af9975984b15578cb6291772d95d4caf799976600ab06349e7de2bf22461aef241458461b5e483a20acb05fda43beba4ebd67c97479fc4b80f4245c7010d0ab904eadca6282db764e7080cf4045b15db14083492e60c585b397d9438cf2bd9e85bd2ba64936b283e58460478098ce53cacf285a90a8488ea226b2ac0748d2b85ec534f2affde6a16c2272b68d309220ca66cc45d9ca51c8f696883c80b21a29990b800747af43a0bbbe126ca8089fe8ac538ac0f9ce66913cfe03a6ee9963ddd522a5b6bf864f690cdfa1438aa8596e16301d0744a2cc463cbfce79b4feafcdfee4d2c39cf6600d1af4b1d7acf5ffc2a736c7ed6d924d7356923e2c51d56041ca6f400296f7efddfd6950e3913d18cabb6ac8f85dc408d10d506d8c37bda39d6f4270ccc0941d3c44a901d935f18d56a3081f85cafbd9262b2fefc5c9dcd15318e8c4ae7c6cdcdcc20b8208b1a0d92597ae7070be48ec98002a30895089221d2395b14944d0f075f01d71a52b1f1bf60764c9cd1f80fa04b9c9517f0a47f0482f8c06ce0092a3fe1e39514674d50dc2774205c5e5aac028fc35e3b6aed8ae3b5fbfb7c60c07106d51dbc11dcdec9f3b8fbe5c721d35f79db738049002ff856e9f3f66c9a6769d32ece0a, 
0x26ecce5184a1a1135bb87c27bc7359967668313a5d12a36275e045620ae1932d, 
0x98dd76b826575fbb239b926fd1d49aabebd56bf7810a7d886788d7740c28f9de
```

## Verification
The proof can be verified using the solidity contract. It can be checked with:
```
cd ../solidty/
forge test -vvv
```


