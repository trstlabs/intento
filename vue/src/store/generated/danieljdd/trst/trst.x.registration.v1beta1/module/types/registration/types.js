/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";
export const protobufPackage = "trst.x.registration.v1beta1";
const baseSeedConfig = { masterCert: "", encryptedKey: "" };
export const SeedConfig = {
    encode(message, writer = Writer.create()) {
        if (message.masterCert !== "") {
            writer.uint32(10).string(message.masterCert);
        }
        if (message.encryptedKey !== "") {
            writer.uint32(18).string(message.encryptedKey);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseSeedConfig };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.masterCert = reader.string();
                    break;
                case 2:
                    message.encryptedKey = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseSeedConfig };
        if (object.masterCert !== undefined && object.masterCert !== null) {
            message.masterCert = String(object.masterCert);
        }
        else {
            message.masterCert = "";
        }
        if (object.encryptedKey !== undefined && object.encryptedKey !== null) {
            message.encryptedKey = String(object.encryptedKey);
        }
        else {
            message.encryptedKey = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.masterCert !== undefined && (obj.masterCert = message.masterCert);
        message.encryptedKey !== undefined &&
            (obj.encryptedKey = message.encryptedKey);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseSeedConfig };
        if (object.masterCert !== undefined && object.masterCert !== null) {
            message.masterCert = object.masterCert;
        }
        else {
            message.masterCert = "";
        }
        if (object.encryptedKey !== undefined && object.encryptedKey !== null) {
            message.encryptedKey = object.encryptedKey;
        }
        else {
            message.encryptedKey = "";
        }
        return message;
    },
};
const baseRegistrationNodeInfo = {};
export const RegistrationNodeInfo = {
    encode(message, writer = Writer.create()) {
        if (message.certificate.length !== 0) {
            writer.uint32(10).bytes(message.certificate);
        }
        if (message.encryptedSeed.length !== 0) {
            writer.uint32(18).bytes(message.encryptedSeed);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseRegistrationNodeInfo };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.certificate = reader.bytes();
                    break;
                case 2:
                    message.encryptedSeed = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseRegistrationNodeInfo };
        if (object.certificate !== undefined && object.certificate !== null) {
            message.certificate = bytesFromBase64(object.certificate);
        }
        if (object.encryptedSeed !== undefined && object.encryptedSeed !== null) {
            message.encryptedSeed = bytesFromBase64(object.encryptedSeed);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.certificate !== undefined &&
            (obj.certificate = base64FromBytes(message.certificate !== undefined
                ? message.certificate
                : new Uint8Array()));
        message.encryptedSeed !== undefined &&
            (obj.encryptedSeed = base64FromBytes(message.encryptedSeed !== undefined
                ? message.encryptedSeed
                : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseRegistrationNodeInfo };
        if (object.certificate !== undefined && object.certificate !== null) {
            message.certificate = object.certificate;
        }
        else {
            message.certificate = new Uint8Array();
        }
        if (object.encryptedSeed !== undefined && object.encryptedSeed !== null) {
            message.encryptedSeed = object.encryptedSeed;
        }
        else {
            message.encryptedSeed = new Uint8Array();
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
