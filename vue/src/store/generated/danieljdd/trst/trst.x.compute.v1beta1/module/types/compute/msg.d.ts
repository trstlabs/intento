import { Writer, Reader } from "protobufjs/minimal";
import { Coin } from "../cosmos/base/v1beta1/coin";
export declare const protobufPackage = "trst.x.compute.v1beta1";
export interface MsgStoreCode {
    sender: Uint8Array;
    /** WASMByteCode can be raw or gzip compressed */
    wasmByteCode: Uint8Array;
    /** Source is a valid absolute HTTPS URI to the contract's source code, optional */
    source: string;
    /** Builder is a valid docker image name with tag, optional */
    builder: string;
    /**
     * InstantiatePermission to apply on contract creation, optional
     *  AccessConfig InstantiatePermission = 5;
     */
    contractPeriod: number;
}
export interface MsgInstantiateContract {
    sender: Uint8Array;
    /**
     * Admin is an optional address that can execute migrations
     *  bytes admin = 2 [(gogoproto.casttype) = "github.com/cosmos/cosmos-sdk/types.AccAddress"];
     */
    callbackCodeHash: string;
    codeId: number;
    contractId: string;
    initMsg: Uint8Array;
    lastMsg: Uint8Array;
    initFunds: Coin[];
    callbackSig: Uint8Array;
}
export interface MsgExecuteContract {
    sender: Uint8Array;
    contract: Uint8Array;
    msg: Uint8Array;
    callbackCodeHash: string;
    sentFunds: Coin[];
    callbackSig: Uint8Array;
}
export declare const MsgStoreCode: {
    encode(message: MsgStoreCode, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgStoreCode;
    fromJSON(object: any): MsgStoreCode;
    toJSON(message: MsgStoreCode): unknown;
    fromPartial(object: DeepPartial<MsgStoreCode>): MsgStoreCode;
};
export declare const MsgInstantiateContract: {
    encode(message: MsgInstantiateContract, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgInstantiateContract;
    fromJSON(object: any): MsgInstantiateContract;
    toJSON(message: MsgInstantiateContract): unknown;
    fromPartial(object: DeepPartial<MsgInstantiateContract>): MsgInstantiateContract;
};
export declare const MsgExecuteContract: {
    encode(message: MsgExecuteContract, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgExecuteContract;
    fromJSON(object: any): MsgExecuteContract;
    toJSON(message: MsgExecuteContract): unknown;
    fromPartial(object: DeepPartial<MsgExecuteContract>): MsgExecuteContract;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
