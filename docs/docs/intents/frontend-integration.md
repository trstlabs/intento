---
sidebar_position: 6
title: Frontend Integration
description: How to integrate intent-based flows into your interchain dApp
---

## IntentoJS

We've built a JS framework called [IntentoJS](https://npmjs.com/package/intentojs) to submit flows to the chain. It contains a message registry that you can use to encode and decode protobuf messages that Intento supports, including CosmWasm and Osmosis messages. An implementation is [TriggerPortal](https://triggerportal.zone).

An example of submitting an MsgSubmitFlow in typescript. A label is optional but recommended to keep track an overview of the flows.

```js

import {
  Coin,
  msgRegistry, Registry,
  toUtf8,
  IntentoChainClient,
} from 'intentojs'


type ExecuteSubmitFlowArgs = {
  owner: string
  FlowData: FlowData
  client: IntentoChainClient
}

export const executeSubmitFlow = async ({
  client,
  FlowData,
  owner,
}: ExecuteSubmitFlowArgs): Promise<any> => {

let msgs = []

  for (let msg of FlowData.msgs) {
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

    await client.tx.intent.submit_flow({
      connectionId: FlowData.connectionId,
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

## Example FlowFee calculation

The following is used in TriggerPortal to estimate fees.

```js
export const getExpectedFlowFee = async (
  client: IntentoChainClient,
  durationSeconds: number,
  lenMsgs: number,
  intervalSeconds?: number
) => {
  try {
    const params = await getFlowParams(client);
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
      (Number(params.FlowFlexFeeMul) / 100) * periodMinutes;
    const FlowFee =
      recurrences * flexFeeForPeriod +
      recurrences * Number(params.FlowConstantFee) * lenMsgs;
    const FlowFeeDenom = convertMicroDenomToDenom(FlowFee, 6);

    return FlowFeeDenom;
  } catch (e) {
    console.error("err(getExpectedFlowFee):", e);
  }
};

async function getFlowParams(client: IntentoChainClient) {
  console.log("getFlowParams");
  try {
    const resp = await client.query.flow.params({});
    return resp.params;
  } catch (e) {
    console.error("err(getFlowParams):", e);
  }
}
```

The function returns a Promise that resolves to the expected transflow fee in Intento chain's native denomination, INTO.

The JavaScript function getExpectedFlowFee calculates the expected transflow fee for a trustless chain transflow based on the duration of the transflow, the length of the messages to be sent, and the recurrence interval (optional). The formula for calculating the fee is:

FlowFee = recurrences _ flexFeeForPeriod + recurrences _ constantFee \* lenMsgs

where:

recurrences is the number of times the transflow will recur during the specified duration. It is calculated as:

recurrences = intervalSeconds && intervalSeconds < durationSeconds ? Math.floor(durationSeconds / intervalSeconds) : 1

flexFeeForPeriod is the flex fee for each recurrence, calculated as:

flexFeeForPeriod = (Number(params.FlowFlexFeeMul) / 100) \* periodMinutes

where params.FlowFlexFeeMul is a parameter retrieved from the Intento client and periodMinutes is the duration of each recurrence in minutes.

constantFee is a constant fee for each message sent in the transflow. It is also retrieved from Intento client.

lenMsgs is the length of the messages to be sent in the transflow.
