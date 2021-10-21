import { Reader, Writer } from "protobufjs/minimal";
import { ContractInfo, ContractInfoWithAddress } from "../compute/types";
import { Duration } from "../google/protobuf/duration";
import { StringEvent } from "../cosmos/base/abci/v1beta1/abci";
import { Empty } from "../google/protobuf/empty";
export declare const protobufPackage = "trst.x.compute.v1beta1";
/** QueryContractInfoRequest is the request type for the Query/ContractInfo RPC method */
export interface QueryContractInfoRequest {
    /** address is the address of the contract to query */
    address: Uint8Array;
}
/** QueryContractInfoResponse is the response type for the Query/ContractInfo RPC method */
export interface QueryContractInfoResponse {
    /** address is the address of the contract */
    address: Uint8Array;
    ContractInfo: ContractInfo | undefined;
}
/** QueryContractResultRequest is the request type for the Query/ContractrResult RPC method */
export interface QueryContractResultRequest {
    /** address is the address of the contract to query */
    address: Uint8Array;
}
/** QueryContractResultResponse is the response type for the Query/ContractrResult RPC method */
export interface QueryContractResultResponse {
    /** address is the address of the contract */
    address: Uint8Array;
    data: Uint8Array;
    log: string;
}
export interface QueryContractHistoryRequest {
    /** address is the address of the contract to query */
    address: Uint8Array;
}
export interface QueryContractsByCodeRequest {
    /** grpc-gateway_out does not support Go style CodID */
    codeId: number;
}
export interface QueryContractsByCodeResponse {
    contractInfos: ContractInfoWithAddress[];
}
export interface QuerySmartContractStateRequest {
    /** address is the address of the contract */
    address: Uint8Array;
    queryData: Uint8Array;
}
export interface QueryContractAddressByContractIdRequest {
    contractId: string;
}
export interface QueryContractKeyRequest {
    /** address is the address of the contract */
    address: Uint8Array;
}
export interface QueryContractHashRequest {
    /** address is the address of the contract */
    address: Uint8Array;
}
export interface QuerySmartContractStateResponse {
    data: Uint8Array;
}
export interface QueryCodeRequest {
    /** grpc-gateway_out does not support Go style CodID */
    codeId: number;
}
export interface CodeInfoResponse {
    /** id for legacy support */
    codeId: number;
    creator: Uint8Array;
    codeHash: Uint8Array;
    source: string;
    builder: string;
    contractDuration: Duration | undefined;
}
export interface QueryCodeResponse {
    codeInfo: CodeInfoResponse | undefined;
    data: Uint8Array;
}
export interface QueryCodesResponse {
    codeInfos: CodeInfoResponse[];
}
export interface QueryContractAddressByContractIdResponse {
    /** address is the address of the contract */
    address: Uint8Array;
}
export interface QueryContractKeyResponse {
    /** address is the address of the contract */
    key: Uint8Array;
}
export interface QueryContractHashResponse {
    codeHash: Uint8Array;
}
/** DecryptedAnswer is a struct that represents a decrypted tx-query */
export interface DecryptedAnswer {
    type: string;
    input: string;
    outputData: string;
    outputDataAsString: string;
    outputLogs: StringEvent[];
    outputError: Uint8Array;
    plaintextError: string;
}
export declare const QueryContractInfoRequest: {
    encode(message: QueryContractInfoRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractInfoRequest;
    fromJSON(object: any): QueryContractInfoRequest;
    toJSON(message: QueryContractInfoRequest): unknown;
    fromPartial(object: DeepPartial<QueryContractInfoRequest>): QueryContractInfoRequest;
};
export declare const QueryContractInfoResponse: {
    encode(message: QueryContractInfoResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractInfoResponse;
    fromJSON(object: any): QueryContractInfoResponse;
    toJSON(message: QueryContractInfoResponse): unknown;
    fromPartial(object: DeepPartial<QueryContractInfoResponse>): QueryContractInfoResponse;
};
export declare const QueryContractResultRequest: {
    encode(message: QueryContractResultRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractResultRequest;
    fromJSON(object: any): QueryContractResultRequest;
    toJSON(message: QueryContractResultRequest): unknown;
    fromPartial(object: DeepPartial<QueryContractResultRequest>): QueryContractResultRequest;
};
export declare const QueryContractResultResponse: {
    encode(message: QueryContractResultResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractResultResponse;
    fromJSON(object: any): QueryContractResultResponse;
    toJSON(message: QueryContractResultResponse): unknown;
    fromPartial(object: DeepPartial<QueryContractResultResponse>): QueryContractResultResponse;
};
export declare const QueryContractHistoryRequest: {
    encode(message: QueryContractHistoryRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractHistoryRequest;
    fromJSON(object: any): QueryContractHistoryRequest;
    toJSON(message: QueryContractHistoryRequest): unknown;
    fromPartial(object: DeepPartial<QueryContractHistoryRequest>): QueryContractHistoryRequest;
};
export declare const QueryContractsByCodeRequest: {
    encode(message: QueryContractsByCodeRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractsByCodeRequest;
    fromJSON(object: any): QueryContractsByCodeRequest;
    toJSON(message: QueryContractsByCodeRequest): unknown;
    fromPartial(object: DeepPartial<QueryContractsByCodeRequest>): QueryContractsByCodeRequest;
};
export declare const QueryContractsByCodeResponse: {
    encode(message: QueryContractsByCodeResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractsByCodeResponse;
    fromJSON(object: any): QueryContractsByCodeResponse;
    toJSON(message: QueryContractsByCodeResponse): unknown;
    fromPartial(object: DeepPartial<QueryContractsByCodeResponse>): QueryContractsByCodeResponse;
};
export declare const QuerySmartContractStateRequest: {
    encode(message: QuerySmartContractStateRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QuerySmartContractStateRequest;
    fromJSON(object: any): QuerySmartContractStateRequest;
    toJSON(message: QuerySmartContractStateRequest): unknown;
    fromPartial(object: DeepPartial<QuerySmartContractStateRequest>): QuerySmartContractStateRequest;
};
export declare const QueryContractAddressByContractIdRequest: {
    encode(message: QueryContractAddressByContractIdRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractAddressByContractIdRequest;
    fromJSON(object: any): QueryContractAddressByContractIdRequest;
    toJSON(message: QueryContractAddressByContractIdRequest): unknown;
    fromPartial(object: DeepPartial<QueryContractAddressByContractIdRequest>): QueryContractAddressByContractIdRequest;
};
export declare const QueryContractKeyRequest: {
    encode(message: QueryContractKeyRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractKeyRequest;
    fromJSON(object: any): QueryContractKeyRequest;
    toJSON(message: QueryContractKeyRequest): unknown;
    fromPartial(object: DeepPartial<QueryContractKeyRequest>): QueryContractKeyRequest;
};
export declare const QueryContractHashRequest: {
    encode(message: QueryContractHashRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractHashRequest;
    fromJSON(object: any): QueryContractHashRequest;
    toJSON(message: QueryContractHashRequest): unknown;
    fromPartial(object: DeepPartial<QueryContractHashRequest>): QueryContractHashRequest;
};
export declare const QuerySmartContractStateResponse: {
    encode(message: QuerySmartContractStateResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QuerySmartContractStateResponse;
    fromJSON(object: any): QuerySmartContractStateResponse;
    toJSON(message: QuerySmartContractStateResponse): unknown;
    fromPartial(object: DeepPartial<QuerySmartContractStateResponse>): QuerySmartContractStateResponse;
};
export declare const QueryCodeRequest: {
    encode(message: QueryCodeRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryCodeRequest;
    fromJSON(object: any): QueryCodeRequest;
    toJSON(message: QueryCodeRequest): unknown;
    fromPartial(object: DeepPartial<QueryCodeRequest>): QueryCodeRequest;
};
export declare const CodeInfoResponse: {
    encode(message: CodeInfoResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): CodeInfoResponse;
    fromJSON(object: any): CodeInfoResponse;
    toJSON(message: CodeInfoResponse): unknown;
    fromPartial(object: DeepPartial<CodeInfoResponse>): CodeInfoResponse;
};
export declare const QueryCodeResponse: {
    encode(message: QueryCodeResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryCodeResponse;
    fromJSON(object: any): QueryCodeResponse;
    toJSON(message: QueryCodeResponse): unknown;
    fromPartial(object: DeepPartial<QueryCodeResponse>): QueryCodeResponse;
};
export declare const QueryCodesResponse: {
    encode(message: QueryCodesResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryCodesResponse;
    fromJSON(object: any): QueryCodesResponse;
    toJSON(message: QueryCodesResponse): unknown;
    fromPartial(object: DeepPartial<QueryCodesResponse>): QueryCodesResponse;
};
export declare const QueryContractAddressByContractIdResponse: {
    encode(message: QueryContractAddressByContractIdResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractAddressByContractIdResponse;
    fromJSON(object: any): QueryContractAddressByContractIdResponse;
    toJSON(message: QueryContractAddressByContractIdResponse): unknown;
    fromPartial(object: DeepPartial<QueryContractAddressByContractIdResponse>): QueryContractAddressByContractIdResponse;
};
export declare const QueryContractKeyResponse: {
    encode(message: QueryContractKeyResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractKeyResponse;
    fromJSON(object: any): QueryContractKeyResponse;
    toJSON(message: QueryContractKeyResponse): unknown;
    fromPartial(object: DeepPartial<QueryContractKeyResponse>): QueryContractKeyResponse;
};
export declare const QueryContractHashResponse: {
    encode(message: QueryContractHashResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryContractHashResponse;
    fromJSON(object: any): QueryContractHashResponse;
    toJSON(message: QueryContractHashResponse): unknown;
    fromPartial(object: DeepPartial<QueryContractHashResponse>): QueryContractHashResponse;
};
export declare const DecryptedAnswer: {
    encode(message: DecryptedAnswer, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): DecryptedAnswer;
    fromJSON(object: any): DecryptedAnswer;
    toJSON(message: DecryptedAnswer): unknown;
    fromPartial(object: DeepPartial<DecryptedAnswer>): DecryptedAnswer;
};
/** Query provides defines the gRPC querier service */
export interface Query {
    ContractInfo(request: QueryContractInfoRequest): Promise<QueryContractInfoResponse>;
    ContractResult(request: QueryContractResultRequest): Promise<QueryContractResultResponse>;
    /**
     * rpc ContractHistory (QueryContractHistoryRequest) returns (QueryContractHistoryResponse) {
     * option (google.api.http).get = "/compute/v1beta1/contract/{address}/history";
     * }
     */
    ContractsByCode(request: QueryContractsByCodeRequest): Promise<QueryContractsByCodeResponse>;
    /**
     * rpc AllContractState (QueryAllContractStateRequest) returns (QueryAllContractStateResponse) {
     * option (google.api.http).get = "/compute/v1beta1/contract/{address}/state";
     * }
     * rpc RawContractState (QueryRawContractStateRequest) returns (QueryRawContractStateResponse) {
     * option (google.api.http).get = "/compute/v1beta1/contract/{address}/raw/{query_data}";
     * }
     */
    SmartContractState(request: QuerySmartContractStateRequest): Promise<QuerySmartContractStateResponse>;
    Code(request: QueryCodeRequest): Promise<QueryCodeResponse>;
    Codes(request: Empty): Promise<QueryCodesResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    ContractInfo(request: QueryContractInfoRequest): Promise<QueryContractInfoResponse>;
    ContractResult(request: QueryContractResultRequest): Promise<QueryContractResultResponse>;
    ContractsByCode(request: QueryContractsByCodeRequest): Promise<QueryContractsByCodeResponse>;
    SmartContractState(request: QuerySmartContractStateRequest): Promise<QuerySmartContractStateResponse>;
    Code(request: QueryCodeRequest): Promise<QueryCodeResponse>;
    Codes(request: Empty): Promise<QueryCodesResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
