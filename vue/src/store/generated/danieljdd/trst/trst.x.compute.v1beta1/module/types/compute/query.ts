/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { ContractInfo, ContractInfoWithAddress } from "../compute/types";
import { Duration } from "../google/protobuf/duration";
import { StringEvent } from "../cosmos/base/abci/v1beta1/abci";
import { Empty } from "../google/protobuf/empty";

export const protobufPackage = "trst.x.compute.v1beta1";

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

const baseQueryContractInfoRequest: object = {};

export const QueryContractInfoRequest = {
  encode(
    message: QueryContractInfoRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractInfoRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractInfoRequest,
    } as QueryContractInfoRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractInfoRequest {
    const message = {
      ...baseQueryContractInfoRequest,
    } as QueryContractInfoRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    return message;
  },

  toJSON(message: QueryContractInfoRequest): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractInfoRequest>
  ): QueryContractInfoRequest {
    const message = {
      ...baseQueryContractInfoRequest,
    } as QueryContractInfoRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    return message;
  },
};

const baseQueryContractInfoResponse: object = {};

export const QueryContractInfoResponse = {
  encode(
    message: QueryContractInfoResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    if (message.ContractInfo !== undefined) {
      ContractInfo.encode(
        message.ContractInfo,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractInfoResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractInfoResponse,
    } as QueryContractInfoResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        case 2:
          message.ContractInfo = ContractInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractInfoResponse {
    const message = {
      ...baseQueryContractInfoResponse,
    } as QueryContractInfoResponse;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    if (object.ContractInfo !== undefined && object.ContractInfo !== null) {
      message.ContractInfo = ContractInfo.fromJSON(object.ContractInfo);
    } else {
      message.ContractInfo = undefined;
    }
    return message;
  },

  toJSON(message: QueryContractInfoResponse): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    message.ContractInfo !== undefined &&
      (obj.ContractInfo = message.ContractInfo
        ? ContractInfo.toJSON(message.ContractInfo)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractInfoResponse>
  ): QueryContractInfoResponse {
    const message = {
      ...baseQueryContractInfoResponse,
    } as QueryContractInfoResponse;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    if (object.ContractInfo !== undefined && object.ContractInfo !== null) {
      message.ContractInfo = ContractInfo.fromPartial(object.ContractInfo);
    } else {
      message.ContractInfo = undefined;
    }
    return message;
  },
};

const baseQueryContractResultRequest: object = {};

export const QueryContractResultRequest = {
  encode(
    message: QueryContractResultRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractResultRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractResultRequest,
    } as QueryContractResultRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractResultRequest {
    const message = {
      ...baseQueryContractResultRequest,
    } as QueryContractResultRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    return message;
  },

  toJSON(message: QueryContractResultRequest): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractResultRequest>
  ): QueryContractResultRequest {
    const message = {
      ...baseQueryContractResultRequest,
    } as QueryContractResultRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    return message;
  },
};

const baseQueryContractResultResponse: object = { log: "" };

export const QueryContractResultResponse = {
  encode(
    message: QueryContractResultResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    if (message.log !== "") {
      writer.uint32(26).string(message.log);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractResultResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractResultResponse,
    } as QueryContractResultResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        case 2:
          message.data = reader.bytes();
          break;
        case 3:
          message.log = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractResultResponse {
    const message = {
      ...baseQueryContractResultResponse,
    } as QueryContractResultResponse;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    if (object.data !== undefined && object.data !== null) {
      message.data = bytesFromBase64(object.data);
    }
    if (object.log !== undefined && object.log !== null) {
      message.log = String(object.log);
    } else {
      message.log = "";
    }
    return message;
  },

  toJSON(message: QueryContractResultResponse): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : new Uint8Array()
      ));
    message.log !== undefined && (obj.log = message.log);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractResultResponse>
  ): QueryContractResultResponse {
    const message = {
      ...baseQueryContractResultResponse,
    } as QueryContractResultResponse;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    if (object.data !== undefined && object.data !== null) {
      message.data = object.data;
    } else {
      message.data = new Uint8Array();
    }
    if (object.log !== undefined && object.log !== null) {
      message.log = object.log;
    } else {
      message.log = "";
    }
    return message;
  },
};

