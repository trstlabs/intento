import { Reader, Writer } from "protobufjs/minimal";
import { MasterCertificate } from "../registration/msg";
export declare const protobufPackage = "trst.x.registration.v1beta1";
export interface QueryMasterKeyRequest {
}
export interface QueryMasterKeyResponse {
    /** [(gogoproto.nullable) = false]; */
    masterKey: MasterCertificate | undefined;
}
export interface QueryEncryptedSeedRequest {
    pubKey: Uint8Array;
}
export interface QueryEncryptedSeedResponse {
    /** [(gogoproto.nullable) = false]; */
    encryptedSeed: Uint8Array;
}
export declare const QueryMasterKeyRequest: {
    encode(_: QueryMasterKeyRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryMasterKeyRequest;
    fromJSON(_: any): QueryMasterKeyRequest;
    toJSON(_: QueryMasterKeyRequest): unknown;
    fromPartial(_: DeepPartial<QueryMasterKeyRequest>): QueryMasterKeyRequest;
};
export declare const QueryMasterKeyResponse: {
    encode(message: QueryMasterKeyResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryMasterKeyResponse;
    fromJSON(object: any): QueryMasterKeyResponse;
    toJSON(message: QueryMasterKeyResponse): unknown;
    fromPartial(object: DeepPartial<QueryMasterKeyResponse>): QueryMasterKeyResponse;
};
export declare const QueryEncryptedSeedRequest: {
    encode(message: QueryEncryptedSeedRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryEncryptedSeedRequest;
    fromJSON(object: any): QueryEncryptedSeedRequest;
    toJSON(message: QueryEncryptedSeedRequest): unknown;
    fromPartial(object: DeepPartial<QueryEncryptedSeedRequest>): QueryEncryptedSeedRequest;
};
export declare const QueryEncryptedSeedResponse: {
    encode(message: QueryEncryptedSeedResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryEncryptedSeedResponse;
    fromJSON(object: any): QueryEncryptedSeedResponse;
    toJSON(message: QueryEncryptedSeedResponse): unknown;
    fromPartial(object: DeepPartial<QueryEncryptedSeedResponse>): QueryEncryptedSeedResponse;
};
/** Query provides defines the gRPC querier service */
export interface Query {
    MasterKey(request: QueryMasterKeyRequest): Promise<QueryMasterKeyResponse>;
    EncryptedSeed(request: QueryEncryptedSeedRequest): Promise<QueryEncryptedSeedResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    MasterKey(request: QueryMasterKeyRequest): Promise<QueryMasterKeyResponse>;
    EncryptedSeed(request: QueryEncryptedSeedRequest): Promise<QueryEncryptedSeedResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
