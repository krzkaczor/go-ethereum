## OFAC Geth

Ethereum-go with censorship built-in!

## Usage

This fork adds 2 new flags:

- `--rpc.blacklist PATH` - use with a path to JSON file containg an array of blacklisted arrays.
- `--rpc.compliance-api-url URL` - use with an [API](https://www.trmlabs.com/) that returns current compliance status for a given account.

### Building and running

```sh
make geth

# run on localnet
./build/bin/geth --rpc.compliance-api-url "X" --rpc.blacklist Y --datadir data --networkid 12345 --http --mine --miner.threads=1 --miner.etherbase=Y
```

## Why?

Just for lulz. Really I was looking for an excuse to dive into geth's codebase.
