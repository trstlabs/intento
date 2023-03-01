# connections for ICA
rly tx connection trstdev1-trstdev2 --override
rly tx connection trstdev1-localjuno --override
# transfer channel 
rly tx channel trstdev1-trstdev2 --override
rly tx channel trstdev1-localjuno --override
rly start trstdev1-trstdev2 trstdev1-localjuno -p events -b 100 --debug