# Luanch Client
## Generating keys
```
./phenixcli keys add <key name>
./phenixcli keys add <key name> --recover
./phenixcli keys list
```

## Send transactions
```
./phenixcli send --from=<from address> --amount=6coin1 --to=<to address> --chain-id=phenix
```

## Query transaction
```
./phenixcli tx <TX HASH> --chain-id=phenix
```

## Query account
```
./phenixcli account <address> --chain-id=phenix
```