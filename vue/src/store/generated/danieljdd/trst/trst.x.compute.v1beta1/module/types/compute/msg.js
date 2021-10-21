/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { Coin } from "../cosmos/base/v1beta1/coin";
export const protobufPackage = "trst.x.compute.v1beta1";
const baseMsgStoreCode = { source: "", builder: "", contractPeriod: 0 };
export const MsgStoreCode = {
    encode(message, writer = Writer.create()) {
        if (message.sender.length !== 0) {
            writer.uint32(10).bytes(message.sender);
        }
        if (message.wasmByteCode.length !== 0) {
            writer.uint32(18).bytes(message.wasmByteCode);
        }
        if (message.source !== "") {
            writer.uint32(26).string(message.source);
        }
        if (message.builder !== "") {
            writer.uint32(34).string(message.builder);
        }
        if (message.contractPeriod !== 0) {
            writer.uint32(48).int64(message.contractPeriod);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgStoreCode };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.bytes();
                    break;
                case 2:
                    message.wasmByteCode = reader.bytes();
                    break;
                case 3:
                    message.source = reader.string();
                    break;
                case 4:
                    message.builder = reader.string();
                    break;
                case 6:
                    message.contractPeriod = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgStoreCode };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = bytesFromBase64(object.sender);
        }
        if (object.wasmByteCode !== undefined && object.wasmByteCode !== null) {
            message.wasmByteCode = bytesFromBase64(object.wasmByteCode);
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
        if (object.contractPeriod !== undefined && object.contractPeriod !== null) {
            message.contractPeriod = Number(object.contractPeriod);
        }
        else {
            message.contractPeriod = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined &&
            (obj.sender = base64FromBytes(message.sender !== undefined ? message.sender : new Uint8Array()));
        message.wasmByteCode !== undefined &&
            (obj.wasmByteCode = base64FromBytes(message.wasmByteCode !== undefined
                ? message.wasmByteCode
                : new Uint8Array()));
        message.source !== undefined && (obj.source = message.source);
        message.builder !== undefined && (obj.builder = message.builder);
        message.contractPeriod !== undefined &&
            (obj.contractPeriod = message.contractPeriod);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgStoreCode };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = new Uint8Array();
        }
        if (object.wasmByteCode !== undefined && object.wasmByteCode !== null) {
            message.wasmByteCode = object.wasmByteCode;
        }
        else {
            message.wasmByteCode = new Uint8Array();
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
        if (object.contractPeriod !== undefined && object.contractPeriod !== null) {
            message.contractPeriod = object.contractPeriod;
        }
        else {
            message.contractPeriod = 0;
        }
        return message;
    },
};
const baseMsgInstantiateContract = {
    callbackCodeHash: "",
    codeId: 0,
    contractId: "",
};
export const MsgInstantiateContract = {
    encode(message, writer = Writer.create()) {
        if (message.sender.length !== 0) {
            writer.uint32(10).bytes(message.sender);
        }
        if (message.callbackCodeHash !== "") {
            writer.uint32(18).string(message.callbackCodeHash);
        }
        if (message.codeId !== 0) {
            writer.uint32(24).uint64(message.codeId);
        }
        if (message.contractId !== "") {
            writer.uint32(34).string(message.contractId);
        }
        if (message.initMsg.length !== 0) {
            writer.uint32(42).bytes(message.initMsg);
        }
        if (message.lastMsg.length !== 0) {
            writer.uint32(50).bytes(message.lastMsg);
        }
        for (const v of message.initFunds) {
            Coin.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (message.callbackSig.length !== 0) {
            writer.uint32(66).bytes(message.callbackSig);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgInstantiateContract };
        message.initFunds = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.bytes();
                    break;
                case 2:
                    message.callbackCodeHash = reader.string();
                    break;
                case 3:
                    message.codeId = longToNumber(reader.uint64());
                    break;
                case 4:
                    message.contractId = reader.string();
                    break;
                case 5:
                    message.initMsg = reader.bytes();
                    break;
                case 6:
                    message.lastMsg = reader.bytes();
                    break;
                case 7:
                    message.initFunds.push(Coin.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.callbackSig = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgInstantiateContract };
        message.initFunds = [];
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = bytesFromBase64(object.sender);
        }
        if (object.callbackCodeHash !== undefined &&
            object.callbackCodeHash !== null) {
            message.callbackCodeHash = String(object.callbackCodeHash);
        }
        else {
            message.callbackCodeHash = "";
        }
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = Number(object.codeId);
        }
        else {
            message.codeId = 0;
        }
        if (object.contractId !== undefined && object.contractId !== null) {
            message.contractId = String(object.contractId);
        }
        else {
            message.contractId = "";
        }
        if (object.initMsg !== undefined && object.initMsg !== null) {
            message.initMsg = bytesFromBase64(object.initMsg);
        }
        if (object.lastMsg !== undefined && object.lastMsg !== null) {
            message.lastMsg = bytesFromBase64(object.lastMsg);
        }
        if (object.initFunds !== undefined && object.initFunds !== null) {
            for (const e of object.initFunds) {
                message.initFunds.push(Coin.fromJSON(e));
            }
        }
        if (object.callbackSig !== undefined && object.callbackSig !== null) {
            message.callbackSig = bytesFromBase64(object.callbackSig);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined &&
            (obj.sender = base64FromBytes(message.sender !== undefined ? message.sender : new Uint8Array()));
        message.callbackCodeHash !== undefined &&
            (obj.callbackCodeHash = message.callbackCodeHash);
        message.codeId !== undefined && (obj.codeId = message.codeId);
        message.contractId !== undefined && (obj.contractId = message.contractId);
        message.initMsg !== undefined &&
            (obj.initMsg = base64FromBytes(message.initMsg !== undefined ? message.initMsg : new Uint8Array()));
        message.lastMsg !== undefined &&
            (obj.lastMsg = base64FromBytes(message.lastMsg !== undefined ? message.lastMsg : new Uint8Array()));
        if (message.initFunds) {
            obj.initFunds = message.initFunds.map((e) => e ? Coin.toJSON(e) : undefined);
        }
        else {
            obj.initFunds = [];
        }
        message.callbackSig !== undefined &&
            (obj.callbackSig = base64FromBytes(message.callbackSig !== undefined
                ? message.callbackSig
                : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgInstantiateContract };
        message.initFunds = [];
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = new Uint8Array();
        }
        if (object.callbackCodeHash !== undefined &&
            object.callbackCodeHash !== null) {
            message.callbackCodeHash = object.callbackCodeHash;
        }
        else {
            message.callbackCodeHash = "";
        }
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = object.codeId;
        }
        else {
            message.codeId = 0;
        }
        if (object.contractId !== undefined && object.contractId !== null) {
            message.contractId = object.contractId;
        }
        else {
            message.contractId = "";
        }
        if (object.initMsg !== undefined && object.initMsg !== null) {
            message.initMsg = object.initMsg;
        }
        else {
            message.initMsg = new Uint8Array();
        }
        if (object.lastMsg !== undefined && object.lastMsg !== null) {
            message.lastMsg = object.lastMsg;
        }
        else {
            message.lastMsg = new Uint8Array();
        }
        if (object.initFunds !== undefined && object.initFunds !== null) {
            for (const e of object.initFunds) {
                message.initFunds.push(Coin.fromPartial(e));
            }
        }
        if (object.callbackSig !== undefined && object.callbackSig !== null) {
            message.callbackSig = object.callbackSig;
        }
        else {
            message.callbackSig = new Uint8Array();
        }
        return message;
    },
};
const baseMsgExecuteContract = { callbackCodeHash: "" };
export const MsgExecuteContract = {
    encode(message, writer = Writer.create()) {
        if (message.sender.length !== 0) {
            writer.uint32(10).bytes(message.sender);
        }
        if (message.contract.length !== 0) {
            writer.uint32(18).bytes(message.contract);
        }
        if (message.msg.length !== 0) {
            writer.uint32(26).bytes(message.msg);
        }
        if (message.callbackCodeHash !== "") {
            writer.uint32(34).string(message.callbackCodeHash);
        }
        for (const v of message.sentFunds) {
            Coin.encode(v, writer.uint32(42).fork()).ldelim();
        }
        if (message.callbackSig.length !== 0) {
            writer.uint32(50).bytes(message.callbackSig);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgExecuteContract };
        message.sentFunds = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.bytes();
                    break;
                case 2:
                    message.contract = reader.bytes();
                    break;
                case 3:
                    message.msg = reader.bytes();
                    break;
                case 4:
                    message.callbackCodeHash = reader.string();
                    break;
                case 5:
                    message.sentFunds.push(Coin.decode(reader, reader.uint32()));
                    break;
                case 6:
                    message.callbackSig = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgExecuteContract };
        message.sentFunds = [];
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = bytesFromBase64(object.sender);
        }
        if (object.contract !== undefined && object.contract !== null) {
            message.contract = bytesFromBase64(object.contract);
        }
        if (object.msg !== undefined && object.msg !== null) {
            message.msg = bytesFromBase64(object.msg);
        }
        if (object.callbackCodeHash !== undefined &&
            object.callbackCodeHash !== null) {
            message.callbackCodeHash = String(object.callbackCodeHash);
        }
        else {
            message.callbackCodeHash = "";
        }
        if (object.sentFunds !== undefined && object.sentFunds !== null) {
            for (const e of object.sentFunds) {
                message.sentFunds.push(Coin.fromJSON(e));
            }
        }
        if (object.callbackSig !== undefined && object.callbackSig !== null) {
            message.callbackSig = bytesFromBase64(object.callbackSig);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined &&
            (obj.sender = base64FromBytes(message.sender !== undefined ? message.sender : new Uint8Array()));
        message.contract !== undefined &&
            (obj.contract = base64FromBytes(message.contract !== undefined ? message.contract : new Uint8Array()));
        message.msg !== undefined &&
            (obj.msg = base64FromBytes(message.msg !== undefined ? message.msg : new Uint8Array()));
        message.callbackCodeHash !== undefined &&
            (obj.callbackCodeHash = message.callbackCodeHash);
        if (message.sentFunds) {
            obj.sentFunds = message.sentFunds.map((e) => e ? Coin.toJSON(e) : undefined);
        }
        else {
            obj.sentFunds = [];
        }
        message.callbackSig !== undefined &&
            (obj.callbackSig = base64FromBytes(message.callbackSig !== undefined
                ? message.callbackSig
                : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgExecuteContract };
        message.sentFunds = [];
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = new Uint8Array();
        }
        if (object.contract !== undefined && object.contract !== null) {
            message.contract = object.contract;
        }
        else {
            message.contract = new Uint8Array();
        }
        if (object.msg !== undefined && object.msg !== null) {
            message.msg = object.msg;
        }
        else {
            message.msg = new Uint8Array();
        }
        if (object.callbackCodeHash !== undefined &&
            object.callbackCodeHash !== null) {
            message.callbackCodeHash = object.callbackCodeHash;
        }
        else {
            message.callbackCodeHash = "";
        }
        if (object.sentFunds !== undefined && object.sentFunds !== null) {
            for (const e of object.sentFunds) {
                message.sentFunds.push(Coin.fromPartial(e));
            }
        }
        if (object.callbackSig !== undefined && object.callbackSig !== null) {
            message.callbackSig = object.callbackSig;
        }
        else {
            message.callbackSig = new Uint8Array();
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
