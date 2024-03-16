---
order: 7
title: Frontend Integration
description: How to integrate automation into your interchain dApp
---

## TrustlessJS

We've built a JS framework called [TrustlessJS](https://npmjs.com/package/trustlessjs) to send AutoTx transactions. An implementation for this is [TriggerPørtal](https://triggerportal.zone). It contains a message registry that you can use to encode and decode protobuf messages that Intento supports, including CosmWasm and Osmosis messages.

An example of submitting an MsgSubmitAutoTx in typescript. A label is optional but recommended to keep track an overview of triggers.
Sta

```js

import {
  Coin,
  msgRegistry, Registry,
  toUtf8,
  TrustlessChainClient,
} from 'trustlessjs'


type ExecuteSubmitAutoTxArgs = {
  owner: string
  AutoTxData: AutoTxData
  client: TrustlessChainClient
}

export const executeSubmitAutoTx = async ({
  client,
  AutoTxData,
  owner,
}: ExecuteSubmitAutoTxArgs): Promise<any> => {

let msgs = []

  for (let msg of AutoTxData.msgs) {
    const masterRegistry = new Registry(msgRegistry);

    let value = JSON.parse(msg)["value"]
    let typeUrl: string = JSON.parse(msg)["typeUrl"].toString()

    const encodeObject = {
      typeUrl,
      value
    }

    let msgAny = masterRegistry.encodeAsAny(encodeObject)
    msgs.push(msgAny)
  }

    await client.tx.auto_tx.submit_auto_tx({
      connectionId: AutoTxData.connectionId, 
      owner,
      msgs,
      label: AutoTxData.label ? AutoTxData.label : "",
      duration,
      interval,
      startAt,
      feeFunds,
    },
      { gasLimit: 100_000 }
    )
    
}
```

## Example AutoTxFee calculation

The following is used in TriggerPørtal to estimate fees.

```js

export const getExpectedAutoTxFee = async (client: TrustlessChainClient, durationSeconds: number, lenMsgs: number, intervalSeconds?: number) => {
    try {
        const params = await getAutoTxParams(client) 
        const recurrences = intervalSeconds && intervalSeconds < durationSeconds ? Math.floor(durationSeconds / intervalSeconds) : 1;
        const periodSeconds = intervalSeconds && intervalSeconds < durationSeconds ? intervalSeconds : durationSeconds;
        const periodMinutes = Math.trunc(periodSeconds / 60)
        const flexFeeForPeriod = (Number(params.AutoTxFlexFeeMul) / 100) * periodMinutes
        const AutoTxFee = recurrences * flexFeeForPeriod + recurrences * Number(params.AutoTxConstantFee) * lenMsgs
        const AutoTxFeeDenom = convertMicroDenomToDenom(AutoTxFee, 6)

        return AutoTxFeeDenom
    } catch (e) { console.error('err(getExpectedAutoTxFee):', e) }
}


async function getAutoTxParams(client: TrustlessChainClient) {
    console.log("getAutoTxParams")
    try {
        const resp = await client.query.auto_tx.params({})
        console.log(resp)
        return resp.params
    } catch (e) { console.error('err(getAutoTxParams):', e) }
}
```
The function returns a Promise that resolves to the expected transaction fee in Intento chain's native denomination, INTO.

The JavaScript function getExpectedAutoTxFee calculates the expected transaction fee for a trustless chain transaction based on the duration of the transaction, the length of the messages to be sent, and the recurrence interval (optional). The formula for calculating the fee is:

AutoTxFee = recurrences * flexFeeForPeriod + recurrences * constantFee * lenMsgs

where:

recurrences is the number of times the transaction will recur during the specified duration. It is calculated as:

recurrences = intervalSeconds && intervalSeconds < durationSeconds ? Math.floor(durationSeconds / intervalSeconds) : 1

flexFeeForPeriod is the flex fee for each recurrence, calculated as:

flexFeeForPeriod = (Number(params.AutoTxFlexFeeMul) / 100) * periodMinutes

where params.AutoTxFlexFeeMul is a parameter retrieved from the Intento client and periodMinutes is the duration of each recurrence in minutes.

constantFee is a constant fee for each message sent in the transaction. It is also retrieved from Intento client.

lenMsgs is the length of the messages to be sent in the transaction.
