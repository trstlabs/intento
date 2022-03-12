# NOTE
The module is under development, There are 2 methods to store code. One by Governance, one directly. In the mainnet chain code will be stored through governance. 

The other store code method is only enabled for testing. 

# Wasm Module

This should be a brief overview of the functionality

## Parameters

Through governance, a number of compute parameters can be adjusted. The default values are the following:
```golang

const (
	DefaultMaxContractPeriod time.Duration = time.Hour * 24 * 30 // 30 days

	// Commission percentage to distribute to community pool for leftover balances (rounded up)

	DefaultCommission int64 = 2 // %

	// MinContractDurationForIncentive to distribute reward to contract

	DefaultMinContractDurationForIncentive time.Duration = time.Hour * 24 // 1 day

	// DefaultMaxContractIncentive max amount of utrst coins to give to a contract as incentive

	DefaultMaxContractIncentive int64 = 500000000 // 500utrst

	// MinContractBalanceForIncentive minimum balance required to be elligable for an incentive

	DefaultMinContractBalanceForIncentive int64 = 50000000 // 50utrst
)
```


## Endblocker (-> AfterEpochEnd)
At the end of each block, the Endblock checks if there are Trustless Contracts that are elligable for an incentive (will be Epoch based once Epoch module is in SDK).
MinContractDurationForIncentive,MaxContractIncentive, MinContractBalanceForIncentive can be adjusted by governane to ensure that the contract incentives are distributed in a fair manner.

## Configuration

You can add the following section to `config/app.toml`. Below is shown with defaults:

```toml
[wasm]
# This is the maximum sdk gas (wasm and storage) that we allow for any x/compute "smart" queries
query_gas_limit = 300000
# This is the number of wasm vm instances we keep cached in memory for speed-up
# Warning: this is currently unstable and may lead to crashes, best to keep for 0 unless testing locally
lru_size = 0
```

## Events

A number of events are returned to allow good indexing of the transactions from smart contracts.

Every call to Instantiate or Execute will be tagged with the info on the contract that was executed and who executed it.
It should look something like this (with different addresses). The module is always `wasm`, and `code_id` is only present
when Instantiating a contract, so you can subscribe to new instances, it is omitted on Execute. There is also an `action` tag
which is auto-added by the Cosmos SDK and has a value of either `store-code`, `instantiate` or `execute` depending on which message
was sent:

```json
{
  "Type": "message",
  "Attr": [
    {
      "key": "module",
      "value": "wasm"
    },
    {
      "key": "action",
      "value": "instantiate"
    },
    {
      "key": "signer",
      "value": "trust1vx8knpllrj7n963p9ttd80w47kpacrhuts497x"
    },
    {
      "key": "code_id",
      "value": "1"
    },
    {
      "key": "contract_address",
      "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n8978vsyg"
    }
  ]
}
```

If any funds were transferred to the contract as part of the message, or if the contract released funds as part of it's executions,
it will receive the typical events associated with sending tokens from bank. In this case, we instantiate the contract and
provide a initial balance in the same `MsgInstantiateContract`. We see the following events in addition to the above one:

```json
[
  {
    "Type": "transfer",
    "Attr": [
      {
        "key": "recipient",
        "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n8978vsyg"
      },
      {
        "key": "sender",
        "value": "trust1ffnqn02ft2psvyv4dyr56nnv6plllf9pm2kpmv"
      },
      {
        "key": "amount",
        "value": "100000denom"
      }
    ]
  }
]
```

Finally, the contract itself can emit a "custom event" on Execute only (not on Init).
There is one event per contract, so if one contract calls a second contract, you may receive
one event for the original contract and one for the re-invoked contract. All attributes from the contract are passed through verbatim,
and we add a `contract_address` attribute that contains the actual contract that emitted that event.
Here is an example from the escrow contract successfully releasing funds to the destination address:

```json
{
  "Type": "wasm",
  "Attr": [
    {
      "key": "contract_address",
      "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n8978vsyg"
    },
    {
      "key": "action",
      "value": "release"
    },
    {
      "key": "destination",
      "value": "trust14k7v7ms4jxkk2etmg9gljxjm4ru3qjdugfsflq"
    }
  ]
}
```

### Pulling this all together

We will invoke an escrow contract to release to the designated beneficiary.
The escrow was previously loaded with `100000denom` (from the above example).
In this transaction, we send `5000denom` along with the `MsgExecuteContract`
and the contract releases the entire funds (`105000denom`) to the beneficiary.

We will see all the following events, where you should be able to reconstruct the actions
(remember there are two events for each transfer). We see (1) the initial transfer of funds
to the contract, (2) the contract custom event that it released funds (3) the transfer of funds
from the contract to the beneficiary and (4) the generic x/compute event stating that the contract
was executed (which always appears, while 2 is optional and has information as reliable as the contract):

```json
[
  {
    "Type": "transfer",
    "Attr": [
      {
        "key": "recipient",
        "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n8978vsyg"
      },
      {
        "key": "sender",
        "value": "trust1zm074khx32hqy20hlshlsd423n07pwlu9cpt37"
      },
      {
        "key": "amount",
        "value": "5000denom"
      }
    ]
  },
  {
    "Type": "wasm",
    "Attr": [
      {
        "key": "contract_address",
        "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n8978vsyg"
      },
      {
        "key": "action",
        "value": "release"
      },
      {
        "key": "destination",
        "value": "trust14k7v7ms4jxkk2etmg9gljxjm4ru3qjdugfsflq"
      }
    ]
  },
  {
    "Type": "transfer",
    "Attr": [
      {
        "key": "recipient",
        "value": "trust14k7v7ms4jxkk2etmg9gljxjm4ru3qjdugfsflq"
      },
      {
        "key": "sender",
        "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n8978vsyg"
      },
      {
        "key": "amount",
        "value": "105000denom"
      }
    ]
  },
  {
    "Type": "message",
    "Attr": [
      {
        "key": "module",
        "value": "wasm"
      },
      {
        "key": "action",
        "value": "execute"
      },
      {
        "key": "signer",
        "value": "trust1zm074khx32hqy20hlshlsd423n07pwlu9cpt37"
      },
      {
        "key": "contract_address",
        "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n8978vsyg"
      }
    ]
  }
]
```

A note on this format. This is what we return from our module. However, it seems to me that many events with the same `Type`
get merged together somewhere along the stack, so in this case, you _may_ end up with one "transfer" event with the info for
both transfers. Double check when evaluating the event logs, I will document better with more experience, especially when I
find out the entire path for the events.
