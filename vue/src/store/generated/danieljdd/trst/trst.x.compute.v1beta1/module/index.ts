// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
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

const registry = new Registry(<any>types);

const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
}

interface SignAndBroadcastOptions {
  fee: StdFee,
  memo?: string
}

const txClient = async (wallet: OfflineSigner, { addr: addr }: TxClientOptions = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw MissingWalletError;

  const client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
  const { address } = (await wallet.getAccounts())[0];

  return {
    signAndBroadcast: (msgs: EncodeObject[], { fee, memo }: SignAndBroadcastOptions = {fee: defaultFee, memo: ""}) => client.signAndBroadcast(address, msgs, fee,memo),
    msgInstantiateContract: (data: MsgInstantiateContract): EncodeObject => ({ typeUrl: "/trst.x.compute.v1beta1.MsgInstantiateContract", value: data }),
    msgStoreCode: (data: MsgStoreCode): EncodeObject => ({ typeUrl: "/trst.x.compute.v1beta1.MsgStoreCode", value: data }),
    msgExecuteContract: (data: MsgExecuteContract): EncodeObject => ({ typeUrl: "/trst.x.compute.v1beta1.MsgExecuteContract", value: data }),
    
  };
};

interface QueryClientOptions {
  addr: string
}

const queryClient = async ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr });
};

export {
  txClient,
  queryClient,
};
