---
order: 3
title: Integration
description: How to integrate automation into your dApp front-end
---

## TrustlessJS

We've built a JS framework called [TrustlessJS](https://npmjs.com/package/trustlessjs) to send autotx transactions. An implementation for this is [Triggerportal](https://triggerportal.netlify.app). It contains a message registry that you can use to encode and decode protobuf messages that Trustless Hub supports, including CosmWasm and Osmosis messages.

## MsgSubmitAutoTx

Submitting a MsgSubmitAutoTx takes the following input:

| Field Name        | Data Type                      | Description                                                                                       |
| ----------------- | ------------------------------ | ------------------------------------------------------------------------------------------------- |
| `owner`           | `string`                       | The owner of the transaction                                                                      |
| `connection_id`   | `string`                       | The ID of the connection to use for the transaction (in YAML format)                               |
| `label`           | `string`                       | A label for the transaction                                                                       |
| `msgs`            | `repeated google.protobuf.Any` | A list of arbitrary messages to include in the transaction                                        |
| `duration`        | `string`                       | The amount of time that the transaction code should run for                                       |
| `start_at`        | `uint64`                       | A Unix timestamp representing the custom start time for execution (if set after block inclusion) |
| `interval`        | `string`                       | The interval between automatic message calls                                                     |
| `fee_funds`       | `repeated cosmos.base.v1beta1.Coin` | Optional funds to be used for transaction fees, limiting the amount of fees incurred |
| `depends_on_tx_ids` | `repeated uint64`           | Optional array of transaction IDs that must be executed before the current transaction is allowed to execute |

An example of submitting an MsgSubmitAutoTx in typescript. 

```js

import {
  Coin,
  msgRegistry, Registry,
  toUtf8,
  TrustlessChainClient,
} from 'trustlessjs'


type ExecuteSubmitAutoTxArgs = {
  owner: string
  autoTxData: AutoTxData
  client: TrustlessChainClient
}

export const executeSubmitAutoTx = async ({
  client,
  autoTxData,
  owner,
}: ExecuteSubmitAutoTxArgs): Promise<any> => {

let msgs = []

  for (let msg of autoTxData.msgs) {
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
      connectionId: autoTxData.connectionId, 
      owner,
      msgs,
      label: autoTxData.label ? autoTxData.label : "",
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

The following is used in Triggerportal to estimate fees.

```js

export const getExpectedAutoTxFee = async (client: TrustlessChainClient, durationSeconds: number, lenMsgs: number, intervalSeconds?: number) => {
    try {
        const params = await getAutoTxParams(client) 
        const recurrences = intervalSeconds && intervalSeconds < durationSeconds ? Math.floor(durationSeconds / intervalSeconds) : 1;
        const periodSeconds = intervalSeconds && intervalSeconds < durationSeconds ? intervalSeconds : durationSeconds;
        const periodMinutes = Math.trunc(periodSeconds / 60)
        const flexFeeForPeriod = (Number(params.AutoTxFlexFeeMul) / 100) * periodMinutes
        const autoTxFee = recurrences * flexFeeForPeriod + recurrences * Number(params.AutoTxConstantFee) * lenMsgs
        const autoTxFeeDenom = convertMicroDenomToDenom(autoTxFee, 6)

        return autoTxFeeDenom
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
