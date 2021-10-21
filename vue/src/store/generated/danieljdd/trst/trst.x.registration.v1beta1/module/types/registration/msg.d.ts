import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "trst.x.registration.v1beta1";
export interface RaAuthenticate {
    sender: Uint8Array;
    certificate: Uint8Array;
}
export interface MasterCertificate {
    bytes: Uint8Array;
}
export declare const RaAuthenticate: {
    encode(message: RaAuthenticate, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): RaAuthenticate;
    fromJSON(object: any): RaAuthenticate;
    toJSON(message: RaAuthenticate): unknown;
    fromPartial(object: DeepPartial<RaAuthenticate>): RaAuthenticate;
};
export declare const MasterCertificate: {
    encode(message: MasterCertificate, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MasterCertificate;
    fromJSON(object: any): MasterCertificate;
    toJSON(message: MasterCertificate): unknown;
    fromPartial(object: DeepPartial<MasterCertificate>): MasterCertificate;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
