/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { MasterCertificate } from "../registration/msg";

export const protobufPackage = "trst.x.registration.v1beta1";

export interface QueryMasterKeyRequest {}

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

const baseQueryMasterKeyRequest: object = {};

export const QueryMasterKeyRequest = {
  encode(_: QueryMasterKeyRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryMasterKeyRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryMasterKeyRequest } as QueryMasterKeyRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryMasterKeyRequest {
    const message = { ...baseQueryMasterKeyRequest } as QueryMasterKeyRequest;
    return message;
  },

  toJSON(_: QueryMasterKeyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<QueryMasterKeyRequest>): QueryMasterKeyRequest {
    const message = { ...baseQueryMasterKeyRequest } as QueryMasterKeyRequest;
    return message;
  },
};

const baseQueryMasterKeyResponse: object = {};

export const QueryMasterKeyResponse = {
  encode(
    message: QueryMasterKeyResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.masterKey !== undefined) {
      MasterCertificate.encode(
        message.masterKey,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryMasterKeyResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryMasterKeyResponse } as QueryMasterKeyResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.masterKey = MasterCertificate.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryMasterKeyResponse {
    const message = { ...baseQueryMasterKeyResponse } as QueryMasterKeyResponse;
    if (object.masterKey !== undefined && object.masterKey !== null) {
      message.masterKey = MasterCertificate.fromJSON(object.masterKey);
    } else {
      message.masterKey = undefined;
    }
    return message;
  },

  toJSON(message: QueryMasterKeyResponse): unknown {
    const obj: any = {};
    message.masterKey !== undefined &&
      (obj.masterKey = message.masterKey
        ? MasterCertificate.toJSON(message.masterKey)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryMasterKeyResponse>
  ): QueryMasterKeyResponse {
    const message = { ...baseQueryMasterKeyResponse } as QueryMasterKeyResponse;
    if (object.masterKey !== undefined && object.masterKey !== null) {
      message.masterKey = MasterCertificate.fromPartial(object.masterKey);
    } else {
      message.masterKey = undefined;
    }
    return message;
  },
};

const baseQueryEncryptedSeedRequest: object = {};

export const QueryEncryptedSeedRequest = {
  encode(
    message: QueryEncryptedSeedRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pubKey.length !== 0) {
      writer.uint32(10).bytes(message.pubKey);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryEncryptedSeedRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryEncryptedSeedRequest,
    } as QueryEncryptedSeedRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pubKey = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryEncryptedSeedRequest {
    const message = {
      ...baseQueryEncryptedSeedRequest,
    } as QueryEncryptedSeedRequest;
    if (object.pubKey !== undefined && object.pubKey !== null) {
      message.pubKey = bytesFromBase64(object.pubKey);
    }
    return message;
  },

  toJSON(message: QueryEncryptedSeedRequest): unknown {
    const obj: any = {};
    message.pubKey !== undefined &&
      (obj.pubKey = base64FromBytes(
        message.pubKey !== undefined ? message.pubKey : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryEncryptedSeedRequest>
  ): QueryEncryptedSeedRequest {
    const message = {
      ...baseQueryEncryptedSeedRequest,
    } as QueryEncryptedSeedRequest;
    if (object.pubKey !== undefined && object.pubKey !== null) {
      message.pubKey = object.pubKey;
    } else {
      message.pubKey = new Uint8Array();
    }
    return message;
  },
};

const baseQueryEncryptedSeedResponse: object = {};

export const QueryEncryptedSeedResponse = {
  encode(
    message: QueryEncryptedSeedResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.encryptedSeed.length !== 0) {
      writer.uint32(10).bytes(message.encryptedSeed);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryEncryptedSeedResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryEncryptedSeedResponse,
    } as QueryEncryptedSeedResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.encryptedSeed = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryEncryptedSeedResponse {
    const message = {
      ...baseQueryEncryptedSeedResponse,
    } as QueryEncryptedSeedResponse;
    if (object.encryptedSeed !== undefined && object.encryptedSeed !== null) {
      message.encryptedSeed = bytesFromBase64(object.encryptedSeed);
    }
    return message;
  },

  toJSON(message: QueryEncryptedSeedResponse): unknown {
    const obj: any = {};
    message.encryptedSeed !== undefined &&
      (obj.encryptedSeed = base64FromBytes(
        message.encryptedSeed !== undefined
          ? message.encryptedSeed
          : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryEncryptedSeedResponse>
  ): QueryEncryptedSeedResponse {
    const message = {
      ...baseQueryEncryptedSeedResponse,
    } as QueryEncryptedSeedResponse;
    if (object.encryptedSeed !== undefined && object.encryptedSeed !== null) {
      message.encryptedSeed = object.encryptedSeed;
    } else {
      message.encryptedSeed = new Uint8Array();
    }
    return message;
  },
};

/** Query provides defines the gRPC querier service */
export interface Query {
  MasterKey(request: QueryMasterKeyRequest): Promise<QueryMasterKeyResponse>;
  EncryptedSeed(
    request: QueryEncryptedSeedRequest
  ): Promise<QueryEncryptedSeedResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  MasterKey(request: QueryMasterKeyRequest): Promise<QueryMasterKeyResponse> {
    const data = QueryMasterKeyRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.registration.v1beta1.Query",
      "MasterKey",
      data
    );
    return promise.then((data) =>
      QueryMasterKeyResponse.decode(new Reader(data))
    );
  }

  EncryptedSeed(
    request: QueryEncryptedSeedRequest
  ): Promise<QueryEncryptedSeedResponse> {
    const data = QueryEncryptedSeedRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.registration.v1beta1.Query",
      "EncryptedSeed",
      data
    );
    return promise.then((data) =>
      QueryEncryptedSeedResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  for (let i = 0; i < arr.byteLength; ++i) {
    bin.push(String.fromCharCode(arr[i]));
  }
  return btoa(bin.join(""));
}

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
