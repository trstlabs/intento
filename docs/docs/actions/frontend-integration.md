---
sidebar_position: 5
title: Frontend Integration
description: How to integrate automation into your interchain dApp
---

## TrustlessJS

We've built a JS framework called [TrustlessJS](https://npmjs.com/package/trustlessjs) to submit actions to the chain. It contains a message registry that you can use to encode and decode protobuf messages that Intento supports, including CosmWasm and Osmosis messages. An implementation is [TriggerPortal](https://triggerportal.zone).

An example of submitting an MsgSubmitAction in typescript. A label is optional but recommended to keep track an overview of the actions.

```js

import {
  Coin,
  msgRegistry, Registry,
  toUtf8,
  TrustlessChainClient,
} from 'trustlessjs'


type ExecuteSubmitActionArgs = {
  owner: string
  ActionData: ActionData
  client: TrustlessChainClient
}

export const executeSubmitAction = async ({
  client,
  ActionData,
  owner,
}: ExecuteSubmitActionArgs): Promise<any> => {

let msgs = []

  for (let msg of ActionData.msgs) {
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

    await client.tx.action.submit_action({
      connectionId: ActionData.connectionId,
      owner,
      msgs,
      label,
      duration,
      interval,
      startAt,
      feeFunds,
    },
      { gasLimit: 100_000 }
    )

}
```

## Example ActionFee calculation

The following is used in TriggerPortal to estimate fees.

```js
export const getExpectedActionFee = async (
  client: TrustlessChainClient,
  durationSeconds: number,
  lenMsgs: number,
  intervalSeconds?: number
) => {
  try {
    const params = await getActionParams(client);
    const recurrences =
      intervalSeconds && intervalSeconds < durationSeconds
        ? Math.floor(durationSeconds / intervalSeconds)
        : 1;
    const periodSeconds =
      intervalSeconds && intervalSeconds < durationSeconds
        ? intervalSeconds
        : durationSeconds;
    const periodMinutes = Math.trunc(periodSeconds / 60);
    const flexFeeForPeriod =
      (Number(params.ActionFlexFeeMul) / 100) * periodMinutes;
    const ActionFee =
      recurrences * flexFeeForPeriod +
      recurrences * Number(params.ActionConstantFee) * lenMsgs;
    const ActionFeeDenom = convertMicroDenomToDenom(ActionFee, 6);

    return ActionFeeDenom;
  } catch (e) {
    console.error("err(getExpectedActionFee):", e);
  }
};

async function getActionParams(client: TrustlessChainClient) {
  console.log("getActionParams");
  try {
    const resp = await client.query.action.params({});
    return resp.params;
  } catch (e) {
    console.error("err(getActionParams):", e);
  }
}
```

The function returns a Promise that resolves to the expected transaction fee in Intento chain's native denomination, INTO.

The JavaScript function getExpectedActionFee calculates the expected transaction fee for a trustless chain transaction based on the duration of the transaction, the length of the messages to be sent, and the recurrence interval (optional). The formula for calculating the fee is:

ActionFee = recurrences _ flexFeeForPeriod + recurrences _ constantFee \* lenMsgs

where:

recurrences is the number of times the transaction will recur during the specified duration. It is calculated as:

recurrences = intervalSeconds && intervalSeconds < durationSeconds ? Math.floor(durationSeconds / intervalSeconds) : 1

flexFeeForPeriod is the flex fee for each recurrence, calculated as:

flexFeeForPeriod = (Number(params.ActionFlexFeeMul) / 100) \* periodMinutes

where params.ActionFlexFeeMul is a parameter retrieved from the Intento client and periodMinutes is the duration of each recurrence in minutes.

constantFee is a constant fee for each message sent in the transaction. It is also retrieved from Intento client.

lenMsgs is the length of the messages to be sent in the transaction.
