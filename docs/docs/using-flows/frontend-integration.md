---
sidebar_position: 6
title: Frontend Integration
description: How to integrate intent-based flows into your interchain dApp
---

## IntentoJS

We've built a JS framework called [IntentoJS](https://npmjs.com/package/intentojs) to submit flows to the chain. It contains a message registry that you can use to encode and decode protobuf messages that Intento supports, including CosmWasm and Osmosis messages. An implementation is [TriggerPortal](https://triggerportal.zone).

An example of submitting an MsgSubmitFlow in typescript. A label is optional but recommended to keep track an overview of the flows.

```js
import { Coin, msgRegistry, Registry } from "intentojs";

/**
 * Encodes transaction messages using the Intento registry.
 */
const encodedMsgs: any[] = [];

for (const msg of FlowData.msgs) {
  const masterRegistry = new Registry(msgRegistry);

  const parsedMsg = JSON.parse(msg);
  const typeUrl: string = parsedMsg["typeUrl"].toString();
  const value = parsedMsg["value"];

  const encodeObject = { typeUrl, value };

  // Encode message into Any format
  const msgAny = masterRegistry.encodeAsAny(encodeObject);
  encodedMsgs.push(msgAny);
}

/**
 * Constructs the submitFlow message for Intento.
 */
const msgSubmitFlow =
  intento.intent.v1beta1.MessageComposer.withTypeUrl.submitFlow({
    label: "My Flow", // Optional flow label
    connectionId: "connection-123",
    owner: "into1wdplq6qjh2xruc7qqagma9ya665q6qhcpse4k6",
    msgs: encodedMsgs,
    duration: "1440h", // Flow duration (24h * 60d)
    interval: "600s", // Execution interval (10 min)
    startAt: "1739781618", // UNIX timestamp for start time
    feeFunds: [{ denom: "uinto", amount: "5000000" }], // Funding for fees
  });

/**
 * Signs and broadcasts the transaction.
 */
client.signAndBroadcast(owner, [msgSubmitFlow], {
  amount: [],
  gas: "300000",
});
```

## Example FlowFee calculation

The following is used in TriggerPortal to estimate fees.

```js
/**
 * Calculates the expected flow fee based on gas usage, number of messages, and duration.
 *
 * @param {Params} intentParams - The parameters for the intent flow module.
 * @param {number} gasUsed - The (expected) total gas used for a flw entry, including the action and conditions (typically ranging from 80_000 for a simple action to 500_000 for a larger one).
 * @param {number} lenMsgs - The number of messages in the flow.
 * @param {number} durationSeconds - The total duration of the flow in seconds.
 * @param {number} [intervalSeconds] - Optional interval in seconds between executions.
 * @returns {number} The calculated flow fee in decimal format.
 */
export const getExpectedFlowFee = (
  intentParams: Params,
  gasUsed: number,
  lenMsgs: number,
  durationSeconds: number,
  intervalSeconds?: number
): number => {
  // Determine the number of times the flow will recur
  const recurrences =
    intervalSeconds && intervalSeconds < durationSeconds
      ? Math.floor(durationSeconds / intervalSeconds)
      : 1;

  // Calculate the flex fee based on gas usage
  const flexFeeForPeriod =
    (Number(intentParams.flowFlexFeeMul) / 100) * gasUsed;

  // Compute the total flow fee
  const flowFee =
    recurrences * flexFeeForPeriod +
    recurrences * Number(intentParams.burnFeePerMsg) * lenMsgs;

  // Convert fee from micro-denomination to standard denomination
  const flowFeeDenom = convertMicroDenomToDenom(flowFee, 6);

  return Number(flowFeeDenom.toFixed(4));
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
