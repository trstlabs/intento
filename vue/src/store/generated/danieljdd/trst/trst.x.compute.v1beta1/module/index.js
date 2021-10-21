// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgInstantiateContract } from "./types/compute/msg";
import { MsgStoreCode } from "./types/compute/msg";
import { MsgExecuteContract } from "./types/compute/msg";
const types = [
    ["/trst.x.compute.v1beta1.MsgInstantiateContract", MsgInstantiateContract],
    ["/trst.x.compute.v1beta1.MsgStoreCode", MsgStoreCode],
    ["/trst.x.compute.v1beta1.MsgExecuteContract", MsgExecuteContract],
];
export const MissingWalletError = new Error("wallet is required");
const registry = new Registry(types);
const defaultFee = {
    amount: [],
    gas: "200000",
};
const txClient = async (wallet, { addr: addr } = { addr: "http://localhost:26657" }) => {
    if (!wallet)
        throw MissingWalletError;
    const client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
    const { address } = (await wallet.getAccounts())[0];
    return {
        signAndBroadcast: (msgs, { fee, memo } = { fee: defaultFee, memo: "" }) => client.signAndBroadcast(address, msgs, fee, memo),
        msgInstantiateContract: (data) => ({ typeUrl: "/trst.x.compute.v1beta1.MsgInstantiateContract", value: data }),
        msgStoreCode: (data) => ({ typeUrl: "/trst.x.compute.v1beta1.MsgStoreCode", value: data }),
        msgExecuteContract: (data) => ({ typeUrl: "/trst.x.compute.v1beta1.MsgExecuteContract", value: data }),
    };
};
const queryClient = async ({ addr: addr } = { addr: "http://localhost:1317" }) => {
    return new Api({ baseUrl: addr });
};
export { txClient, queryClient, };
