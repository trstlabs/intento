/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { MasterCertificate } from "../registration/msg";
export const protobufPackage = "trst.x.registration.v1beta1";
const baseQueryMasterKeyRequest = {};
export const QueryMasterKeyRequest = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseQueryMasterKeyRequest };
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
    fromJSON(_) {
        const message = { ...baseQueryMasterKeyRequest };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseQueryMasterKeyRequest };
        return message;
    },
};
const baseQueryMasterKeyResponse = {};
export const QueryMasterKeyResponse = {
    encode(message, writer = Writer.create()) {
        if (message.masterKey !== undefined) {
            MasterCertificate.encode(message.masterKey, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseQueryMasterKeyResponse };
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
    fromJSON(object) {
        const message = { ...baseQueryMasterKeyResponse };
        if (object.masterKey !== undefined && object.masterKey !== null) {
            message.masterKey = MasterCertificate.fromJSON(object.masterKey);
        }
        else {
            message.masterKey = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.masterKey !== undefined &&
            (obj.masterKey = message.masterKey
                ? MasterCertificate.toJSON(message.masterKey)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseQueryMasterKeyResponse };
        if (object.masterKey !== undefined && object.masterKey !== null) {
            message.masterKey = MasterCertificate.fromPartial(object.masterKey);
        }
        else {
            message.masterKey = undefined;
        }
        return message;
    },
};
const baseQueryEncryptedSeedRequest = {};
export const QueryEncryptedSeedRequest = {
    encode(message, writer = Writer.create()) {
        if (message.pubKey.length !== 0) {
            writer.uint32(10).bytes(message.pubKey);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryEncryptedSeedRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryEncryptedSeedRequest,
        };
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = bytesFromBase64(object.pubKey);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.pubKey !== undefined &&
            (obj.pubKey = base64FromBytes(message.pubKey !== undefined ? message.pubKey : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryEncryptedSeedRequest,
        };
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = object.pubKey;
        }
        else {
            message.pubKey = new Uint8Array();
        }
        return message;
    },
};
const baseQueryEncryptedSeedResponse = {};
export const QueryEncryptedSeedResponse = {
    encode(message, writer = Writer.create()) {
        if (message.encryptedSeed.length !== 0) {
            writer.uint32(10).bytes(message.encryptedSeed);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryEncryptedSeedResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryEncryptedSeedResponse,
        };
        if (object.encryptedSeed !== undefined && object.encryptedSeed !== null) {
            message.encryptedSeed = bytesFromBase64(object.encryptedSeed);
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.encryptedSeed !== undefined &&
            (obj.encryptedSeed = base64FromBytes(message.encryptedSeed !== undefined
                ? message.encryptedSeed
                : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryEncryptedSeedResponse,
        };
        if (object.encryptedSeed !== undefined && object.encryptedSeed !== null) {
            message.encryptedSeed = object.encryptedSeed;
        }
        else {
            message.encryptedSeed = new Uint8Array();
        }
        return message;
    },
};
export class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    MasterKey(request) {
        const data = QueryMasterKeyRequest.encode(request).finish();
        const promise = this.rpc.request("trst.x.registration.v1beta1.Query", "MasterKey", data);
        return promise.then((data) => QueryMasterKeyResponse.decode(new Reader(data)));
    }
    EncryptedSeed(request) {
        const data = QueryEncryptedSeedRequest.encode(request).finish();
        const promise = this.rpc.request("trst.x.registration.v1beta1.Query", "EncryptedSeed", data);
        return promise.then((data) => QueryEncryptedSeedResponse.decode(new Reader(data)));
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
