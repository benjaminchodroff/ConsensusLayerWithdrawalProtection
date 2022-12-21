
# Consensus Layer Withdrawal Protection

Submit a pull request with your Ethereum validator change withdrawal credential message to have the CLWP community broadcast it as early as possible for security.

## Acknowledgements

 - [EthDo](https://github.com/wealdtech/ethdo)
 - [EthStaker](https://ethstaker.cc/)

## Documentation

For full documentation on ethdo usage, please refer to [Change Withdrawal Credentials](https://github.com/wealdtech/ethdo/blob/master/docs/changingwithdrawalcredentials.md). Below are the excerpt steps required to join the CLWP voluntary broadcast protection. For your own protection, please run the steps using your own node and execute on an airgapped machine for signing processes. However, for demonstration purposes, I have included an excerpt of the offline-preparation.json files in compressed format.

If your validator has not set an execution layer withdrawal address, your withdrawal credentials will start with "0x00". If they have been set, they will start with "0x01". You may easily check your validator status using beaconchain, such as here:
https://goerli.beaconcha.in/validator/99a29d72501fc49a748d11367b0b2b80be2e5c93cc28a512e06fb40142666e206590ee637ba1bf1e8adfd0e9de3665d5#deposits

At the launch of the Capella/Shanghai hardfork, every validator with Withdrawal Credentials starting with 0x00 will be allowed to perform **a one time** operation to change withdrawal credentials from 0x00 to an execution layer address. You will need validator mnemonic seed phrase (or withdrawal private key) to sign this transaction.

## Steps

### Mainnet
```
# Download this repository
git clone https://github.com/benjaminchodroff/ConsensusLayerWithdrawalProtection
cd ConsensusLayerWithdrawalProtection
# Download ethdo
git clone https://github.com/wealdtech/ethdo.git
cd ethdo
# Build ethdo (ensure you have go >= 1.13)
go build .
# Unpack the offline-preparation.json file (or you may generate your own using --prepare-offline and your own beacon node)
tar -zxf ../offline-preparation.json.mainnet.tar.gz
cp offline-preparation.json.mainnet offline-preparation.json
./ethdo validator credentials set --offline  --json --fork-version 0x03000000 --withdrawal-address 0xAnExecutionLayerAddress --mnemonic "your seed phrase"
history -c 
```
Combine your resulting change-operations.json file into change-operations-clwp-mainnet.json file, and submit a pull request to have it included in CLWP protection. 
Volunteer to run the broadcast of change-operations-clwp-mainnet.json on your node to help protect the community.  

### Goerli
The steps above are tested to work with Goerli as well. Use the change-operations-clwp-goerli.json file. 


## Support

For community support, create an issue or join our OffChain Discord channel #ethereum-consensus-layer-withdrawal-protection https://discord.gg/pwuPA6K4zg


