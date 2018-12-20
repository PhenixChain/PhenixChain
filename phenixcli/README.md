# Luanch Client
## Generating keys
```
./phenixcli keys add <key name>
./phenixcli keys add <key name> --recover
./phenixcli keys list
```

## Send transactions
```
./phenixcli tx send --from=<from address> --amount=6mycoin --to=<to address> --chain-id=phenix
```

## Query transaction
```
./phenixcli query tx <TX HASH> --chain-id=phenix
```

## Query account
```
./phenixcli query account <address> --chain-id=phenix
```