import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "trst.x.registration.v1beta1";
export interface SeedConfig {
    masterCert: string;
    encryptedKey: string;
}
export interface RegistrationNodeInfo {
    certificate: Uint8Array;
    encryptedSeed: Uint8Array;
}
export declare const SeedConfig: {
    encode(message: SeedConfig, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SeedConfig;
    fromJSON(object: any): SeedConfig;
    toJSON(message: SeedConfig): unknown;
    fromPartial(object: DeepPartial<SeedConfig>): SeedConfig;
};
export declare const RegistrationNodeInfo: {
    encode(message: RegistrationNodeInfo, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): RegistrationNodeInfo;
    fromJSON(object: any): RegistrationNodeInfo;
    toJSON(message: RegistrationNodeInfo): unknown;
    fromPartial(object: DeepPartial<RegistrationNodeInfo>): RegistrationNodeInfo;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
