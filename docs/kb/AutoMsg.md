---
order: 2
parent:
  title: Automated Message
  order: 2
---

# Automated Message

- A predefined end-time is provided when storing the code
- An automated message is encrypted at instantiation by the creator
- This is then sent to the same enviromnemt as the instantiation message
- Then a *callback signature* is created. This is a signature with a hash containing the address of the contract and the message, so that only the TRST blockchain is able to execute the message. 
- When the contract ends, the automated message is retrieved along with the callback signature, and then this message is executed automatically 