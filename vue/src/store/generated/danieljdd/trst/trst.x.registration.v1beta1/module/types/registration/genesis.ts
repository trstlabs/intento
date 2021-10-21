/* eslint-disable */
import { RegistrationNodeInfo } from "../registration/types";
import { MasterCertificate } from "../registration/msg";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "trst.x.registration.v1beta1";

export interface GenesisState {
  registration: RegistrationNodeInfo[];
  nodeExchMasterCertificate: MasterCertificate | undefined;
  ioMasterCertificate: MasterCertificate | undefined;
}

const baseGenesisState: object = {};

export const GenesisState = {
  encode(message: GenesisState, writer: Writer = Writer.create()): Writer {
    for (const v of message.registration) {
      RegistrationNodeInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.nodeExchMasterCertificate !== undefined) {
      MasterCertificate.encode(
        message.nodeExchMasterCertificate,
        writer.uint32(18).fork()
      ).ldelim();
    }
    if (message.ioMasterCertificate !== undefined) {
      MasterCertificate.encode(
        message.ioMasterCertificate,
        writer.uint32(26).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenesisState } as GenesisState;
    message.registration = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.registration.push(
            RegistrationNodeInfo.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.nodeExchMasterCertificate = MasterCertificate.decode(
            reader,
            reader.uint32()
          );
          break;
        case 3:
          message.ioMasterCertificate = MasterCertificate.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.registration = [];
    if (object.registration !== undefined && object.registration !== null) {
      for (const e of object.registration) {
        message.registration.push(RegistrationNodeInfo.fromJSON(e));
      }
    }
    if (
      object.nodeExchMasterCertificate !== undefined &&
      object.nodeExchMasterCertificate !== null
    ) {
      message.nodeExchMasterCertificate = MasterCertificate.fromJSON(
        object.nodeExchMasterCertificate
      );
    } else {
      message.nodeExchMasterCertificate = undefined;
    }
    if (
      object.ioMasterCertificate !== undefined &&
      object.ioMasterCertificate !== null
    ) {
      message.ioMasterCertificate = MasterCertificate.fromJSON(
        object.ioMasterCertificate
      );
    } else {
      message.ioMasterCertificate = undefined;
    }
    return message;
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    if (message.registration) {
      obj.registration = message.registration.map((e) =>
        e ? RegistrationNodeInfo.toJSON(e) : undefined
      );
    } else {
      obj.registration = [];
    }
    message.nodeExchMasterCertificate !== undefined &&
      (obj.nodeExchMasterCertificate = message.nodeExchMasterCertificate
        ? MasterCertificate.toJSON(message.nodeExchMasterCertificate)
        : undefined);
    message.ioMasterCertificate !== undefined &&
      (obj.ioMasterCertificate = message.ioMasterCertificate
        ? MasterCertificate.toJSON(message.ioMasterCertificate)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.registration = [];
    if (object.registration !== undefined && object.registration !== null) {
      for (const e of object.registration) {
        message.registration.push(RegistrationNodeInfo.fromPartial(e));
      }
    }
    if (
      object.nodeExchMasterCertificate !== undefined &&
      object.nodeExchMasterCertificate !== null
    ) {
      message.nodeExchMasterCertificate = MasterCertificate.fromPartial(
        object.nodeExchMasterCertificate
      );
    } else {
      message.nodeExchMasterCertificate = undefined;
    }
    if (
      object.ioMasterCertificate !== undefined &&
      object.ioMasterCertificate !== null
    ) {
      message.ioMasterCertificate = MasterCertificate.fromPartial(
        object.ioMasterCertificate
      );
    } else {
      message.ioMasterCertificate = undefined;
    }
    return message;
  },
};

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
