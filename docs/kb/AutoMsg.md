---
order: 2
parent:
  title: AutoMessage
  order: 2
---

# AutoMessage

- A predefined end-time is provided when storing the code
- An AutoMessage is encrypted at instantiation by the creator. 
- This can only be sent to the same contract as instantiated by the creator
- Then a *callback signature* is created. This is a signature with a hash containing the address of the contract and the message, so that only the chain is able to execute the message. 
- When the contract is defined to execute, the AutoMessage is retrieved along with the callback signature

## Security
As the AutoMessage is encrypted with the newly created contract address it can not be use my malicious nodes to expose information of other contracts.
