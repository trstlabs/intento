# #!/bin/bash
# SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
# source ${SCRIPT_DIR}/../../config.sh

# ### AIRDROP TESTING FLOW Pt 1 (STANDARD)

# # This script tests airdrop claiming on Intento

# # To run:
# #   1. Start the network with `make start-dockernet build=tgr`
# #   2. Run this script with `bash dockernet/scripts/airdrop/airdrop1_standard.sh`

# # NOTE: First, store the keys using the following mnemonics
# echo "Registering accounts..."
# # distributor address: into1z835j3j65nqr6ng257q0xkkc9gta72gfl340ws
# # distributor mnemonic: barrel salmon half click confirm crunch sense defy salute process cart fiscal sport clump weasel render private manage picture spell wreck hill frozen before
# echo "barrel salmon half click confirm crunch sense defy salute process cart fiscal sport clump weasel render private manage picture spell wreck hill frozen before" | \
#     $INTO_MAIN_CMD keys add distributor-test --recover

# # airdrop-test address: into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv
# # airdrop claimer mnemonic: royal auction state december october hip monster hotel south help bulk supreme history give deliver pigeon license gold carpet rabbit raw wool fatigue donate
# echo "royal auction state december october hip monster hotel south help bulk supreme history give deliver pigeon license gold carpet rabbit raw wool fatigue donate" | \
#     $INTO_MAIN_CMD keys add airdrop-test --recover

# ## AIRDROP SETUP
# echo "Funding accounts..."
# # Transfer uatom from gaia to stride, so that we can liquid stake later
# $GAIA_MAIN_CMD tx bank send into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv 1000000uatom --from ${GAIA_VAL_PREFIX}1 -y | TRIM_TX
# sleep 5
# # Fund the distributor account
# # $INTO_MAIN_CMD tx bank send val1 into1z835j3j65nqr6ng257q0xkkc9gta72gfl340ws 600000uinto --from val1 -y | TRIM_TX
# # sleep 5
# # Fund the airdrop account
# $INTO_MAIN_CMD tx bank send val1 into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv 1000000000uinto --from val1 -y | TRIM_TX
# sleep 5
# # Create the airdrop, so that the airdrop account can claim tokens
# # $INTO_MAIN_CMD tx claim create-airdrop gaia GAIA uinto 1679715340 40000000 false --from distributor-test -y | TRIM_TX
# # sleep 5
# # Set airdrop allocations
# # $INTO_MAIN_CMD tx claim set-airdrop-allocations gaia into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv 1 --from distributor-test -y | TRIM_TX
# # sleep 5

# # AIRDROP CLAIMS
# # Check balances before claims
# echo -e "\nInitial balance before claim [1000000000uinto expected]:"
# $INTO_MAIN_CMD query bank balances into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv --denom uinto
# # NOTE: You can claim here using the CLI, or from the frontend!
# # Claim 20% of the free tokens
# echo -e "\nClaiming free amount..."
# $INTO_MAIN_CMD tx claim claim-claimable --from airdrop-test --gas 400000 -y | TRIM_TX
# sleep 5
# echo -e "\nBalance after claim [1000002000uinto expected]:" 
# $INTO_MAIN_CMD query bank balances into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv --denom uinto

# # Stake, to claim another 20%
# echo -e "\nStaking..."
# $INTO_MAIN_CMD tx staking delegate stridevaloper1nnurja9zt97huqvsfuartetyjx63tc5zrj5x9f 100uinto --from airdrop-test --gas 400000 -y | TRIM_TX
# sleep 5
# echo -e "\nBalance after stake [1000239900uinto expected]:" 
# $INTO_MAIN_CMD query bank balances into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv --denom uinto

# # AutoTx, to claim 60% of claimable tokens
# echo -e "\Submitting AutoTx..."
# $INTO_MAIN_CMD tx autoibctx submit-auto-tx 1000 uatom --from airdrop-test --gas 400000 -y | TRIM_TX
# sleep 5
# echo -e "\nBalance after submit auto tx [1000599900uinto expected]:" 
# $INTO_MAIN_CMD query bank balances into1nf6v2paty9m22l3ecm7dpakq2c92ueyue2d5yv --denom uinto

