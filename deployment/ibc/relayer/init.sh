rly chains add --file $HOME/trst/deployment/ibc/relayer/interchain-acc-config/chains/trstdev-1.json trstdev-1

rly chains add --file $HOME/trst/deployment/ibc/relayer/interchain-acc-config/chains/trstdev-2.json trstdev-2

rly keys restore trstdev-1 testkey "grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"
rly keys restore trstdev-2 testkey "jelly shadow frog dirt dragon use armed praise universe win jungle close inmate rain oil canvas beauty pioneer chef soccer icon dizzy thunder meadow"

rly paths add trstdev-1 trstdev-2 trstdev1-trstdev2 --file $HOME/trst/deployment/ibc/relayer/interchain-acc-config/paths/trstdev1-trstdev2.json

rly tx connection trstdev1-trstdev2 --override

rly start trstdev1-trstdev2 -p events -b 100 --debug