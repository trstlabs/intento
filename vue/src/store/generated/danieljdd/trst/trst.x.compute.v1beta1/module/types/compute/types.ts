/* eslint-disable */
import { Timestamp } from "../google/protobuf/timestamp";
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";

export const protobufPackage = "trst.x.compute.v1beta1";

export enum AccessType {
  UNDEFINED = 0,
  NOBODY = 1,
  ONLY_ADDRESS = 2,
  EVERYBODY = 3,
  UNRECOGNIZED = -1,
}

export function accessTypeFromJSON(object: any): AccessType {
  switch (object) {
    case 0:
    case "UNDEFINED":
      return AccessType.UNDEFINED;
    case 1:
    case "NOBODY":
      return AccessType.NOBODY;
    case 2:
    case "ONLY_ADDRESS":
      return AccessType.ONLY_ADDRESS;
    case 3:
    case "EVERYBODY":
      return AccessType.EVERYBODY;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AccessType.UNRECOGNIZED;
  }
}

export function accessTypeToJSON(object: AccessType): string {
  switch (object) {
    case AccessType.UNDEFINED:
      return "UNDEFINED";
    case AccessType.NOBODY:
      return "NOBODY";
    case AccessType.ONLY_ADDRESS:
      return "ONLY_ADDRESS";
    case AccessType.EVERYBODY:
      return "EVERYBODY";
    default:
      return "UNKNOWN";
  }
}

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

const baseAccessTypeParam: object = { value: 0 };

