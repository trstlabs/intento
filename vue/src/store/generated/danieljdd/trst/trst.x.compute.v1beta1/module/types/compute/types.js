/* eslint-disable */
import { Timestamp } from "../google/protobuf/timestamp";
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
export const protobufPackage = "trst.x.compute.v1beta1";
export var AccessType;
(function (AccessType) {
    AccessType[AccessType["UNDEFINED"] = 0] = "UNDEFINED";
    AccessType[AccessType["NOBODY"] = 1] = "NOBODY";
    AccessType[AccessType["ONLY_ADDRESS"] = 2] = "ONLY_ADDRESS";
    AccessType[AccessType["EVERYBODY"] = 3] = "EVERYBODY";
    AccessType[AccessType["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(AccessType || (AccessType = {}));
export function accessTypeFromJSON(object) {
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
export function accessTypeToJSON(object) {
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
const baseAccessTypeParam = { value: 0 };
export const AccessTypeParam = {
    encode(message, writer = Writer.create()) {
        if (message.value !== 0) {
            writer.uint32(8).int32(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseAccessTypeParam };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.value = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseAccessTypeParam };
        if (object.value !== undefined && object.value !== null) {
            message.value = accessTypeFromJSON(object.value);
        }
        else {
            message.value = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.value !== undefined &&
            (obj.value = accessTypeToJSON(message.value));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseAccessTypeParam };
        if (object.value !== undefined && object.value !== null) {
            message.value = object.value;
        }
        else {
            message.value = 0;
        }
        return message;
    },
};
const baseCodeInfo = { source: "", builder: "" };
export const CodeInfo = {
    encode(message, writer = Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseCodeInfo };
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
    fromJSON(object) {
        const message = { ...baseCodeInfo };
        if (object.codeHash !== undefined && object.codeHash !== null) {
            message.codeHash = bytesFromBase64(object.codeHash);
        }
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = bytesFromBase64(object.creator);
        }
        if (object.source !== undefined && object.source !== null) {
            message.source = String(object.source);
        }
        else {
            message.source = "";
        }
        if (object.builder !== undefined && object.builder !== null) {
            message.builder = String(object.builder);
        }
        else {
            message.builder = "";
        }
        if (object.endTime !== undefined && object.endTime !== null) {
            message.endTime = Duration.fromJSON(object.endTime);
        }
        else {
            message.endTime = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.codeHash !== undefined &&
            (obj.codeHash = base64FromBytes(message.codeHash !== undefined ? message.codeHash : new Uint8Array()));
        message.creator !== undefined &&
            (obj.creator = base64FromBytes(message.creator !== undefined ? message.creator : new Uint8Array()));
        message.source !== undefined && (obj.source = message.source);
        message.builder !== undefined && (obj.builder = message.builder);
        message.endTime !== undefined &&
            (obj.endTime = message.endTime
                ? Duration.toJSON(message.endTime)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseCodeInfo };
        if (object.codeHash !== undefined && object.codeHash !== null) {
            message.codeHash = object.codeHash;
        }
        else {
            message.codeHash = new Uint8Array();
        }
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = object.creator;
        }
        else {
            message.creator = new Uint8Array();
        }
        if (object.source !== undefined && object.source !== null) {
            message.source = object.source;
        }
        else {
            message.source = "";
        }
        if (object.builder !== undefined && object.builder !== null) {
            message.builder = object.builder;
        }
        else {
            message.builder = "";
        }
        if (object.endTime !== undefined && object.endTime !== null) {
            message.endTime = Duration.fromPartial(object.endTime);
        }
        else {
            message.endTime = undefined;
        }
        return message;
    },
};
const baseContractInfo = { codeId: 0, contractId: "" };
export const ContractInfo = {
    encode(message, writer = Writer.create()) {
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
            AbsoluteTxPosition.encode(message.created, writer.uint32(42).fork()).ldelim();
        }
        if (message.endTime !== undefined) {
            Timestamp.encode(toTimestamp(message.endTime), writer.uint32(50).fork()).ldelim();
        }
        if (message.lastMsg.length !== 0) {
            writer.uint32(58).bytes(message.lastMsg);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseContractInfo };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeId = longToNumber(reader.uint64());
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
                    message.endTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
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
    fromJSON(object) {
        const message = { ...baseContractInfo };
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = Number(object.codeId);
        }
        else {
            message.codeId = 0;
        }
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = bytesFromBase64(object.creator);
        }
        if (object.contractId !== undefined && object.contractId !== null) {
            message.contractId = String(object.contractId);
        }
        else {
            message.contractId = "";
        }
        if (object.created !== undefined && object.created !== null) {
            message.created = AbsoluteTxPosition.fromJSON(object.created);
        }
        else {
            message.created = undefined;
        }
        if (object.endTime !== undefined && object.endTime !== null) {
            message.endTime = fromJsonTimestamp(object.endTime);
        }
        else {
            message.endTime = undefined;
        }
        if (object.lastMsg !== undefined && object.lastMsg !== null) {
            message.lastMsg = bytesFromBase64(object.lastMsg);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.codeId !== undefined && (obj.codeId = message.codeId);
        message.creator !== undefined &&
            (obj.creator = base64FromBytes(message.creator !== undefined ? message.creator : new Uint8Array()));
        message.contractId !== undefined && (obj.contractId = message.contractId);
        message.created !== undefined &&
            (obj.created = message.created
                ? AbsoluteTxPosition.toJSON(message.created)
                : undefined);
        message.endTime !== undefined &&
            (obj.endTime =
                message.endTime !== undefined ? message.endTime.toISOString() : null);
        message.lastMsg !== undefined &&
            (obj.lastMsg = base64FromBytes(message.lastMsg !== undefined ? message.lastMsg : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseContractInfo };
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = object.codeId;
        }
        else {
            message.codeId = 0;
        }
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = object.creator;
        }
        else {
            message.creator = new Uint8Array();
        }
        if (object.contractId !== undefined && object.contractId !== null) {
            message.contractId = object.contractId;
        }
        else {
            message.contractId = "";
        }
        if (object.created !== undefined && object.created !== null) {
            message.created = AbsoluteTxPosition.fromPartial(object.created);
        }
        else {
            message.created = undefined;
        }
        if (object.endTime !== undefined && object.endTime !== null) {
            message.endTime = object.endTime;
        }
        else {
            message.endTime = undefined;
        }
        if (object.lastMsg !== undefined && object.lastMsg !== null) {
            message.lastMsg = object.lastMsg;
        }
        else {
            message.lastMsg = new Uint8Array();
        }
        return message;
    },
};
const baseContractInfoWithAddress = {};
export const ContractInfoWithAddress = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        if (message.ContractInfo !== undefined) {
            ContractInfo.encode(message.ContractInfo, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseContractInfoWithAddress,
        };
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
    fromJSON(object) {
        const message = {
            ...baseContractInfoWithAddress,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        if (object.ContractInfo !== undefined && object.ContractInfo !== null) {
            message.ContractInfo = ContractInfo.fromJSON(object.ContractInfo);
        }
        else {
            message.ContractInfo = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        message.ContractInfo !== undefined &&
            (obj.ContractInfo = message.ContractInfo
                ? ContractInfo.toJSON(message.ContractInfo)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseContractInfoWithAddress,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        if (object.ContractInfo !== undefined && object.ContractInfo !== null) {
            message.ContractInfo = ContractInfo.fromPartial(object.ContractInfo);
        }
        else {
            message.ContractInfo = undefined;
        }
        return message;
    },
};
const baseAbsoluteTxPosition = { blockHeight: 0, txIndex: 0 };
export const AbsoluteTxPosition = {
    encode(message, writer = Writer.create()) {
        if (message.blockHeight !== 0) {
            writer.uint32(8).int64(message.blockHeight);
        }
        if (message.txIndex !== 0) {
            writer.uint32(16).uint64(message.txIndex);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseAbsoluteTxPosition };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.blockHeight = longToNumber(reader.int64());
                    break;
                case 2:
                    message.txIndex = longToNumber(reader.uint64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseAbsoluteTxPosition };
        if (object.blockHeight !== undefined && object.blockHeight !== null) {
            message.blockHeight = Number(object.blockHeight);
        }
        else {
            message.blockHeight = 0;
        }
        if (object.txIndex !== undefined && object.txIndex !== null) {
            message.txIndex = Number(object.txIndex);
        }
        else {
            message.txIndex = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.blockHeight !== undefined &&
            (obj.blockHeight = message.blockHeight);
        message.txIndex !== undefined && (obj.txIndex = message.txIndex);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseAbsoluteTxPosition };
        if (object.blockHeight !== undefined && object.blockHeight !== null) {
            message.blockHeight = object.blockHeight;
        }
        else {
            message.blockHeight = 0;
        }
        if (object.txIndex !== undefined && object.txIndex !== null) {
            message.txIndex = object.txIndex;
        }
        else {
            message.txIndex = 0;
        }
        return message;
    },
};
const baseModel = {};
export const Model = {
    encode(message, writer = Writer.create()) {
        if (message.Key.length !== 0) {
            writer.uint32(10).bytes(message.Key);
        }
        if (message.Value.length !== 0) {
            writer.uint32(18).bytes(message.Value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseModel };
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
    fromJSON(object) {
        const message = { ...baseModel };
        if (object.Key !== undefined && object.Key !== null) {
            message.Key = bytesFromBase64(object.Key);
        }
        if (object.Value !== undefined && object.Value !== null) {
            message.Value = bytesFromBase64(object.Value);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.Key !== undefined &&
            (obj.Key = base64FromBytes(message.Key !== undefined ? message.Key : new Uint8Array()));
        message.Value !== undefined &&
            (obj.Value = base64FromBytes(message.Value !== undefined ? message.Value : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseModel };
        if (object.Key !== undefined && object.Key !== null) {
            message.Key = object.Key;
        }
        else {
            message.Key = new Uint8Array();
        }
        if (object.Value !== undefined && object.Value !== null) {
            message.Value = object.Value;
        }
        else {
            message.Value = new Uint8Array();
        }
        return message;
    },
};
var globalThis = (() => {
    if (typeof globalThis !== "undefined")
        return globalThis;
    if (typeof self !== "undefined")
        return self;
    if (typeof window !== "undefined")
        return window;
    if (typeof global !== "undefined")
        return global;
    throw "Unable to locate global object";
})();
const atob = globalThis.atob ||
    ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64) {
    const bin = atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
        arr[i] = bin.charCodeAt(i);
    }
    return arr;
}
const btoa = globalThis.btoa ||
    ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr) {
    const bin = [];
    for (let i = 0; i < arr.byteLength; ++i) {
        bin.push(String.fromCharCode(arr[i]));
    }
    return btoa(bin.join(""));
}
function toTimestamp(date) {
    const seconds = date.getTime() / 1000;
    const nanos = (date.getTime() % 1000) * 1000000;
    return { seconds, nanos };
}
function fromTimestamp(t) {
    let millis = t.seconds * 1000;
    millis += t.nanos / 1000000;
    return new Date(millis);
}
function fromJsonTimestamp(o) {
    if (o instanceof Date) {
        return o;
    }
    else if (typeof o === "string") {
        return new Date(o);
    }
    else {
        return fromTimestamp(Timestamp.fromJSON(o));
    }
}
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
if (util.Long !== Long) {
    util.Long = Long;
    configure();
}
