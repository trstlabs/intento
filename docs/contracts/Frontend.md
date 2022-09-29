---
order: 5
title: Frontend
description: Frontend with TrustlessJS
---

# Frontend with TrustlessJS

With TrustlessJS you can easily integrate support for Trustless Contracts to your SuperApp.
It is a Typescript wrapper around all of the major interactions with Trustless Contracts and other chain operations.



## NPM

The NPM package is available [here](https://www.npmjs.com/package/trustlessjs) and you can check it out on [Github](https://github.com/danieljdd/TrustlessJS) 

```bash
    npm i trustlessjs
```

## Getting started


``` javascript
 import { TrustlessChainClient } from 'trustlessjs';

 // To create a readonly trustlesscontracts.js client, pass in a gRPC endpoint. 
//To expose the endpoint, the node should NGINX with the right headers. Preflight should include x-grpc-web and content-type, as well as the allow-orgin from header to be [*] or point to the app url.
const trustlessjs = await TrustlessChainClient.create({
  grpcWebUrl: "{grpcendpoint}.{node}.com",
});

const {
  balance: { amount },
} = await trustlessjs.query.bank.balance({
  address: "trust1ap26qrlp8mcq2pg6r47w43l0y8zkqm8a450s03",
  denom: "utrst",
});
```


## Broadcasting Transactions

```ts
import { Wallet, TrustlessChainClient , MsgSend, MsgMultiSend } from "trustlessjs";

const wallet = new Wallet(
  "grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar",
);
const myAddress = wallet.address;

// To create a signer trustlesscontracts.js client, also pass in a wallet
const trustlessjs = await TrustlessChainClient.create({
  grpcWebUrl: "https://grpc.someTrustlessHubnode.com/",
  wallet: wallet,
  walletAddress: myAddress,
  chainId: "pulsar-2",
});

const bob = "trust1dgqnta7fwjj6x9kusyz7n8vpl73l7wsm0gaamk";
const msg = new MsgSend({
  fromAddress: myAddress,
  toAddress: bob,
  amount: [{ denom: "utrst", amount: "1" }],
});

const tx = await trustlessjs.tx.broadcast([msg], {
  gasLimit: 20_000,
  gasPriceInFeeDenom: 0.25,
  feeDenom: "utrst",
});
```

## Keplr Wallet

The recommended way to integrate Keplr is by using `window.keplr.getOfflineSignerOnlyAmino()`:

```ts
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

while (
  !window.keplr ||
  !window.getEnigmaUtils ||
  !window.getOfflineSignerOnlyAmino
) {
  await sleep(100);
}

const CHAIN_ID = "trst_chain_1";
const CHAIN_GRPC = "https://cli.trustlesshub.com/";

await window.keplr.enable(CHAIN_ID);

const keplrOfflineSigner = window.getOfflineSignerOnlyAmino(CHAIN_ID);
const [{ address: myAddress }] = await keplrOfflineSigner.getAccounts();

const trustlessjs = await TrustlessChainClient.create({
  grpcWebUrl: CHAIN_GRPC,
  chainId: CHAIN_ID,
  wallet: keplrOfflineSigner,
  walletAddress: myAddress,
  encryptionUtils: window.getEnigmaUtils(CHAIN_ID),
});

// Note: Using `window.getEnigmaUtils` is optional, it will allow
// Keplr to use the same encryption seed across sessions for the account.
// The benefit of this is that `trustlessjs.query.getTx()` will be able to decrypt
// the response across sessions.
```

### Executing Contracts

Below are examples of executing contracts. 

``` javascript
const tx = await this.trustlessjs.tx.compute.executeContract({
        sender: this.props.user.address,
        contract: tokenAddress,
        msg: {
          send: {
            recipient: recipientAddr,
            amount: decAmount
          },
        },
        funds: [],
      },

        {
          gasLimit: GAS_FOR_APPROVE
        }
      );
      console.log(tx)
      if (tx.code != 0) {
        this.props.notify('error', `Error sending tokens for token address ${tokenAddress}: ${tx.rawLog}`);

      }

trustlessjs.tx.compute.executeContract({
        sender: this.props.user.address,
        contract: tokenAddress,
        msg: {
          increase_allowance: {
            spender: pair.contract_addr,
            amount: UINT128_MAX,
          },
        },

        funds: [],
      },
      {
        gasLimit: GAS_FOR_APPROVE
}
        );
```



``` javascript
const tx = await this.trustlessjs.tx.compute.executeContract({
        sender: this.props.user.address,
        contract: tokenAddress,
        msg: {
          send: {
            recipient: recipientAddr,
            amount: decAmount
          },
        },
        funds: [],
      },

        {
          gasLimit: GAS_FOR_APPROVE
        }
      );
      console.log(tx)
      if (tx.code != 0) {
        this.props.notify('error', `Error sending tokens for token address ${tokenAddress}: ${tx.rawLog}`);

      }

trustlessjs.tx.compute.executeContract({
        sender: this.props.user.address,
        contract: tokenAddress,
        msg: {
          increase_allowance: {
            spender: pair.contract_addr,
            amount: UINT128_MAX,
          },
        },

        funds: [],
      },
      {
        gasLimit: GAS_FOR_APPROVE
}
        );
```


More advanced execution can be integrated too. Here we instantiate a DCA strateg by sending a message to a TIP20 token. With 1-click authorization is given to the newly instantates contract through a reply to the TIP20 contract after instantiation. It is a complex operation on the backend, but simple on the frontend and to the end user.

``` javascript
result = await this.props.trustlessjs.tx.compute.executeContract({
                       sender: this.props.user.address,
                       contract: fromToken,
                       msg: {
                         instantiate_with_allowance:
                          {
                          duration: this.state.duration + this.state.time,
                            code_id: globalThis.config.DCA_CODE_ID,
                           interval: this.state.interval + this.state.timeI,
                          contract_id: 'BackswapDCA RandomID: '+random.toString(),
                         max_allowance: canonicalizeBalance( new BigNumber(this.state.maxAllowance), fromDecimals ).toNumber().toString(),
                          code_hash: globalThis.config.DCA_CODE_HASH,
                          auto_msg: btoa(JSON.stringify({ auto_msg: {}})),
                          msg: dcaData,

                          }

                       },

                       funds: [{"amount": decFee, "denom": "utrst"}],
                     },
                     {
                       gasLimit: GAS_FOR_AUTOSWAP
                     }
                    );

```


### Querying

``` javascript

trustlessjs.query.compute.queryContractPrivateState({address: token, query: {
      transfer_history: {
        key: viewingKey,
        address: this.props.user.address,
        page: 0,
        page_size: 1000,
      },
    }});
```

TrustlessJS will try to decrypt the encrypted outputs locallt given that enigmaUtils is defined. 

