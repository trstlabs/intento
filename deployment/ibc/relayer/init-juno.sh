rly chains add --file $HOME/trst/deployment/ibc/relayer/interchain-acc-config/chains/trstdev-1.json trstdev-1
rly chains add --file $HOME/trst/deployment/ibc/relayer/interchain-acc-config/chains/juno.json localjuno

rly keys restore trstdev-1 testkey "grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"
rly keys restore localjuno testkey "crop staff genuine enjoy dial pact sorry bless note fall abuse more cheese clutch make ripple machine this gravity lend thank marine sell print"

rly paths add trstdev-1 testing trstdev1-localjuno --file $HOME/trst/deployment/ibc/relayer/interchain-acc-config/paths/trstdev1-localjuno.json
rly tx connection trstdev1-localjuno --override
rly start trstdev1-localjuno -p events -b 100 --debug > rly.log