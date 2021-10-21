/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";
export const protobufPackage = "trst.x.registration.v1beta1";
const baseRaAuthenticate = {};
export const RaAuthenticate = {
    encode(message, writer = Writer.create()) {
        if (message.sender.length !== 0) {
            writer.uint32(10).bytes(message.sender);
        }
        if (message.certificate.length !== 0) {
            writer.uint32(18).bytes(message.certificate);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseRaAuthenticate };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.bytes();
                    break;
                case 2:
                    message.certificate = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseRaAuthenticate };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = bytesFromBase64(object.sender);
        }
        if (object.certificate !== undefined && object.certificate !== null) {
            message.certificate = bytesFromBase64(object.certificate);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined &&
            (obj.sender = base64FromBytes(message.sender !== undefined ? message.sender : new Uint8Array()));
        message.certificate !== undefined &&
            (obj.certificate = base64FromBytes(message.certificate !== undefined
                ? message.certificate
                : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseRaAuthenticate };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = new Uint8Array();
        }
        if (object.certificate !== undefined && object.certificate !== null) {
            message.certificate = object.certificate;
        }
        else {
            message.certificate = new Uint8Array();
        }
        return message;
    },
};
const baseMasterCertificate = {};
export const MasterCertificate = {
    encode(message, writer = Writer.create()) {
        if (message.bytes.length !== 0) {
            writer.uint32(10).bytes(message.bytes);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMasterCertificate };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.bytes = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMasterCertificate };
        if (object.bytes !== undefined && object.bytes !== null) {
            message.bytes = bytesFromBase64(object.bytes);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.bytes !== undefined &&
            (obj.bytes = base64FromBytes(message.bytes !== undefined ? message.bytes : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMasterCertificate };
        if (object.bytes !== undefined && object.bytes !== null) {
            message.bytes = object.bytes;
        }
        else {
            message.bytes = new Uint8Array();
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