const baseQueryContractHistoryRequest: object = {};

export const QueryContractHistoryRequest = {
  encode(
    message: QueryContractHistoryRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractHistoryRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractHistoryRequest,
    } as QueryContractHistoryRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractHistoryRequest {
    const message = {
      ...baseQueryContractHistoryRequest,
    } as QueryContractHistoryRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    return message;
  },

  toJSON(message: QueryContractHistoryRequest): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractHistoryRequest>
  ): QueryContractHistoryRequest {
    const message = {
      ...baseQueryContractHistoryRequest,
    } as QueryContractHistoryRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    return message;
  },
};

const baseQueryContractsByCodeRequest: object = { codeId: 0 };

export const QueryContractsByCodeRequest = {
  encode(
    message: QueryContractsByCodeRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.codeId !== 0) {
      writer.uint32(8).uint64(message.codeId);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractsByCodeRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractsByCodeRequest,
    } as QueryContractsByCodeRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractsByCodeRequest {
    const message = {
      ...baseQueryContractsByCodeRequest,
    } as QueryContractsByCodeRequest;
    if (object.codeId !== undefined && object.codeId !== null) {
      message.codeId = Number(object.codeId);
    } else {
      message.codeId = 0;
    }
    return message;
  },

  toJSON(message: QueryContractsByCodeRequest): unknown {
    const obj: any = {};
    message.codeId !== undefined && (obj.codeId = message.codeId);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractsByCodeRequest>
  ): QueryContractsByCodeRequest {
    const message = {
      ...baseQueryContractsByCodeRequest,
    } as QueryContractsByCodeRequest;
    if (object.codeId !== undefined && object.codeId !== null) {
      message.codeId = object.codeId;
    } else {
      message.codeId = 0;
    }
    return message;
  },
};

const baseQueryContractsByCodeResponse: object = {};

