---
sidebar_position: 2
title: Message Authentication
description: Message Authentication in the Intent Module
---

In the Intent module, message authentication is a critical part of ensuring that flows are authorized and executed by the correct entities. This document explains the authentication process for various message types, with a focus on local messages, `MsgExec` messages, and ICA messages.

## Authentication Overview

The module processes three primary message types:

1. **Local Messages**:
   - These are standard messages sent directly by users.
   - Authentication ensures that the signer matches the flow owner specified in the request.

2. **Authz `MsgExec` Messages**:
   - These messages allow one account to execute messages on behalf of another.
   - Authentication involves verifying the signers for all contained messages to ensure proper delegation.

3. **ICA (Interchain Account) Messages**:
   - These are messages sent via IBC (Inter-Blockchain Communication) through an interchain account.
   - Trustless Execution Agent messages can only call `MsgExec` for security reasons, while Self-Hosted ICA messages are allowed without restrictions.

---

## Authentication Logic

### Local Messages

Local messages require a direct match between the signer and the flow owner. The module validates this relationship explicitly:

```go
if isLocalMessage(flowInfo) {
    return k.validateSigners(ctx, codec, flowInfo, message)
}
```

### Authz `MsgExec` Messages

For `MsgExec` messages, authentication is applied to each inner message:

```go
if isAuthzMsgExec(message) {
    return k.validateAuthzMsg(ctx, codec, flowInfo, message)
}
```

This ensures that delegation rules are respected and all flows performed on behalf of another account are authorized.

### Trustless Execution Agent Messages

Trustless Execution Agent messages are authenticated differently. Since ICA operates via IBC, the packet sender is already verified by the IBC protocol. Additionally, our module controls which Trustless Execution Agent to use through the configuration provided by the user. Because of this controlled environment, further signer validation is unnecessary for Trustless Execution Agent messages. However, Trustless Execution Agent messages are restricted to `MsgExec` only for added security:

```go
if isHostedICAMessage(flowInfo) {
    if message.TypeUrl != sdk.MsgTypeURL(&authztypes.MsgExec{}) {
        return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only MsgExec is allowed for Trustless Execution Agent messages")
    }
    return nil
}
```

### Self-hosted ICA Messages

Self-hosted ICA messages do not require additional restrictions. They are trusted based on the direct control of the owner. These messages are allowed without further validation:

```go
if isSelfHostedICAMessage(flowInfo) {
    return nil
}
```

---

## Validation

1. **ICA Authentication**:
   - IBC ensures that the packet sender is authenticated via AuthenticateTx as part of its [protocol](https://tutorials.cosmos.network/academy/3-ibc/8-ica.html#authentication).
   - This removes the need for additional signer checks within the Intent module.

2. **Controlled Configuration**:
   - The flow submission from the user specifies which Trustless Execution Agent (and the fee configuration thereof) is used, and this configuration is already expected and verified during setup.

3. **Security for `MsgExec`**:
   - Restricting Trustless Execution Agent messages to `MsgExec` ensures that only delegated flows are performed, maintaining the security model.
  
---

## Example Code: `validateMessage`

The following code demonstrates how the module handles authentication for different message types:

```go
func (k Keeper) validateMessage(ctx sdk.Context, codec codec.Codec, flowInfo types.FlowInfo, message *codectypes.Any) error {
    var sdkMsg sdk.Msg
    if err := codec.UnpackAny(message, &sdkMsg); err != nil {
        return errorsmod.Wrap(err, "failed to unpack message")
    }

    switch {
    case isAuthzMsgExec(message):
        // Validate Authz MsgExec messages.
        return k.validateAuthzMsg(ctx, codec, flowInfo, message)

    case isLocalMessage(flowInfo):
        // Validate local messages to ensure the signer matches the owner.
        return k.validateSigners(ctx, codec, flowInfo, message)

    case isHostedICAMessage(flowInfo):
        // Restrict Hosted ICA messages to MsgExec for security.
        if message.TypeUrl != sdk.MsgTypeURL(&authztypes.MsgExec{}) {
            return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only MsgExec is allowed for Hosted ICA messages")
        }
        return nil

    case isSelfHostedICAMessage(flowInfo):
        // Allow Self-hosted ICA messages without additional validation.
        return nil

    default:
        // Unsupported message type.
        return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unsupported message type")
    }
}
```

---

## Conclusion

Message authentication is handled carefully within the Intent module to ensure security and correctness. While local and `MsgExec` messages require explicit validation, ICA messages rely on IBC’s inherent authentication mechanisms and user-controlled configurations. Restricting Trustless Execution Agent messages to `MsgExec` adds an additional layer of security, while Self-hosted ICAs are trusted based on owner control. This approach balances security with efficiency, adhering to the principles of our module’s design.
