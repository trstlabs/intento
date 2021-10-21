import { RegistrationNodeInfo } from "../registration/types";
import { MasterCertificate } from "../registration/msg";
import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "trst.x.registration.v1beta1";
export interface GenesisState {
    registration: RegistrationNodeInfo[];
    nodeExchMasterCertificate: MasterCertificate | undefined;
    ioMasterCertificate: MasterCertificate | undefined;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
