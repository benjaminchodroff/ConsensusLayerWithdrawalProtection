
# Consensus Layer Withdrawal Protection

Submit a pull request with your Ethereum validator change withdrawal credential message to have the CLWP community broadcast it as early as possible for security.

## Acknowledgements

 - [EthDo](https://github.com/wealdtech/ethdo)
 - [EthStaker](https://ethstaker.cc/)

## Documentation

Consensus Layer Withdrawal Protection (CLWP) is an optional way to set your Ethereum validator withdrawal address as securely and early as possible. By submitting your signed validator withdrawal address to this repository, we will broadcast this message to many volunteer and large node operators using low latency bots. As each BLS signature is cryptographically generated, these messages may be verified by anyone in the community and are impossible to forge unless the seed phrase is compromised. If the seed phrase is compromised and a contested signature is received by this repository, we will remove all submissions for the given validator and arbitrate the issue using a Kleros curated list. 

CLWP submissions are generated using ethdo or staking-deposit-cli. For full documentation on ethdo usage, please refer to [Change Withdrawal Credentials](https://github.com/wealdtech/ethdo/blob/master/docs/changingwithdrawalcredentials.md). Below are the excerpt steps required to join the CLWP voluntary broadcast protection. For your own protection, please run the steps using your own node and execute on an airgapped machine for signing processes. However, if you do not have access to your own Ethereum beacon node, we have provided an offline-preparation.json file that you may use instead. 

If your validator has not set an execution layer withdrawal address, your withdrawal credentials will start with "0x00". If they have been set, they will start with "0x01". You may easily check your validator status using beaconchain under the "Withdrawal Credentials", such as here:
https://goerli.beaconcha.in/validator/99a29d72501fc49a748d11367b0b2b80be2e5c93cc28a512e06fb40142666e206590ee637ba1bf1e8adfd0e9de3665d5#deposits

At the launch of the Capella/Shanghai hardfork, every validator with Withdrawal Credentials starting with 0x00 will be allowed to perform **a one time** operation to change withdrawal credentials from 0x00 to an execution layer address. You will need validator mnemonic seed phrase (or withdrawal private key) to sign this transaction. For more details on how withdrawals and the set withdrawal address work, we recommend reading https://notes.ethereum.org/@launchpad/withdrawals-faq

## Demo

CLWP Video Demo: https://www.youtube.com/watch?v=EWkGyorgpAg

CLWP Presentation: https://docs.google.com/presentation/d/1qV0NP2-5UZI51Ja7Vf_ANOR9AqIPvNpfuTzxbJVCn90/edit

## Steps

2023-01-21 - We are accepting CLWP submissions using ethdo version 1.27 or later. An offline-preparation.json file has been updated for Capella

### Mainnet

1. Download an ethdo release, an open source Ethereum command line interface for validator actions, onto a clean computer
https://github.com/wealdtech/ethdo/releases/tag/v1.27.1 
2. Run  ethdo to generate a “change-operations.json” file. Choose either the “Easy without Node” or “Offline with Node” approach.

    * Easy without Node - use a cached “prepare offline” beacon node list from GitHub (no beacon node required, but needs secure offline computer):
     
        ```
        # Download https://github.com/benjaminchodroff/ConsensusLayerWithdrawalProtection/blob/main/offline-preparation.json.mainnet.tar.gz 
        tar -zxf offline-preparation.json.mainnet.tar.gz
        cp offline-preparation.json.mainnet offline-preparation.json

        # In the same directory as offline-preparation.json file, run ethdo (Triple check your withdrawal address)
        ./ethdo validator credentials set --offline --mnemonic="abandon … art" --withdrawal-address=0x0123…cdef
        ```

    * Offline with Node - prepare your own offline-preparation.json file using the “--offline” flag (Advanced)
     
        ```
        # On Beacon Node: 
        ./ethdo validator credentials set --prepare-offline
        
        # Copy the offline-preparation.json and ethdo to your airgapped secure machine
        
        # Secure Machine: 
        ./ethdo validator credentials set --offline --mnemonic="abandon … art" --withdrawal-address=0x0123…cdef
        
        # Copy the resulting change-operations.json file back to your online computer
        ```

Don’t forget to wipe your command line history to prevent your seed phrase from being stored in memory:

```
history -c
```

3. Inspect the resulting change-operations.json, and move it to a validatorIndex.json file (such as 123456.json) per validator and submit a pull request to have it included (ask for help in Discord) https://github.com/benjaminchodroff/ConsensusLayerWithdrawalProtection

    * validator_index: did it find all your validators? Make sure all your validators are found, then create one file per validator index. 
    * from_bls_pubkey: This is your public key of your withdrawal address, not your validator public address (Not easy to verify, ignore it)
    * to_execution_address: Triple check this is your intended withdrawal address! Not case sensitive. 
    * Signature: This is the BLS signature of your change withdrawal address operation. We will verify if it matches. Do not attempt to change it or it will be invalid.

If you have multiple validators, you will need a text editor to split the file - review the format of existing submissions as an example.
Never modify the validator_index, from_bls_pubkey, to_execution_address, or signature or it will invalidate the submission. 

4. Done! We will review your submission, merge it, and many CLWP node operators will volunteer to broadcast your set withdrawal address. Submissions must be received by February 28, 2023. You may broadcast your own submissions and we welcome others to help broadcast all the CLWP submissions to their own beacon nodes. If we receive conflicting change operations from multiple parties, we will require all parties to arbitrate the issue on a Kleros curated list before adding the winner back to this GitHub repository:
https://curate.kleros.io/tcr/1/0x479083b5343aB89bb39608e3176D750c8A6957B5

Move the resulting change-operations.json file into the mainnet folder with individual files for each validatorIndex.json, verify the withdrawal address again, and submit a pull request to have it included in CLWP protection. If you need help or prefer not to link your GitHub account to your validator, reach out to an admin on the Support below and we can assist. If you have multiple validators, you will need to split the file to have a single submission manually using a text editor. Never modify the validator index, public key, or withdrawal address, or signature or it will invalidate the submission. If you have many validators, you may install "jq" on linux and split the change-operations.json file using it:

```
for ((i=0; i<`jq -ec '.|length' change-operations.json`;i++)); do validator=`jq -ec ".[${i}].message.validator_index|tonumber" change-operations.json`; echo "`jq -ec "["".[${i}]""]" change-operations.json`" > ${validator}.json;done
``` 
 
Volunteer to run the broadcast of change-operations-clwp-mainnet.json on your node to help protect the community.  

### Goerli

The steps above are tested to work with Goerli as well. Use the offline-preparation.json.goerli.tar.gz file, and place your change-operations.json file in the goerli directory as validatorIndex.json per each validator. 


## Support

For community support, create an issue or join our OffChain Discord channel #ethereum-consensus-layer-withdrawal-protection https://discord.gg/pwuPA6K4zg

## Volunteer

We welcome every node operator to volunteer by loading CLWP submissions into their node in advance of the Capella hard fork on each chain. Please sign up to our mailing list on https://clwp.xyz to receive notice when we have detailed instructions in mid-February. There is no cost (you don't even need to stake), and there is no penalty even if an attacker "wins the race" against a CLWP submission. Your beacon chain client will simply ignore the local submission and use the on chain consensus. All CLWP submissions may be independently verified and, even if a submission in this repository was invalid, your local beacon chain client would refuse to process it without penalty. 
