import { Writer, Reader } from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
export declare const protobufPackage = "trst.x.compute.v1beta1";
export declare enum AccessType {
    UNDEFINED = 0,
    NOBODY = 1,
    ONLY_ADDRESS = 2,
    EVERYBODY = 3,
    UNRECOGNIZED = -1
}
export declare function accessTypeFromJSON(object: any): AccessType;
export declare function accessTypeToJSON(object: AccessType): string;
export interface AccessTypeParam {
    value: AccessType;
}
/** CodeInfo is data for the uploaded contract WASM code */
export interface CodeInfo {
    codeHash: Uint8Array;
    creator: Uint8Array;
    source: string;
    builder: string;
    /**
     * AccessConfig instantiate_config = 5 [(gogoproto.nullable) = false];
     *   google.protobuf.Timestamp end_time = 5;
     */
    endTime: Duration | undefined;
}
/** ContractInfo stores a WASM contract instance */
export interface ContractInfo {
    codeId: number;
    creator: Uint8Array;
    /** bytes admin = 3 [(gogoproto.casttype) = "github.com/cosmos/cosmos-sdk/types.AccAddress"]; */
    contractId: string;
    /**
     * never show this in query results, just use for sorting
     * (Note: when using json tag "-" amino refused to serialize it...)
     */
    created: AbsoluteTxPosition | undefined;
    endTime: Date | undefined;
    /**
     * bytes init_msg = 5 [(gogoproto.casttype) = "encoding/json.RawMessage"];
     *
     *    AbsoluteTxPosition last_updated = 7;
     *    uint64 previous_code_id = 8 [(gogoproto.customname) = "PreviousCodeID"];
     */
    lastMsg: Uint8Array;
}
/** ContractInfoWithAddress adds the address (key) to the ContractInfo representation */
export interface ContractInfoWithAddress {
    address: Uint8Array;
    ContractInfo: ContractInfo | undefined;
}
/** AbsoluteTxPosition can be used to sort contracts */
export interface AbsoluteTxPosition {
    /** BlockHeight is the block the contract was created at */
    blockHeight: number;
    /** TxIndex is a monotonic counter within the block (actual transaction index, or gas consumed) */
    txIndex: number;
}
/** Model is a struct that holds a KV pair */
export interface Model {
    /** hex-encode key to read it better (this is often ascii) */
    Key: Uint8Array;
    /** base64-encode raw value */
    Value: Uint8Array;
}
export declare const AccessTypeParam: {
    encode(message: AccessTypeParam, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): AccessTypeParam;
    fromJSON(object: any): AccessTypeParam;
    toJSON(message: AccessTypeParam): unknown;
    fromPartial(object: DeepPartial<AccessTypeParam>): AccessTypeParam;
};
export declare const CodeInfo: {
    encode(message: CodeInfo, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): CodeInfo;
    fromJSON(object: any): CodeInfo;
    toJSON(message: CodeInfo): unknown;
    fromPartial(object: DeepPartial<CodeInfo>): CodeInfo;
};
export declare const ContractInfo: {
    encode(message: ContractInfo, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): ContractInfo;
    fromJSON(object: any): ContractInfo;
    toJSON(message: ContractInfo): unknown;
    fromPartial(object: DeepPartial<ContractInfo>): ContractInfo;
};
export declare const ContractInfoWithAddress: {
    encode(message: ContractInfoWithAddress, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): ContractInfoWithAddress;
    fromJSON(object: any): ContractInfoWithAddress;
    toJSON(message: ContractInfoWithAddress): unknown;
    fromPartial(object: DeepPartial<ContractInfoWithAddress>): ContractInfoWithAddress;
};
export declare const AbsoluteTxPosition: {
    encode(message: AbsoluteTxPosition, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): AbsoluteTxPosition;
    fromJSON(object: any): AbsoluteTxPosition;
    toJSON(message: AbsoluteTxPosition): unknown;
    fromPartial(object: DeepPartial<AbsoluteTxPosition>): AbsoluteTxPosition;
};
export declare const Model: {
    encode(message: Model, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Model;
    fromJSON(object: any): Model;
    toJSON(message: Model): unknown;
    fromPartial(object: DeepPartial<Model>): Model;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
