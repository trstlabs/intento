/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { ContractInfo, ContractInfoWithAddress } from "../compute/types";
import { Duration } from "../google/protobuf/duration";
import { StringEvent } from "../cosmos/base/abci/v1beta1/abci";
import { Empty } from "../google/protobuf/empty";
export const protobufPackage = "trst.x.compute.v1beta1";
const baseQueryContractInfoRequest = {};
export const QueryContractInfoRequest = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractInfoRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractInfoRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractInfoRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        return message;
    },
};
const baseQueryContractInfoResponse = {};
export const QueryContractInfoResponse = {
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
            ...baseQueryContractInfoResponse,
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
            ...baseQueryContractInfoResponse,
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
            ...baseQueryContractInfoResponse,
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
const baseQueryContractResultRequest = {};
export const QueryContractResultRequest = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractResultRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractResultRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractResultRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        return message;
    },
};
const baseQueryContractResultResponse = { log: "" };
export const QueryContractResultResponse = {
    encode(message, writer = Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractResultResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractResultResponse,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        if (object.data !== undefined && object.data !== null) {
            message.data = bytesFromBase64(object.data);
        }
        if (object.log !== undefined && object.log !== null) {
            message.log = String(object.log);
        }
        else {
            message.log = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        message.data !== undefined &&
            (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
        message.log !== undefined && (obj.log = message.log);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractResultResponse,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        if (object.data !== undefined && object.data !== null) {
            message.data = object.data;
        }
        else {
            message.data = new Uint8Array();
        }
        if (object.log !== undefined && object.log !== null) {
            message.log = object.log;
        }
        else {
            message.log = "";
        }
        return message;
    },
};
const baseQueryContractHistoryRequest = {};
export const QueryContractHistoryRequest = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractHistoryRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractHistoryRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractHistoryRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        return message;
    },
};
const baseQueryContractsByCodeRequest = { codeId: 0 };
export const QueryContractsByCodeRequest = {
    encode(message, writer = Writer.create()) {
        if (message.codeId !== 0) {
            writer.uint32(8).uint64(message.codeId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractsByCodeRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeId = longToNumber(reader.uint64());
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
            ...baseQueryContractsByCodeRequest,
        };
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = Number(object.codeId);
        }
        else {
            message.codeId = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.codeId !== undefined && (obj.codeId = message.codeId);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractsByCodeRequest,
        };
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = object.codeId;
        }
        else {
            message.codeId = 0;
        }
        return message;
    },
};
const baseQueryContractsByCodeResponse = {};
export const QueryContractsByCodeResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.contractInfos) {
            ContractInfoWithAddress.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractsByCodeResponse,
        };
        message.contractInfos = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.contractInfos.push(ContractInfoWithAddress.decode(reader, reader.uint32()));
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
            ...baseQueryContractsByCodeResponse,
        };
        message.contractInfos = [];
        if (object.contractInfos !== undefined && object.contractInfos !== null) {
            for (const e of object.contractInfos) {
                message.contractInfos.push(ContractInfoWithAddress.fromJSON(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.contractInfos) {
            obj.contractInfos = message.contractInfos.map((e) => e ? ContractInfoWithAddress.toJSON(e) : undefined);
        }
        else {
            obj.contractInfos = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractsByCodeResponse,
        };
        message.contractInfos = [];
        if (object.contractInfos !== undefined && object.contractInfos !== null) {
            for (const e of object.contractInfos) {
                message.contractInfos.push(ContractInfoWithAddress.fromPartial(e));
            }
        }
        return message;
    },
};
const baseQuerySmartContractStateRequest = {};
export const QuerySmartContractStateRequest = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        if (message.queryData.length !== 0) {
            writer.uint32(18).bytes(message.queryData);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQuerySmartContractStateRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQuerySmartContractStateRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        if (object.queryData !== undefined && object.queryData !== null) {
            message.queryData = bytesFromBase64(object.queryData);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        message.queryData !== undefined &&
            (obj.queryData = base64FromBytes(message.queryData !== undefined ? message.queryData : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQuerySmartContractStateRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        if (object.queryData !== undefined && object.queryData !== null) {
            message.queryData = object.queryData;
        }
        else {
            message.queryData = new Uint8Array();
        }
        return message;
    },
};
const baseQueryContractAddressByContractIdRequest = { contractId: "" };
export const QueryContractAddressByContractIdRequest = {
    encode(message, writer = Writer.create()) {
        if (message.contractId !== "") {
            writer.uint32(10).string(message.contractId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractAddressByContractIdRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractAddressByContractIdRequest,
        };
        if (object.contractId !== undefined && object.contractId !== null) {
            message.contractId = String(object.contractId);
        }
        else {
            message.contractId = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.contractId !== undefined && (obj.contractId = message.contractId);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractAddressByContractIdRequest,
        };
        if (object.contractId !== undefined && object.contractId !== null) {
            message.contractId = object.contractId;
        }
        else {
            message.contractId = "";
        }
        return message;
    },
};
const baseQueryContractKeyRequest = {};
export const QueryContractKeyRequest = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractKeyRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractKeyRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractKeyRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        return message;
    },
};
const baseQueryContractHashRequest = {};
export const QueryContractHashRequest = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractHashRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractHashRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractHashRequest,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        return message;
    },
};
const baseQuerySmartContractStateResponse = {};
export const QuerySmartContractStateResponse = {
    encode(message, writer = Writer.create()) {
        if (message.data.length !== 0) {
            writer.uint32(10).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQuerySmartContractStateResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQuerySmartContractStateResponse,
        };
        if (object.data !== undefined && object.data !== null) {
            message.data = bytesFromBase64(object.data);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined &&
            (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQuerySmartContractStateResponse,
        };
        if (object.data !== undefined && object.data !== null) {
            message.data = object.data;
        }
        else {
            message.data = new Uint8Array();
        }
        return message;
    },
};
const baseQueryCodeRequest = { codeId: 0 };
export const QueryCodeRequest = {
    encode(message, writer = Writer.create()) {
        if (message.codeId !== 0) {
            writer.uint32(8).uint64(message.codeId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseQueryCodeRequest };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeId = longToNumber(reader.uint64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseQueryCodeRequest };
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = Number(object.codeId);
        }
        else {
            message.codeId = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.codeId !== undefined && (obj.codeId = message.codeId);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseQueryCodeRequest };
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = object.codeId;
        }
        else {
            message.codeId = 0;
        }
        return message;
    },
};
const baseCodeInfoResponse = { codeId: 0, source: "", builder: "" };
export const CodeInfoResponse = {
    encode(message, writer = Writer.create()) {
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
            Duration.encode(message.contractDuration, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseCodeInfoResponse };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeId = longToNumber(reader.uint64());
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
    fromJSON(object) {
        const message = { ...baseCodeInfoResponse };
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = Number(object.codeId);
        }
        else {
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
        if (object.contractDuration !== undefined &&
            object.contractDuration !== null) {
            message.contractDuration = Duration.fromJSON(object.contractDuration);
        }
        else {
            message.contractDuration = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.codeId !== undefined && (obj.codeId = message.codeId);
        message.creator !== undefined &&
            (obj.creator = base64FromBytes(message.creator !== undefined ? message.creator : new Uint8Array()));
        message.codeHash !== undefined &&
            (obj.codeHash = base64FromBytes(message.codeHash !== undefined ? message.codeHash : new Uint8Array()));
        message.source !== undefined && (obj.source = message.source);
        message.builder !== undefined && (obj.builder = message.builder);
        message.contractDuration !== undefined &&
            (obj.contractDuration = message.contractDuration
                ? Duration.toJSON(message.contractDuration)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseCodeInfoResponse };
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
        if (object.codeHash !== undefined && object.codeHash !== null) {
            message.codeHash = object.codeHash;
        }
        else {
            message.codeHash = new Uint8Array();
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
        if (object.contractDuration !== undefined &&
            object.contractDuration !== null) {
            message.contractDuration = Duration.fromPartial(object.contractDuration);
        }
        else {
            message.contractDuration = undefined;
        }
        return message;
    },
};
const baseQueryCodeResponse = {};
export const QueryCodeResponse = {
    encode(message, writer = Writer.create()) {
        if (message.codeInfo !== undefined) {
            CodeInfoResponse.encode(message.codeInfo, writer.uint32(10).fork()).ldelim();
        }
        if (message.data.length !== 0) {
            writer.uint32(18).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseQueryCodeResponse };
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
    fromJSON(object) {
        const message = { ...baseQueryCodeResponse };
        if (object.codeInfo !== undefined && object.codeInfo !== null) {
            message.codeInfo = CodeInfoResponse.fromJSON(object.codeInfo);
        }
        else {
            message.codeInfo = undefined;
        }
        if (object.data !== undefined && object.data !== null) {
            message.data = bytesFromBase64(object.data);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.codeInfo !== undefined &&
            (obj.codeInfo = message.codeInfo
                ? CodeInfoResponse.toJSON(message.codeInfo)
                : undefined);
        message.data !== undefined &&
            (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseQueryCodeResponse };
        if (object.codeInfo !== undefined && object.codeInfo !== null) {
            message.codeInfo = CodeInfoResponse.fromPartial(object.codeInfo);
        }
        else {
            message.codeInfo = undefined;
        }
        if (object.data !== undefined && object.data !== null) {
            message.data = object.data;
        }
        else {
            message.data = new Uint8Array();
        }
        return message;
    },
};
const baseQueryCodesResponse = {};
export const QueryCodesResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.codeInfos) {
            CodeInfoResponse.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseQueryCodesResponse };
        message.codeInfos = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeInfos.push(CodeInfoResponse.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseQueryCodesResponse };
        message.codeInfos = [];
        if (object.codeInfos !== undefined && object.codeInfos !== null) {
            for (const e of object.codeInfos) {
                message.codeInfos.push(CodeInfoResponse.fromJSON(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.codeInfos) {
            obj.codeInfos = message.codeInfos.map((e) => e ? CodeInfoResponse.toJSON(e) : undefined);
        }
        else {
            obj.codeInfos = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseQueryCodesResponse };
        message.codeInfos = [];
        if (object.codeInfos !== undefined && object.codeInfos !== null) {
            for (const e of object.codeInfos) {
                message.codeInfos.push(CodeInfoResponse.fromPartial(e));
            }
        }
        return message;
    },
};
const baseQueryContractAddressByContractIdResponse = {};
export const QueryContractAddressByContractIdResponse = {
    encode(message, writer = Writer.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractAddressByContractIdResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractAddressByContractIdResponse,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = bytesFromBase64(object.address);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = base64FromBytes(message.address !== undefined ? message.address : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractAddressByContractIdResponse,
        };
        if (object.address !== undefined && object.address !== null) {
            message.address = object.address;
        }
        else {
            message.address = new Uint8Array();
        }
        return message;
    },
};
const baseQueryContractKeyResponse = {};
export const QueryContractKeyResponse = {
    encode(message, writer = Writer.create()) {
        if (message.key.length !== 0) {
            writer.uint32(10).bytes(message.key);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractKeyResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractKeyResponse,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = bytesFromBase64(object.key);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined &&
            (obj.key = base64FromBytes(message.key !== undefined ? message.key : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractKeyResponse,
        };
        if (object.key !== undefined && object.key !== null) {
            message.key = object.key;
        }
        else {
            message.key = new Uint8Array();
        }
        return message;
    },
};
const baseQueryContractHashResponse = {};
export const QueryContractHashResponse = {
    encode(message, writer = Writer.create()) {
        if (message.codeHash.length !== 0) {
            writer.uint32(10).bytes(message.codeHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryContractHashResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryContractHashResponse,
        };
        if (object.codeHash !== undefined && object.codeHash !== null) {
            message.codeHash = bytesFromBase64(object.codeHash);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.codeHash !== undefined &&
            (obj.codeHash = base64FromBytes(message.codeHash !== undefined ? message.codeHash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryContractHashResponse,
        };
        if (object.codeHash !== undefined && object.codeHash !== null) {
            message.codeHash = object.codeHash;
        }
        else {
            message.codeHash = new Uint8Array();
        }
        return message;
    },
};
const baseDecryptedAnswer = {
    type: "",
    input: "",
    outputData: "",
    outputDataAsString: "",
    plaintextError: "",
};
export const DecryptedAnswer = {
    encode(message, writer = Writer.create()) {
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
            StringEvent.encode(v, writer.uint32(42).fork()).ldelim();
        }
        if (message.outputError.length !== 0) {
            writer.uint32(50).bytes(message.outputError);
        }
        if (message.plaintextError !== "") {
            writer.uint32(58).string(message.plaintextError);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseDecryptedAnswer };
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
    fromJSON(object) {
        const message = { ...baseDecryptedAnswer };
        message.outputLogs = [];
        if (object.type !== undefined && object.type !== null) {
            message.type = String(object.type);
        }
        else {
            message.type = "";
        }
        if (object.input !== undefined && object.input !== null) {
            message.input = String(object.input);
        }
        else {
            message.input = "";
        }
        if (object.outputData !== undefined && object.outputData !== null) {
            message.outputData = String(object.outputData);
        }
        else {
            message.outputData = "";
        }
        if (object.outputDataAsString !== undefined &&
            object.outputDataAsString !== null) {
            message.outputDataAsString = String(object.outputDataAsString);
        }
        else {
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
        }
        else {
            message.plaintextError = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.type !== undefined && (obj.type = message.type);
        message.input !== undefined && (obj.input = message.input);
        message.outputData !== undefined && (obj.outputData = message.outputData);
        message.outputDataAsString !== undefined &&
            (obj.outputDataAsString = message.outputDataAsString);
        if (message.outputLogs) {
            obj.outputLogs = message.outputLogs.map((e) => e ? StringEvent.toJSON(e) : undefined);
        }
        else {
            obj.outputLogs = [];
        }
        message.outputError !== undefined &&
            (obj.outputError = base64FromBytes(message.outputError !== undefined
                ? message.outputError
                : new Uint8Array()));
        message.plaintextError !== undefined &&
            (obj.plaintextError = message.plaintextError);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseDecryptedAnswer };
        message.outputLogs = [];
        if (object.type !== undefined && object.type !== null) {
            message.type = object.type;
        }
        else {
            message.type = "";
        }
        if (object.input !== undefined && object.input !== null) {
            message.input = object.input;
        }
        else {
            message.input = "";
        }
        if (object.outputData !== undefined && object.outputData !== null) {
            message.outputData = object.outputData;
        }
        else {
            message.outputData = "";
        }
        if (object.outputDataAsString !== undefined &&
            object.outputDataAsString !== null) {
            message.outputDataAsString = object.outputDataAsString;
        }
        else {
            message.outputDataAsString = "";
        }
        if (object.outputLogs !== undefined && object.outputLogs !== null) {
            for (const e of object.outputLogs) {
                message.outputLogs.push(StringEvent.fromPartial(e));
            }
        }
        if (object.outputError !== undefined && object.outputError !== null) {
            message.outputError = object.outputError;
        }
        else {
            message.outputError = new Uint8Array();
        }
        if (object.plaintextError !== undefined && object.plaintextError !== null) {
            message.plaintextError = object.plaintextError;
        }
        else {
            message.plaintextError = "";
        }
        return message;
    },
};
export class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    ContractInfo(request) {
        const data = QueryContractInfoRequest.encode(request).finish();
        const promise = this.rpc.request("trst.x.compute.v1beta1.Query", "ContractInfo", data);
        return promise.then((data) => QueryContractInfoResponse.decode(new Reader(data)));
    }
    ContractResult(request) {
        const data = QueryContractResultRequest.encode(request).finish();
        const promise = this.rpc.request("trst.x.compute.v1beta1.Query", "ContractResult", data);
        return promise.then((data) => QueryContractResultResponse.decode(new Reader(data)));
    }
    ContractsByCode(request) {
        const data = QueryContractsByCodeRequest.encode(request).finish();
        const promise = this.rpc.request("trst.x.compute.v1beta1.Query", "ContractsByCode", data);
        return promise.then((data) => QueryContractsByCodeResponse.decode(new Reader(data)));
    }
    SmartContractState(request) {
        const data = QuerySmartContractStateRequest.encode(request).finish();
        const promise = this.rpc.request("trst.x.compute.v1beta1.Query", "SmartContractState", data);
        return promise.then((data) => QuerySmartContractStateResponse.decode(new Reader(data)));
    }
    Code(request) {
        const data = QueryCodeRequest.encode(request).finish();
        const promise = this.rpc.request("trst.x.compute.v1beta1.Query", "Code", data);
        return promise.then((data) => QueryCodeResponse.decode(new Reader(data)));
    }
    Codes(request) {
        const data = Empty.encode(request).finish();
        const promise = this.rpc.request("trst.x.compute.v1beta1.Query", "Codes", data);
        return promise.then((data) => QueryCodesResponse.decode(new Reader(data)));
    }
}
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
