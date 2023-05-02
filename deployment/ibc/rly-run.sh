#!/bin/bash

# Set up chains, paths, connections, and channels using rly

# Chain config paths
chain_configs=(
    "/rlyconfig/chains/trstdev-1.json"
    "/rlyconfig/chains/trstdev-2.json"
    "/rlyconfig/chains/juno.json"
)

# Chain names
chain_names=(
    "trstdev-1"
    "trstdev-2"
    "localjuno"
)

# Path config paths
path_configs=(
    "/rlyconfig/paths/trstdev1-trstdev2.json"
    "/rlyconfig/paths/trstdev1-localjuno.json"
)

# Path names
path_names=(
    "trstdev1-trstdev2"
    "trstdev1-localjuno"
)

# Keys
keys=(
    "grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"
    "jelly shadow frog dirt dragon use armed praise universe win jungle close inmate rain oil canvas beauty pioneer chef soccer icon dizzy thunder meadow"
    "crop staff genuine enjoy dial pact sorry bless note fall abuse more cheese clutch make ripple machine this gravity lend thank marine sell print"
)

# Initialize rly config
rly config init

# Add chains
for ((i=0;i<${#chain_configs[@]};++i)); do
    rly chains add --file "${chain_configs[$i]}" "${chain_names[$i]}"
done

# Restore keys for each chain
for ((i=0;i<${#keys[@]};++i)); do
    for ((j=0;j<${#chain_names[@]};++j)); do
        rly keys restore "${chain_names[$j]}" testkey "${keys[$i]}"
    done
done

# Add paths
for ((i=0;i<${#path_configs[@]};++i)); do
    rly paths add "${chain_names[$i]}" "${path_names[$i]}" --file "${path_configs[$i]}"
done

# Create connections
for ((i=0;i<${#path_names[@]};++i)); do
    rly tx connection "${path_names[$i]}"
done

# Create channels
for ((i=0;i<${#path_names[@]};++i)); do
    rly tx channel "${path_names[$i]}"
done

# Start relayer
rly start "${path_names[@]}" -p events -b 100 --debug > rly.log


# rly config init
# rly chains add --file /rlyconfig/chains/trstdev-1.json trstdev-1
# rly chains add --file /rlyconfig/chains/trstdev-2.json trstdev-2
# rly chains add --file /rlyconfig/chains/juno.json localjuno

# rly keys restore trstdev-1 testkey "grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"
# rly keys restore localjuno testkey "crop staff genuine enjoy dial pact sorry bless note fall abuse more cheese clutch make ripple machine this gravity lend thank marine sell print"
# rly keys restore trstdev-2 testkey "jelly shadow frog dirt dragon use armed praise universe win jungle close inmate rain oil canvas beauty pioneer chef soccer icon dizzy thunder meadow"

# rly paths add trstdev-1 trstdev-1 trstdev1-trstdev2 --file /rlyconfig/paths/trstdev1-trstdev2.json
# rly paths add trstdev-1 testing trstdev1-localjuno --file /rlyconfig/paths/trstdev1-localjuno.json

# # connections for ICA
# rly tx connection trstdev1-trstdev2 
# rly tx connection trstdev1-localjuno
# # transfer channel 
# rly tx channel trstdev1-trstdev2
# rly tx channel trstdev1-localjuno
# rly start trstdev1-trstdev2 trstdev1-localjuno -p events -b 100  --debug > rly.log
