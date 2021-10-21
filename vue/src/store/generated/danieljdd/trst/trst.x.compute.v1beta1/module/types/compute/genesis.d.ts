import { Writer, Reader } from "protobufjs/minimal";
import { CodeInfo, ContractInfo, Model } from "../compute/types";
export declare const protobufPackage = "trst.x.compute.v1beta1";
/** GenesisState - genesis state of x/wasm */
export interface GenesisState {
    /** Params params = 1 [(gogoproto.nullable) = false]; */
    codes: Code[];
    contracts: Contract[];
    sequences: Sequence[];
}
/** Code struct encompasses CodeInfo and CodeBytes */
export interface Code {
    codeId: number;
    codeInfo: CodeInfo | undefined;
    codeBytes: Uint8Array;
}
/** Contract struct encompasses ContractAddress, ContractInfo, and ContractState */
export interface Contract {
    contractAddress: Uint8Array;
    contractInfo: ContractInfo | undefined;
    contractState: Model[];
}
/** Sequence id and value of a counter */
export interface Sequence {
    idKey: Uint8Array;
    value: number;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
export declare const Code: {
    encode(message: Code, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Code;
    fromJSON(object: any): Code;
    toJSON(message: Code): unknown;
    fromPartial(object: DeepPartial<Code>): Code;
};
export declare const Contract: {
    encode(message: Contract, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Contract;
    fromJSON(object: any): Contract;
    toJSON(message: Contract): unknown;
    fromPartial(object: DeepPartial<Contract>): Contract;
};
export declare const Sequence: {
    encode(message: Sequence, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Sequence;
    fromJSON(object: any): Sequence;
    toJSON(message: Sequence): unknown;
    fromPartial(object: DeepPartial<Sequence>): Sequence;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
