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
0x28b48077440f4e35cc0e016b2fd82cd033c51878c723e5b4d07eccde88a04b781b69068a3bc78bfc61b7841ed736cc3da873c412dc7c2d6fe840aefa354a9900122499906cb628cf604001999e78c33c0e09e76df67b20d9d6ffd14d29849453132fd12dfe29501b5d3c25d33927ce031258d6cd18c226bc0f1fd8360dbe8804276795f22ba16a0fa80e450e03491cebece4a546843b3353aed86efeed40fd400bae2ed5cc5f0556d4df1ae6261fc5de362f084734962de941a5cbd18fd97615179441d078a34e5e1d67e6c0d33269b15410559087b5f9fb1734c57a16036e760cbcbbfe0f38854d79af68d6cc8e11c5ec7e47897da51a5f535d9efbd9d30a3302a6fdbcbf6d7a53ef211640af682a5b5ae624f41fd008a47a5f275fcea2d117262aa66480a736a8c776bd2db5032cf7b63dbc5fc60cedb1ab934f088deff96d175a4112879c35bf731034719934aa88ac5653f8a106eca3b66ec042992fb57d0ce737c93860ba3f99607f8aa5825bdb7891cd63554cb74e19489f9b71b38b14121f262b7ba2771805e9ac46fdc39e94e9d68fa9dc51cd0b1b0e2f7baa48331d2831ca78172cd01d9a3540fcb32d473963c96a91c49c9e68417bf328c308a36b15b5cdc7738f9662046ef4c80042af1a4b29ba6ec3a8f8a504381821cd2f32ed2f4fe8e49116e605a39a5b44806fa416fc284b0010ac421046c4a9063285885707f1764dbd6e6865e427b94a6af3442806aacf5f0e5a90b9431f615f3dace77002970a5a3be6d6278283e42b8687e1e030226087f749c8901c3e1731ba08e30601a4565407d53ae6e4b70ba68cd92ebecf54c3c29f0eff277648461e9cd997300c5c9c8b52ae1223ccafa73d7479553f82eb8cead1f094dc9750ff1b3679b6fa2a1a76640978c7b0b8bda3656019ef5bc5f367075c78756c28d96b1bd424a643231866b4be8dfc2cd524edd161bc53130567588c8759c4c16c9c9af20b6c53ba2df5b7971e61a2638e4b1e2ea92ad621ad2f2f0e8d35901205921cc21f74f53d18db556bfcbf990c243991d27d0d47fc019250256aafaec903af4e37603db35722d6f16f43339491cee9a3c9ce8472265000f18fe83ff9ef20a2b6b06c7fc5d9192307a56a899805a0c07510952dde235e16c22d2e3af4369a396709649064f90fa473ddace03be6e6fc1b6d08db8f01058b00f0a03bd737866affec62617751, 
0x294c92ad18c17f3907bced5e71ec39ba0a5a0f3e8d9769e343ba5a8094220753, 
0x98dd76b826575fbb239b926fd1d49aabebd56bf7810a7d886788d7740c28f9de
```

## Verification
The proof can be verified using the solidity contract. It can be checked with:
```
cd ../solidty/
forge test -vvv
```


