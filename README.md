# ConsensusLayerWithdrawalProtection

ConsensusLayerWithdrawalProtection

Ethereum Consensus Layer Withdrawal Protection provides an additional verification to validators

# Consensus Layer Withdrawal Protection
First proposed by [Jim McDonald](https://ethresear.ch/u/jgm/summary), with contributions from [Pietje Puk](https://ethresear.ch/u/pietjepuk/summary), [Benjamin Chodroff](https://ethresear.ch/u/benjaminchodroff/summary) 

Thanks to the [Ethereum Research](https://ethresear.ch/) team, [ETHStaker](https://ethstaker.cc/) solo staking community, and the hackers who encourage us to make Ethereum even better

# Background
The Consensus Layer change withdrawal credentials operation is [not yet fully specified](https://github.com/ethereum/consensus-specs/pull/2759), but is likely to have at least the following fields:
* Validator index
* Current withdrawal public key
* Proposed execution layer withdrawal address
* Signature by withdrawal private key over the prior fields

The consensus layer change withdrawal credentials proposal is secure for a single user who has certainty their keys and mnemonic have not been compromised. However, as withdrawals on the consensus layer are not yet possible, no user can have absolute certainty that their keys are not compromised until the change withdrawal address is on chain, and too late to change. 

In a situation where multiple users have access to the consensus layer withdrawal private key, it is impossible for the network to differentiate the “legitimate” holder of the key from the “illegitimate”. However, there are signals that could be considered in a wider sense. The most verifiable one is that the legitimate holder was also in control of the execution layer deposit address, in which case the suggestion is to require the proposed withdrawal address to be verified by the execution layer deposit address. Setting a withdrawal address to an execution layer address was not supported by the eth2.0-deposit-cli until v1.1.1 on March 23, 2021, leaving some early adopters wishing they could force change their execution layer address earlier. Forcing this change is not something that can be enforced in-protocol, partly due to lack of information on the beacon chain about the execution layer deposit address and partly due to the fact that this was never listed as a requirement. It is also possible that execution layer deposit addresses is no longer under the control of the legitimate holder of the withdrawal private key. 

However, it is possible for individual nodes to locally restrict the changes they wish to include in blocks they propose, and which they propagate around the network, in way that does not change consensus. It is also possible for client nodes to help broadcast signed change withdrawal address requests to ensure as many nodes witness this change as soon as possible in a fair manner. Further, such change withdrawal signatures can be preloaded in advance to further help nodes filter attacking requests. 

This proposal is to provide an optional mechanism where nodes may load a list of accepted withdrawal addresses and withdrawal credentials to provide additional protection. It aims to request nodes set a priority on withdrawal credential claims that favor a verifiable execution layer deposit address in the event of two conflicting change withdrawal credentials. It also establishes a list of change withdrawal address signed messages to help broadcast "as soon as possible" when the network supports it, and encourage client teams to help use these lists to honor filter and prioritize accepting requests by REST and transmitting them via P2P. This will not change consensus, but may help prevent propagating an attack where a withdrawal key has been knowingly or unknowingly been compromised. 

# Proposal

## Change Withdrawal Address Grace Period
Beacon node clients optionally implement a "change withdrawal address rebroadcast delay" that creates an optional delay in rebroadcasting change withdrawal addresses (suggested to default to 2000 seconds (>5 epochs), or allowing setting to 0 for instructing client to only rebroadcast requests matching the below Change Withdrawal Address Broadcast list) for change withdrawal address requests to allow peer replication of client accepted valid requests, and collecting multiple valid change withdrawal address requests per withdrawal credentials before finalizing the change withdrawal address request. This can prevent a "first to arrive" critical race condition for conflicting change withdraw address. 

## Change Withdrawal Address Broadcast
Support a "Withdrawal Change Address Broadcast" file of lines specifying "validator index,current withdrawal public key,proposed execution layer withdrawal address,consensus layer signature" which will instruct nodes to automatically submit a one time change withdrawal address broadcast message for each valid line. 

This file will give all node operators an opportunity to ensure their valid change withdrawal address messages are broadcast and heard by nodes during the first epoch supporting the change withdrawal address operation. It will also instruct nodes to perpetually prefer accepting and repeating signatures matching the signature in the file, and reject rebroadcasting messages which do not match the signature. 

## Change Withdrawal Address Acceptance
Support a "Withdrawal Address Acceptance" file in the format "withdrawal credentials, execution layer address" which allows clients to load, or packaged by default, a verifiable list matching the consensus layer withdrawal credentials and the original execution layer deposit address. 

While any withdrawal address can be supported, it can be used to help enforce a deposit address is given preference in rebroadcasting even if other clients do not support or have loaded a Change Withdrawal Address Broadcast file. 

## Change Withdrawal Address Handling
First rely on preloaded "Change Withdrawal Address Broadcast" file of verifiable signatures, then fallback to a "Change Withdrawal Address Acceptance" file intended to be loaded with all validator original deposit address information, and then fallback to only "first request" is accepted but delay in rebroadcasting it via P2P. All of this proposal is optional, but we encourage all client teams to include a copy of the uncontested verification file and enable it by default to protect the community. 

If a node is presented with a change withdrawal address operation via the REST API or P2P:

A) Withdrawal credential found in "Change Withdrawal Address Broadcast" file:
  1. Signature Match: If a valid change withdrawal request signature is received for a withdrawal credential that matches the first signature found in the "Change Withdrawal Address Broadcast" file, accept it via REST API, rebroadcast it via P2P, and drop any pending “first preferred” if existing. 
  2. Signature Mismatch: If a valid change withdrawal request is received for a withdrawal credential that does not match the first signature found in the "Change Withdrawal Address Broadcast" file, reject it via REST API, and drop it to prevent rebroadcasting it via P2P.

B) Withdrawal credential not found in or no "Change Withdrawal Address Broadcast" file:
1. Matching withdraw credential and withdraw address in "Change Withdrawal Address Acceptance" file: If a valid change withdrawal address request is received for a withdrawal credential that matches the first found withdrawal address provided in the "Change Withdrawal Address Acceptance" file, accept it via REST API, rebroadcast it via P2P, and drop any pending “first preferred” if existing. 
2. Mismatching withdraw credential and withdraw address in "Change Withdrawal Address Acceptance" file: If a valid change withdrawal request is received for a withdrawal credential that does not match the first found withdrawal address provided in the "withdrawal address" file, reject it via REST API, and drop it to prevent rebroadcasting it via P2P.
3. Missing withdraw address in or no "Change Withdrawal Address Acceptance" file: 

    i. First Preferred: If first valid change withdrawal request is received for a not finalized withdrawal credential that does not not have any listed withdrawal credential entry in the "Change Withdrawal Address Acceptance" file, accept it via REST API, but do not yet rebroadcast it via P2P (“grace period”). Once the client “Change Withdrawal Address Grace Period” has expired and no other messages have invalidated this message, rebroadcast the request via P2P. 
  
    ii. Subsequent Rejected: If an existing valid "First Preferred" request exists for a not finalized withdrawal credential, reject it via REST API, and drop it to prevent rebroadcasting via P2P. 


Note that these restrictions do not apply to withdrawal credential change operations in blocks. If an operation has been included on-chain it is by definition valid regardless of its contents.

It is critical to understand that this proposal is not a consensus change. Nothing in this proposal restricts the validity of withdrawal credential change operations within the protocol. It is a voluntary change by client teams to build this functionality in to their beacon nodes, and a voluntary change by node operators to accept any or all of the restrictions and broadcasting capabilities suggested by end users.

Because of the above, even if fully implemented, it will be down to chance as to which validators propose blocks, and which voluntary restrictions those validators’ beacon nodes are running. Node operators can do what they will to help the community prevent attacks on any compromised consensus layer keys, but there are no guarantees of success this will prevent a successful attack. 

# Documentation 
## Change Withdrawal Address Acceptance File
A file intended to be preloaded with all consensus layer withdrawal credentials and verifiable execution layer deposit address. This file will be generated by a script and able to be independently verified by community members using the consensus and execution layers, and intended to be included by all clients, enabled by default. Client nodes are encouraged to enable packaging this independently verifiable list with the client software, and enable it by default to help further protect the community from unsuspected attacks. 

depositAddress.txt format:
```withdrawal credential, withdrawal address```

Example depositAddress.txt:
```
000092c20062cee70389f1cb4fa566a2be5e2319ff43965db26dbaa3ce90b9df99,01c34eb7e3f34e54646d7cd140bb7c20a466b3e852
0000d66cf353931500a54cbd0bc59cbaac6690cb0932f42dc8afeddc88feeaad6f,01c34eb7e3f34e54646d7cd140bb7c20a466b3e852
0000d6b91fbbce0146739afb0f541d6c21e8c41e92b97874828f402597bf530ce4,01c34eb7e3f34e54646d7cd140bb7c20a466b3e852
000037ca9a1cf2223d8b9f81a14d4937fef94890ae4fcdfbba928a4dc2ff7fcf3b,01c34eb7e3f34e54646d7cd140bb7c20a466b3e852
0000344b6c73f71b11c56aba0d01b7d8ad83559f209d0a4101a515f6ad54c89771,01f19b1c91faacf8071bd4bb5ab99db0193809068e
```

## Change Withdrawal Address Broadcast Claim
A community collected and independently verifiable list of "Change Withdrawal Address Broadcasts" containing verifiable claims will be collected. Client teams and node operators may verify these claims independently and are suggestted to include "Uncontested and Verified" claims enabled by default in their package. 

To make a verifiable claim, users must upload using their GitHub ID with the following contents to the [CLWP repository](https://github.com/benjaminchodroff/ConsensusLayerWithdrawalProtection) in a text file "claims/validatorIndex-gitHubUser.txt" such as "123456-benjaminchodroff.txt"

123456-benjaminchodroff.txt:
```
current_withdrawal_public_key=b03c5ea17b017cffd22b6031575c4453f20a4737393de16a626fb0a8b0655fe66472765720abed97e8022680204d3868
proposed_withdrawal_address=0108f2e9Ce74d5e787428d261E01b437dC579a5164
consensus_layer_withdrawal_signature=
execution_layer_deposit_signature=
execution_layer_withdrawal_signature=
email=noreply@ethereum.org
```

| key | value | 
| ----| ------|
| current_withdrawal_public_key | The "pubkey" field found in deposit_data json file matching the validator index|
| proposed_withdrawal_address | The address in ethereum you wish to authorize withdrawals to, prefaced by "01" without any "0x", such that an address "0x08f2e9Ce74d5e787428d261E01b437dC579a5164" turns into "0108f2e9Ce74d5e787428d261E01b437dC579a5164 |
| consensus_layer_withdrawal_signature | The verifiable signature generated by signing "validator_index,current_withdrawal_public_key,proposed_withdrawal_address" using the consensus layer private key |
| consensus_layer_withdrawal_credentials | The verifiable original "withdrawal_credentials" found in deposit_data json file, should start with "00" or "01".
| execution_layer_deposit_signature | The verifiable signature generated by signing "validator_index,current_withdrawal_public_key,proposed_withdrawal_address" using the execution layer deposit address private key | 
| execution_layer_withdrawal_signature | The verifiable signature generated by signing "validator_index,current_withdrawal_public_key,proposed_withdrawal_address" using the execution layer proposed withdrawal address private key. This may be the same result as the "execution_layer_deposit_signature" if the user intends to withdraw to the same execution layer deposit address. | 
| email | Any actively monitored email address to notify in the event of contention (required, but may specify noreply@ethereum.org to opt-out of all contentions and notifications - not recommended) |

## Claim Acceptance
In order for a submission to be merged into CLWP GitHub repository, the submission must have:
1. Valid filename in the format validatorIndex-githubUsername.txt
2. Valid validator index which is deposited, pending, or active on the consensus layer 
3. Matching Github username in file name to the user submitting the request
4. Verifiable consensus_layer_withdrawal_signature, execution_layer_deposit_signature, and execution_layer_withdrawal_signature
5. All required fields in the file with no other content present

All merge requests that fail will be provided a reason from above which must be addressed prior to merge. Any future verifiable amendments to accepted claims must be proposed by the same GitHub user, or it will be treated as a contention.

## Change Withdrawal Address Broadcast
Anyone in the community will be able to generate the following verifiable files from the claims provided:
	A. UncontestedVerified - Community collected list of all verifiable uncontested change withdrawal address final requests (no conflicting withdrawal credentials allowed from different GitHub users)
	B. ContestedVerified - Community collected list of all contested verifiable change withdrawal address requests (will contain only verifiable but conflicting withdrawal credentials from different GitHub users)

If a contested but verified "Change Withdrawal Address Broadcast" request arrives to the GitHub community, all parties will be notified via GitHub, forced into the ContestedVerified list, and may try to convince the wider community be providing any off chain evidence supporting their claim to then include their claim in nodes.All node operators are encouraged to load the UncontestedVerified signatures file as enabled, and optionally append only ContestedVerified signatures that they have been convinced are the rightful owner in a manner to further strengthen the community. 
