---
sidebar_position: 2
title: Message Authentication
description: Message Authentication in the Intent Module
---

In the Intent engine, message authentication is a critical part of ensuring that flows are authorized and executed by the correct entities. This document explains the authentication process for various message types, with a focus on local messages, `MsgExec` messages, and ICA messages.

## Authentication Overview

The module processes three primary message types:

1. **Local Messages**:
   - These are standard messages sent directly by users.
   - Authentication ensures that the signer matches the flow owner specified in the request.

2. **Authz `MsgExec` Messages**:
   - These messages allow one account to execute messages on behalf of another.
   - Authentication involves verifying the signers for all contained messages to ensure proper delegation.

3. **Proxy Account Messages**:
   - These are messages sent via IBC (Inter-Blockchain Communication) through a Proxy Account (ICA or Union Proxy).
   - **Trustless Agent** messages can only call `MsgExec` for security reasons (requiring AuthZ grants on the host), while **Self-Hosted Proxy** messages are allowed without restrictions (as the user controls the account).

---

## Authentication Logic

### Local Messages

Local messages require a direct match between the signer and the flow owner. The module validates this relationship explicitly:

```go
if isLocalMessage(flow) {
    return k.validateSigners(ctx, codec, flow, message)
}
```

### Authz `MsgExec` Messages

For `MsgExec` messages, authentication is applied to each inner message:

```go
if isAuthzMsgExec(message) {
    return k.validateAuthzMsg(ctx, codec, flow, message)
}
```

This ensures that delegation rules are respected and all flows performed on behalf of another account are authorized.

### Trustless Agent Messages

Trustless Agent messages are authenticated differently. Since the proxy execution happens via IBC (either ICA packet or Union ZK proof), the packet source is authenticated by the transport layer. Our module controls which Trustless Agent to use through the configuration provided by the user.

For Trustless Agents (Cosmos), messages are restricted to `MsgExec` only. This ensures the agent can only perform actions for which the user has explicitly granted authorization (via AuthZ) on the host chain.

```go
if isTrustlessAgentMessage(flow) {
    if message.TypeUrl != sdk.MsgTypeURL(&authztypes.MsgExec{}) {
        return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only MsgExec is allowed for Trustless Agent messages")
    }
    return nil
}
```

### Self-Hosted Proxy Messages

Self-Hosted Proxy messages (e.g., standard ICA or Union EVM execution) do not require additional restrictions from the Intent Engine's perspective. They are trusted based on the ownership of the flow.

```go
if isSelfHostedICAMessage(flow) {
    return nil
}
```

---

## Validation

1. **ICA Authentication**:
   - IBC ensures that the packet sender is authenticated via AuthenticateTx as part of its [protocol](https://tutorials.cosmos.network/academy/3-ibc/8-ica.html#authentication).
   - This removes the need for additional signer checks within the Intent engine.

2. **Controlled Configuration**:
   - The flow submission from the user specifies which Trustless Execution Agent (and the fee configuration thereof) is used, and this configuration is already expected and verified during setup.

3. **Security for `MsgExec`**:
   - Restricting Trustless Execution Agent messages to `MsgExec` ensures that only delegated flows are performed, maintaining the security model.
  
---

## Example Code: `validateMessage`

The following code demonstrates how the module handles authentication for different message types:

```go
func (k Keeper) validateMessage(ctx sdk.Context, codec codec.Codec, flow types.Flow, message *codectypes.Any) error {
    var sdkMsg sdk.Msg
    if err := codec.UnpackAny(message, &sdkMsg); err != nil {
        return errorsmod.Wrap(err, "failed to unpack message")
    }

    switch {
    case isAuthzMsgExec(message):
        // Validate Authz MsgExec messages.
        return k.validateAuthzMsg(ctx, codec, flow, message)

    case isLocalMessage(flow):
        // Validate local messages to ensure the signer matches the owner.
        return k.validateSigners(ctx, codec, flow, message)

    case isTrustlessAgentMessage(flow):
		// Restrict Trustless Agent messages to MsgExec for security.
		if message.TypeUrl != sdk.MsgTypeURL(&authztypes.MsgExec{}) {
			return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only MsgExec is allowed for Trustless Agent messages")
		}
		return nil

    case isSelfHostedICAMessage(flow):
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

Message authentication is handled carefully within the Intent engine to ensure security and correctness. While local and `MsgExec` messages require explicit validation, ICA messages rely on IBC’s inherent authentication mechanisms and user-controlled configurations. Restricting Trustless Execution Agent messages to `MsgExec` adds an additional layer of security, while Self-hosted ICAs are trusted based on owner control. This approach balances security with efficiency, adhering to the principles of our module’s design.