export const QueryContractsByCodeResponse = {
  encode(
    message: QueryContractsByCodeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.contractInfos) {
      ContractInfoWithAddress.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractsByCodeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractsByCodeResponse,
    } as QueryContractsByCodeResponse;
    message.contractInfos = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contractInfos.push(
            ContractInfoWithAddress.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractsByCodeResponse {
    const message = {
      ...baseQueryContractsByCodeResponse,
    } as QueryContractsByCodeResponse;
    message.contractInfos = [];
    if (object.contractInfos !== undefined && object.contractInfos !== null) {
      for (const e of object.contractInfos) {
        message.contractInfos.push(ContractInfoWithAddress.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: QueryContractsByCodeResponse): unknown {
    const obj: any = {};
    if (message.contractInfos) {
      obj.contractInfos = message.contractInfos.map((e) =>
        e ? ContractInfoWithAddress.toJSON(e) : undefined
      );
    } else {
      obj.contractInfos = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractsByCodeResponse>
  ): QueryContractsByCodeResponse {
    const message = {
      ...baseQueryContractsByCodeResponse,
    } as QueryContractsByCodeResponse;
    message.contractInfos = [];
    if (object.contractInfos !== undefined && object.contractInfos !== null) {
      for (const e of object.contractInfos) {
        message.contractInfos.push(ContractInfoWithAddress.fromPartial(e));
      }
    }
    return message;
  },
};

const baseQuerySmartContractStateRequest: object = {};

export const QuerySmartContractStateRequest = {
  encode(
    message: QuerySmartContractStateRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    if (message.queryData.length !== 0) {
      writer.uint32(18).bytes(message.queryData);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QuerySmartContractStateRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQuerySmartContractStateRequest,
    } as QuerySmartContractStateRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        case 2:
          message.queryData = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QuerySmartContractStateRequest {
    const message = {
      ...baseQuerySmartContractStateRequest,
    } as QuerySmartContractStateRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    if (object.queryData !== undefined && object.queryData !== null) {
      message.queryData = bytesFromBase64(object.queryData);
    }
    return message;
  },

  toJSON(message: QuerySmartContractStateRequest): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    message.queryData !== undefined &&
      (obj.queryData = base64FromBytes(
        message.queryData !== undefined ? message.queryData : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QuerySmartContractStateRequest>
  ): QuerySmartContractStateRequest {
    const message = {
      ...baseQuerySmartContractStateRequest,
    } as QuerySmartContractStateRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    if (object.queryData !== undefined && object.queryData !== null) {
      message.queryData = object.queryData;
    } else {
      message.queryData = new Uint8Array();
    }
    return message;
  },
};

const baseQueryContractAddressByContractIdRequest: object = { contractId: "" };

export const QueryContractAddressByContractIdRequest = {
  encode(
    message: QueryContractAddressByContractIdRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.contractId !== "") {
      writer.uint32(10).string(message.contractId);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractAddressByContractIdRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractAddressByContractIdRequest,
    } as QueryContractAddressByContractIdRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contractId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractAddressByContractIdRequest {
    const message = {
      ...baseQueryContractAddressByContractIdRequest,
    } as QueryContractAddressByContractIdRequest;
    if (object.contractId !== undefined && object.contractId !== null) {
      message.contractId = String(object.contractId);
    } else {
      message.contractId = "";
    }
    return message;
  },

  toJSON(message: QueryContractAddressByContractIdRequest): unknown {
    const obj: any = {};
    message.contractId !== undefined && (obj.contractId = message.contractId);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractAddressByContractIdRequest>
  ): QueryContractAddressByContractIdRequest {
    const message = {
      ...baseQueryContractAddressByContractIdRequest,
    } as QueryContractAddressByContractIdRequest;
    if (object.contractId !== undefined && object.contractId !== null) {
      message.contractId = object.contractId;
    } else {
      message.contractId = "";
    }
    return message;
  },
};

const baseQueryContractKeyRequest: object = {};

export const QueryContractKeyRequest = {
  encode(
    message: QueryContractKeyRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryContractKeyRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractKeyRequest,
    } as QueryContractKeyRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractKeyRequest {
    const message = {
      ...baseQueryContractKeyRequest,
    } as QueryContractKeyRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    return message;
  },

  toJSON(message: QueryContractKeyRequest): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractKeyRequest>
  ): QueryContractKeyRequest {
    const message = {
      ...baseQueryContractKeyRequest,
    } as QueryContractKeyRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    return message;
  },
};

const baseQueryContractHashRequest: object = {};

export const QueryContractHashRequest = {
  encode(
    message: QueryContractHashRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractHashRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractHashRequest,
    } as QueryContractHashRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractHashRequest {
    const message = {
      ...baseQueryContractHashRequest,
    } as QueryContractHashRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    return message;
  },

  toJSON(message: QueryContractHashRequest): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractHashRequest>
  ): QueryContractHashRequest {
    const message = {
      ...baseQueryContractHashRequest,
    } as QueryContractHashRequest;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    return message;
  },
};

const baseQuerySmartContractStateResponse: object = {};

export const QuerySmartContractStateResponse = {
  encode(
    message: QuerySmartContractStateResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QuerySmartContractStateResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQuerySmartContractStateResponse,
    } as QuerySmartContractStateResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QuerySmartContractStateResponse {
    const message = {
      ...baseQuerySmartContractStateResponse,
    } as QuerySmartContractStateResponse;
    if (object.data !== undefined && object.data !== null) {
      message.data = bytesFromBase64(object.data);
    }
    return message;
  },

  toJSON(message: QuerySmartContractStateResponse): unknown {
    const obj: any = {};
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QuerySmartContractStateResponse>
  ): QuerySmartContractStateResponse {
    const message = {
      ...baseQuerySmartContractStateResponse,
    } as QuerySmartContractStateResponse;
    if (object.data !== undefined && object.data !== null) {
      message.data = object.data;
    } else {
      message.data = new Uint8Array();
    }
    return message;
  },
};

const baseQueryCodeRequest: object = { codeId: 0 };

export const QueryCodeRequest = {
  encode(message: QueryCodeRequest, writer: Writer = Writer.create()): Writer {
    if (message.codeId !== 0) {
      writer.uint32(8).uint64(message.codeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryCodeRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryCodeRequest } as QueryCodeRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryCodeRequest {
    const message = { ...baseQueryCodeRequest } as QueryCodeRequest;
    if (object.codeId !== undefined && object.codeId !== null) {
      message.codeId = Number(object.codeId);
    } else {
      message.codeId = 0;
    }
    return message;
  },

  toJSON(message: QueryCodeRequest): unknown {
    const obj: any = {};
    message.codeId !== undefined && (obj.codeId = message.codeId);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryCodeRequest>): QueryCodeRequest {
    const message = { ...baseQueryCodeRequest } as QueryCodeRequest;
    if (object.codeId !== undefined && object.codeId !== null) {
      message.codeId = object.codeId;
    } else {
      message.codeId = 0;
    }
    return message;
  },
};

const baseCodeInfoResponse: object = { codeId: 0, source: "", builder: "" };

export const CodeInfoResponse = {
  encode(message: CodeInfoResponse, writer: Writer = Writer.create()): Writer {
    if (message.codeId !== 0) {
      writer.uint32(8).uint64(message.codeId);
    }
    if (message.creator.length !== 0) {
      writer.uint32(18).bytes(message.creator);
    }
    if (message.codeHash.length !== 0) {
      writer.uint32(26).bytes(message.codeHash);
    }
    if (message.source !== "") {
      writer.uint32(34).string(message.source);
    }
    if (message.builder !== "") {
      writer.uint32(42).string(message.builder);
    }
    if (message.contractDuration !== undefined) {
      Duration.encode(
        message.contractDuration,
        writer.uint32(50).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CodeInfoResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCodeInfoResponse } as CodeInfoResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeId = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.creator = reader.bytes();
          break;
        case 3:
          message.codeHash = reader.bytes();
          break;
        case 4:
          message.source = reader.string();
          break;
        case 5:
          message.builder = reader.string();
          break;
        case 6:
          message.contractDuration = Duration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CodeInfoResponse {
    const message = { ...baseCodeInfoResponse } as CodeInfoResponse;
    if (object.codeId !== undefined && object.codeId !== null) {
      message.codeId = Number(object.codeId);
    } else {
      message.codeId = 0;
    }
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = bytesFromBase64(object.creator);
    }
    if (object.codeHash !== undefined && object.codeHash !== null) {
      message.codeHash = bytesFromBase64(object.codeHash);
    }
    if (object.source !== undefined && object.source !== null) {
      message.source = String(object.source);
    } else {
      message.source = "";
    }
    if (object.builder !== undefined && object.builder !== null) {
      message.builder = String(object.builder);
    } else {
      message.builder = "";
    }
    if (
      object.contractDuration !== undefined &&
      object.contractDuration !== null
    ) {
      message.contractDuration = Duration.fromJSON(object.contractDuration);
    } else {
      message.contractDuration = undefined;
    }
    return message;
  },

  toJSON(message: CodeInfoResponse): unknown {
    const obj: any = {};
    message.codeId !== undefined && (obj.codeId = message.codeId);
    message.creator !== undefined &&
      (obj.creator = base64FromBytes(
        message.creator !== undefined ? message.creator : new Uint8Array()
      ));
    message.codeHash !== undefined &&
      (obj.codeHash = base64FromBytes(
        message.codeHash !== undefined ? message.codeHash : new Uint8Array()
      ));
    message.source !== undefined && (obj.source = message.source);
    message.builder !== undefined && (obj.builder = message.builder);
    message.contractDuration !== undefined &&
      (obj.contractDuration = message.contractDuration
        ? Duration.toJSON(message.contractDuration)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<CodeInfoResponse>): CodeInfoResponse {
    const message = { ...baseCodeInfoResponse } as CodeInfoResponse;
    if (object.codeId !== undefined && object.codeId !== null) {
      message.codeId = object.codeId;
    } else {
      message.codeId = 0;
    }
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = new Uint8Array();
    }
    if (object.codeHash !== undefined && object.codeHash !== null) {
      message.codeHash = object.codeHash;
    } else {
      message.codeHash = new Uint8Array();
    }
    if (object.source !== undefined && object.source !== null) {
      message.source = object.source;
    } else {
      message.source = "";
    }
    if (object.builder !== undefined && object.builder !== null) {
      message.builder = object.builder;
    } else {
      message.builder = "";
    }
    if (
      object.contractDuration !== undefined &&
      object.contractDuration !== null
    ) {
      message.contractDuration = Duration.fromPartial(object.contractDuration);
    } else {
      message.contractDuration = undefined;
    }
    return message;
  },
};

const baseQueryCodeResponse: object = {};

export const QueryCodeResponse = {
  encode(message: QueryCodeResponse, writer: Writer = Writer.create()): Writer {
    if (message.codeInfo !== undefined) {
      CodeInfoResponse.encode(
        message.codeInfo,
        writer.uint32(10).fork()
      ).ldelim();
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryCodeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryCodeResponse } as QueryCodeResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeInfo = CodeInfoResponse.decode(reader, reader.uint32());
          break;
        case 2:
          message.data = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryCodeResponse {
    const message = { ...baseQueryCodeResponse } as QueryCodeResponse;
    if (object.codeInfo !== undefined && object.codeInfo !== null) {
      message.codeInfo = CodeInfoResponse.fromJSON(object.codeInfo);
    } else {
      message.codeInfo = undefined;
    }
    if (object.data !== undefined && object.data !== null) {
      message.data = bytesFromBase64(object.data);
    }
    return message;
  },

  toJSON(message: QueryCodeResponse): unknown {
    const obj: any = {};
    message.codeInfo !== undefined &&
      (obj.codeInfo = message.codeInfo
        ? CodeInfoResponse.toJSON(message.codeInfo)
        : undefined);
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<QueryCodeResponse>): QueryCodeResponse {
    const message = { ...baseQueryCodeResponse } as QueryCodeResponse;
    if (object.codeInfo !== undefined && object.codeInfo !== null) {
      message.codeInfo = CodeInfoResponse.fromPartial(object.codeInfo);
    } else {
      message.codeInfo = undefined;
    }
    if (object.data !== undefined && object.data !== null) {
      message.data = object.data;
    } else {
      message.data = new Uint8Array();
    }
    return message;
  },
};

const baseQueryCodesResponse: object = {};

export const QueryCodesResponse = {
  encode(
    message: QueryCodesResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.codeInfos) {
      CodeInfoResponse.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryCodesResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryCodesResponse } as QueryCodesResponse;
    message.codeInfos = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeInfos.push(
            CodeInfoResponse.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryCodesResponse {
    const message = { ...baseQueryCodesResponse } as QueryCodesResponse;
    message.codeInfos = [];
    if (object.codeInfos !== undefined && object.codeInfos !== null) {
      for (const e of object.codeInfos) {
        message.codeInfos.push(CodeInfoResponse.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: QueryCodesResponse): unknown {
    const obj: any = {};
    if (message.codeInfos) {
      obj.codeInfos = message.codeInfos.map((e) =>
        e ? CodeInfoResponse.toJSON(e) : undefined
      );
    } else {
      obj.codeInfos = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<QueryCodesResponse>): QueryCodesResponse {
    const message = { ...baseQueryCodesResponse } as QueryCodesResponse;
    message.codeInfos = [];
    if (object.codeInfos !== undefined && object.codeInfos !== null) {
      for (const e of object.codeInfos) {
        message.codeInfos.push(CodeInfoResponse.fromPartial(e));
      }
    }
    return message;
  },
};

const baseQueryContractAddressByContractIdResponse: object = {};

export const QueryContractAddressByContractIdResponse = {
  encode(
    message: QueryContractAddressByContractIdResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.address.length !== 0) {
      writer.uint32(10).bytes(message.address);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractAddressByContractIdResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractAddressByContractIdResponse,
    } as QueryContractAddressByContractIdResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractAddressByContractIdResponse {
    const message = {
      ...baseQueryContractAddressByContractIdResponse,
    } as QueryContractAddressByContractIdResponse;
    if (object.address !== undefined && object.address !== null) {
      message.address = bytesFromBase64(object.address);
    }
    return message;
  },

  toJSON(message: QueryContractAddressByContractIdResponse): unknown {
    const obj: any = {};
    message.address !== undefined &&
      (obj.address = base64FromBytes(
        message.address !== undefined ? message.address : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractAddressByContractIdResponse>
  ): QueryContractAddressByContractIdResponse {
    const message = {
      ...baseQueryContractAddressByContractIdResponse,
    } as QueryContractAddressByContractIdResponse;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = new Uint8Array();
    }
    return message;
  },
};

const baseQueryContractKeyResponse: object = {};

export const QueryContractKeyResponse = {
  encode(
    message: QueryContractKeyResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.key.length !== 0) {
      writer.uint32(10).bytes(message.key);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractKeyResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractKeyResponse,
    } as QueryContractKeyResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractKeyResponse {
    const message = {
      ...baseQueryContractKeyResponse,
    } as QueryContractKeyResponse;
    if (object.key !== undefined && object.key !== null) {
      message.key = bytesFromBase64(object.key);
    }
    return message;
  },

  toJSON(message: QueryContractKeyResponse): unknown {
    const obj: any = {};
    message.key !== undefined &&
      (obj.key = base64FromBytes(
        message.key !== undefined ? message.key : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractKeyResponse>
  ): QueryContractKeyResponse {
    const message = {
      ...baseQueryContractKeyResponse,
    } as QueryContractKeyResponse;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = new Uint8Array();
    }
    return message;
  },
};

const baseQueryContractHashResponse: object = {};

export const QueryContractHashResponse = {
  encode(
    message: QueryContractHashResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.codeHash.length !== 0) {
      writer.uint32(10).bytes(message.codeHash);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryContractHashResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryContractHashResponse,
    } as QueryContractHashResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeHash = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryContractHashResponse {
    const message = {
      ...baseQueryContractHashResponse,
    } as QueryContractHashResponse;
    if (object.codeHash !== undefined && object.codeHash !== null) {
      message.codeHash = bytesFromBase64(object.codeHash);
    }
    return message;
  },

  toJSON(message: QueryContractHashResponse): unknown {
    const obj: any = {};
    message.codeHash !== undefined &&
      (obj.codeHash = base64FromBytes(
        message.codeHash !== undefined ? message.codeHash : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryContractHashResponse>
  ): QueryContractHashResponse {
    const message = {
      ...baseQueryContractHashResponse,
    } as QueryContractHashResponse;
    if (object.codeHash !== undefined && object.codeHash !== null) {
      message.codeHash = object.codeHash;
    } else {
      message.codeHash = new Uint8Array();
    }
    return message;
  },
};

const baseDecryptedAnswer: object = {
  type: "",
  input: "",
  outputData: "",
  outputDataAsString: "",
  plaintextError: "",
};

export const DecryptedAnswer = {
  encode(message: DecryptedAnswer, writer: Writer = Writer.create()): Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.input !== "") {
      writer.uint32(18).string(message.input);
    }
    if (message.outputData !== "") {
      writer.uint32(26).string(message.outputData);
    }
    if (message.outputDataAsString !== "") {
      writer.uint32(34).string(message.outputDataAsString);
    }
    for (const v of message.outputLogs) {
      StringEvent.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.outputError.length !== 0) {
      writer.uint32(50).bytes(message.outputError);
    }
    if (message.plaintextError !== "") {
      writer.uint32(58).string(message.plaintextError);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): DecryptedAnswer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDecryptedAnswer } as DecryptedAnswer;
    message.outputLogs = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.string();
          break;
        case 2:
          message.input = reader.string();
          break;
        case 3:
          message.outputData = reader.string();
          break;
        case 4:
          message.outputDataAsString = reader.string();
          break;
        case 5:
          message.outputLogs.push(StringEvent.decode(reader, reader.uint32()));
          break;
        case 6:
          message.outputError = reader.bytes();
          break;
        case 7:
          message.plaintextError = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DecryptedAnswer {
    const message = { ...baseDecryptedAnswer } as DecryptedAnswer;
    message.outputLogs = [];
    if (object.type !== undefined && object.type !== null) {
      message.type = String(object.type);
    } else {
      message.type = "";
    }
    if (object.input !== undefined && object.input !== null) {
      message.input = String(object.input);
    } else {
      message.input = "";
    }
    if (object.outputData !== undefined && object.outputData !== null) {
      message.outputData = String(object.outputData);
    } else {
      message.outputData = "";
    }
    if (
      object.outputDataAsString !== undefined &&
      object.outputDataAsString !== null
    ) {
      message.outputDataAsString = String(object.outputDataAsString);
    } else {
      message.outputDataAsString = "";
    }
    if (object.outputLogs !== undefined && object.outputLogs !== null) {
      for (const e of object.outputLogs) {
        message.outputLogs.push(StringEvent.fromJSON(e));
      }
    }
    if (object.outputError !== undefined && object.outputError !== null) {
      message.outputError = bytesFromBase64(object.outputError);
    }
    if (object.plaintextError !== undefined && object.plaintextError !== null) {
      message.plaintextError = String(object.plaintextError);
    } else {
      message.plaintextError = "";
    }
    return message;
  },

  toJSON(message: DecryptedAnswer): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.input !== undefined && (obj.input = message.input);
    message.outputData !== undefined && (obj.outputData = message.outputData);
    message.outputDataAsString !== undefined &&
      (obj.outputDataAsString = message.outputDataAsString);
    if (message.outputLogs) {
      obj.outputLogs = message.outputLogs.map((e) =>
        e ? StringEvent.toJSON(e) : undefined
      );
    } else {
      obj.outputLogs = [];
    }
    message.outputError !== undefined &&
      (obj.outputError = base64FromBytes(
        message.outputError !== undefined
          ? message.outputError
          : new Uint8Array()
      ));
    message.plaintextError !== undefined &&
      (obj.plaintextError = message.plaintextError);
    return obj;
  },

  fromPartial(object: DeepPartial<DecryptedAnswer>): DecryptedAnswer {
    const message = { ...baseDecryptedAnswer } as DecryptedAnswer;
    message.outputLogs = [];
    if (object.type !== undefined && object.type !== null) {
      message.type = object.type;
    } else {
      message.type = "";
    }
    if (object.input !== undefined && object.input !== null) {
      message.input = object.input;
    } else {
      message.input = "";
    }
    if (object.outputData !== undefined && object.outputData !== null) {
      message.outputData = object.outputData;
    } else {
      message.outputData = "";
    }
    if (
      object.outputDataAsString !== undefined &&
      object.outputDataAsString !== null
    ) {
      message.outputDataAsString = object.outputDataAsString;
    } else {
      message.outputDataAsString = "";
    }
    if (object.outputLogs !== undefined && object.outputLogs !== null) {
      for (const e of object.outputLogs) {
        message.outputLogs.push(StringEvent.fromPartial(e));
      }
    }
    if (object.outputError !== undefined && object.outputError !== null) {
      message.outputError = object.outputError;
    } else {
      message.outputError = new Uint8Array();
    }
    if (object.plaintextError !== undefined && object.plaintextError !== null) {
      message.plaintextError = object.plaintextError;
    } else {
      message.plaintextError = "";
    }
    return message;
  },
};

/** Query provides defines the gRPC querier service */
export interface Query {
  ContractInfo(
    request: QueryContractInfoRequest
  ): Promise<QueryContractInfoResponse>;
  ContractResult(
    request: QueryContractResultRequest
  ): Promise<QueryContractResultResponse>;
  /**
   * rpc ContractHistory (QueryContractHistoryRequest) returns (QueryContractHistoryResponse) {
   * option (google.api.http).get = "/compute/v1beta1/contract/{address}/history";
   * }
   */
  ContractsByCode(
    request: QueryContractsByCodeRequest
  ): Promise<QueryContractsByCodeResponse>;
  /**
   * rpc AllContractState (QueryAllContractStateRequest) returns (QueryAllContractStateResponse) {
   * option (google.api.http).get = "/compute/v1beta1/contract/{address}/state";
   * }
   * rpc RawContractState (QueryRawContractStateRequest) returns (QueryRawContractStateResponse) {
   * option (google.api.http).get = "/compute/v1beta1/contract/{address}/raw/{query_data}";
   * }
   */
  SmartContractState(
    request: QuerySmartContractStateRequest
  ): Promise<QuerySmartContractStateResponse>;
  Code(request: QueryCodeRequest): Promise<QueryCodeResponse>;
  Codes(request: Empty): Promise<QueryCodesResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  ContractInfo(
    request: QueryContractInfoRequest
  ): Promise<QueryContractInfoResponse> {
    const data = QueryContractInfoRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.compute.v1beta1.Query",
      "ContractInfo",
      data
    );
    return promise.then((data) =>
      QueryContractInfoResponse.decode(new Reader(data))
    );
  }

  ContractResult(
    request: QueryContractResultRequest
  ): Promise<QueryContractResultResponse> {
    const data = QueryContractResultRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.compute.v1beta1.Query",
      "ContractResult",
      data
    );
    return promise.then((data) =>
      QueryContractResultResponse.decode(new Reader(data))
    );
  }

  ContractsByCode(
    request: QueryContractsByCodeRequest
  ): Promise<QueryContractsByCodeResponse> {
    const data = QueryContractsByCodeRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.compute.v1beta1.Query",
      "ContractsByCode",
      data
    );
    return promise.then((data) =>
      QueryContractsByCodeResponse.decode(new Reader(data))
    );
  }

  SmartContractState(
    request: QuerySmartContractStateRequest
  ): Promise<QuerySmartContractStateResponse> {
    const data = QuerySmartContractStateRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.compute.v1beta1.Query",
      "SmartContractState",
      data
    );
    return promise.then((data) =>
      QuerySmartContractStateResponse.decode(new Reader(data))
    );
  }

  Code(request: QueryCodeRequest): Promise<QueryCodeResponse> {
    const data = QueryCodeRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.compute.v1beta1.Query",
      "Code",
      data
    );
    return promise.then((data) => QueryCodeResponse.decode(new Reader(data)));
  }

  Codes(request: Empty): Promise<QueryCodesResponse> {
    const data = Empty.encode(request).finish();
    const promise = this.rpc.request(
      "trst.x.compute.v1beta1.Query",
      "Codes",
      data
    );
    return promise.then((data) => QueryCodesResponse.decode(new Reader(data)));
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

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
