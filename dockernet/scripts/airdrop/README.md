# Airdrop Integration Tests

Each airdrop testing script (1 through 4) tests different aspects of the airdrop

## Overview

* **Part 1: Standard**: Tests basic airdrop claims and actions
* **Part 2: Vesting**: Tests that the airdrop vests properly

### Instructions

* Only the GAIA host zone is required. Start dockernet with:

```bash
make start-dockernet build=tgr
```

* Run the corresponding script

```bash
bash dockernet/scripts/airdrop/airdrop{1/2/3/4}.sh
```

* **NOTE**: Each script must be run independently, meaning you must restart dockernet between runs (`make start-dockernet build=tgr`)