export const AccessTypeParam = {
  encode(message: AccessTypeParam, writer: Writer = Writer.create()): Writer {
    if (message.value !== 0) {
      writer.uint32(8).int32(message.value);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AccessTypeParam {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAccessTypeParam } as AccessTypeParam;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.value = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AccessTypeParam {
    const message = { ...baseAccessTypeParam } as AccessTypeParam;
    if (object.value !== undefined && object.value !== null) {
      message.value = accessTypeFromJSON(object.value);
    } else {
      message.value = 0;
    }
    return message;
  },

  toJSON(message: AccessTypeParam): unknown {
    const obj: any = {};
    message.value !== undefined &&
      (obj.value = accessTypeToJSON(message.value));
    return obj;
  },

  fromPartial(object: DeepPartial<AccessTypeParam>): AccessTypeParam {
    const message = { ...baseAccessTypeParam } as AccessTypeParam;
    if (object.value !== undefined && object.value !== null) {
      message.value = object.value;
    } else {
      message.value = 0;
    }
    return message;
  },
};

const baseCodeInfo: object = { source: "", builder: "" };

export const CodeInfo = {
  encode(message: CodeInfo, writer: Writer = Writer.create()): Writer {
    if (message.codeHash.length !== 0) {
      writer.uint32(10).bytes(message.codeHash);
    }
    if (message.creator.length !== 0) {
      writer.uint32(18).bytes(message.creator);
    }
    if (message.source !== "") {
      writer.uint32(26).string(message.source);
    }
    if (message.builder !== "") {
      writer.uint32(34).string(message.builder);
    }
    if (message.endTime !== undefined) {
      Duration.encode(message.endTime, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CodeInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCodeInfo } as CodeInfo;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeHash = reader.bytes();
          break;
        case 2:
          message.creator = reader.bytes();
          break;
        case 3:
          message.source = reader.string();
          break;
        case 4:
          message.builder = reader.string();
          break;
        case 5:
          message.endTime = Duration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CodeInfo {
    const message = { ...baseCodeInfo } as CodeInfo;
    if (object.codeHash !== undefined && object.codeHash !== null) {
      message.codeHash = bytesFromBase64(object.codeHash);
    }
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = bytesFromBase64(object.creator);
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
    if (object.endTime !== undefined && object.endTime !== null) {
      message.endTime = Duration.fromJSON(object.endTime);
    } else {
      message.endTime = undefined;
    }
    return message;
  },

  toJSON(message: CodeInfo): unknown {
    const obj: any = {};
    message.codeHash !== undefined &&
      (obj.codeHash = base64FromBytes(
        message.codeHash !== undefined ? message.codeHash : new Uint8Array()
      ));
    message.creator !== undefined &&
      (obj.creator = base64FromBytes(
        message.creator !== undefined ? message.creator : new Uint8Array()
      ));
    message.source !== undefined && (obj.source = message.source);
    message.builder !== undefined && (obj.builder = message.builder);
    message.endTime !== undefined &&
      (obj.endTime = message.endTime
        ? Duration.toJSON(message.endTime)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<CodeInfo>): CodeInfo {
    const message = { ...baseCodeInfo } as CodeInfo;
    if (object.codeHash !== undefined && object.codeHash !== null) {
      message.codeHash = object.codeHash;
    } else {
      message.codeHash = new Uint8Array();
    }
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = new Uint8Array();
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
    if (object.endTime !== undefined && object.endTime !== null) {
      message.endTime = Duration.fromPartial(object.endTime);
    } else {
      message.endTime = undefined;
    }
    return message;
  },
};

const baseContractInfo: object = { codeId: 0, contractId: "" };

export const ContractInfo = {
  encode(message: ContractInfo, writer: Writer = Writer.create()): Writer {
    if (message.codeId !== 0) {
      writer.uint32(8).uint64(message.codeId);
    }
    if (message.creator.length !== 0) {
      writer.uint32(18).bytes(message.creator);
    }
    if (message.contractId !== "") {
      writer.uint32(34).string(message.contractId);
    }
    if (message.created !== undefined) {
      AbsoluteTxPosition.encode(
        message.created,
        writer.uint32(42).fork()
      ).ldelim();
    }
    if (message.endTime !== undefined) {
      Timestamp.encode(
        toTimestamp(message.endTime),
        writer.uint32(50).fork()
      ).ldelim();
    }
    if (message.lastMsg.length !== 0) {
      writer.uint32(58).bytes(message.lastMsg);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ContractInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContractInfo } as ContractInfo;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeId = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.creator = reader.bytes();
          break;
        case 4:
          message.contractId = reader.string();
          break;
        case 5:
          message.created = AbsoluteTxPosition.decode(reader, reader.uint32());
          break;
        case 6:
          message.endTime = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 7:
          message.lastMsg = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractInfo {
    const message = { ...baseContractInfo } as ContractInfo;
    if (object.codeId !== undefined && object.codeId !== null) {
      message.codeId = Number(object.codeId);
    } else {
      message.codeId = 0;
    }
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = bytesFromBase64(object.creator);
    }
    if (object.contractId !== undefined && object.contractId !== null) {
      message.contractId = String(object.contractId);
    } else {
      message.contractId = "";
    }
    if (object.created !== undefined && object.created !== null) {
      message.created = AbsoluteTxPosition.fromJSON(object.created);
    } else {
      message.created = undefined;
    }
    if (object.endTime !== undefined && object.endTime !== null) {
      message.endTime = fromJsonTimestamp(object.endTime);
    } else {
      message.endTime = undefined;
    }
    if (object.lastMsg !== undefined && object.lastMsg !== null) {
      message.lastMsg = bytesFromBase64(object.lastMsg);
    }
    return message;
  },

  toJSON(message: ContractInfo): unknown {
    const obj: any = {};
    message.codeId !== undefined && (obj.codeId = message.codeId);
    message.creator !== undefined &&
      (obj.creator = base64FromBytes(
        message.creator !== undefined ? message.creator : new Uint8Array()
      ));
    message.contractId !== undefined && (obj.contractId = message.contractId);
    message.created !== undefined &&
      (obj.created = message.created
        ? AbsoluteTxPosition.toJSON(message.created)
        : undefined);
    message.endTime !== undefined &&
      (obj.endTime =
        message.endTime !== undefined ? message.endTime.toISOString() : null);
    message.lastMsg !== undefined &&
      (obj.lastMsg = base64FromBytes(
        message.lastMsg !== undefined ? message.lastMsg : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<ContractInfo>): ContractInfo {
    const message = { ...baseContractInfo } as ContractInfo;
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
    if (object.contractId !== undefined && object.contractId !== null) {
      message.contractId = object.contractId;
    } else {
      message.contractId = "";
    }
    if (object.created !== undefined && object.created !== null) {
      message.created = AbsoluteTxPosition.fromPartial(object.created);
    } else {
      message.created = undefined;
    }
    if (object.endTime !== undefined && object.endTime !== null) {
      message.endTime = object.endTime;
    } else {
      message.endTime = undefined;
    }
    if (object.lastMsg !== undefined && object.lastMsg !== null) {
      message.lastMsg = object.lastMsg;
    } else {
      message.lastMsg = new Uint8Array();
    }
    return message;
  },
};

const baseContractInfoWithAddress: object = {};

export const ContractInfoWithAddress = {
  encode(
    message: ContractInfoWithAddress,
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

  decode(input: Reader | Uint8Array, length?: number): ContractInfoWithAddress {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseContractInfoWithAddress,
    } as ContractInfoWithAddress;
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

  fromJSON(object: any): ContractInfoWithAddress {
    const message = {
      ...baseContractInfoWithAddress,
    } as ContractInfoWithAddress;
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

  toJSON(message: ContractInfoWithAddress): unknown {
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
    object: DeepPartial<ContractInfoWithAddress>
  ): ContractInfoWithAddress {
    const message = {
      ...baseContractInfoWithAddress,
    } as ContractInfoWithAddress;
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

const baseAbsoluteTxPosition: object = { blockHeight: 0, txIndex: 0 };

export const AbsoluteTxPosition = {
  encode(
    message: AbsoluteTxPosition,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.blockHeight !== 0) {
      writer.uint32(8).int64(message.blockHeight);
    }
    if (message.txIndex !== 0) {
      writer.uint32(16).uint64(message.txIndex);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AbsoluteTxPosition {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAbsoluteTxPosition } as AbsoluteTxPosition;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockHeight = longToNumber(reader.int64() as Long);
          break;
        case 2:
          message.txIndex = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AbsoluteTxPosition {
    const message = { ...baseAbsoluteTxPosition } as AbsoluteTxPosition;
    if (object.blockHeight !== undefined && object.blockHeight !== null) {
      message.blockHeight = Number(object.blockHeight);
    } else {
      message.blockHeight = 0;
    }
    if (object.txIndex !== undefined && object.txIndex !== null) {
      message.txIndex = Number(object.txIndex);
    } else {
      message.txIndex = 0;
    }
    return message;
  },

  toJSON(message: AbsoluteTxPosition): unknown {
    const obj: any = {};
    message.blockHeight !== undefined &&
      (obj.blockHeight = message.blockHeight);
    message.txIndex !== undefined && (obj.txIndex = message.txIndex);
    return obj;
  },

  fromPartial(object: DeepPartial<AbsoluteTxPosition>): AbsoluteTxPosition {
    const message = { ...baseAbsoluteTxPosition } as AbsoluteTxPosition;
    if (object.blockHeight !== undefined && object.blockHeight !== null) {
      message.blockHeight = object.blockHeight;
    } else {
      message.blockHeight = 0;
    }
    if (object.txIndex !== undefined && object.txIndex !== null) {
      message.txIndex = object.txIndex;
    } else {
      message.txIndex = 0;
    }
    return message;
  },
};

const baseModel: object = {};

export const Model = {
  encode(message: Model, writer: Writer = Writer.create()): Writer {
    if (message.Key.length !== 0) {
      writer.uint32(10).bytes(message.Key);
    }
    if (message.Value.length !== 0) {
      writer.uint32(18).bytes(message.Value);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Model {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseModel } as Model;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Key = reader.bytes();
          break;
        case 2:
          message.Value = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Model {
    const message = { ...baseModel } as Model;
    if (object.Key !== undefined && object.Key !== null) {
      message.Key = bytesFromBase64(object.Key);
    }
    if (object.Value !== undefined && object.Value !== null) {
      message.Value = bytesFromBase64(object.Value);
    }
    return message;
  },

  toJSON(message: Model): unknown {
    const obj: any = {};
    message.Key !== undefined &&
      (obj.Key = base64FromBytes(
        message.Key !== undefined ? message.Key : new Uint8Array()
      ));
    message.Value !== undefined &&
      (obj.Value = base64FromBytes(
        message.Value !== undefined ? message.Value : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<Model>): Model {
    const message = { ...baseModel } as Model;
    if (object.Key !== undefined && object.Key !== null) {
      message.Key = object.Key;
    } else {
      message.Key = new Uint8Array();
    }
    if (object.Value !== undefined && object.Value !== null) {
      message.Value = object.Value;
    } else {
      message.Value = new Uint8Array();
    }
    return message;
  },
};

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

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

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
