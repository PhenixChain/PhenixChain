# Luanch Client
## Generating keys
```
./phenixcli keys add <key name>
./phenixcli keys add <key name> --recover
./phenixcli keys list
```

## Send transactions
```
./phenixcli tx send <to address> 6mycoin --from=<from address> --chain-id=phenix
```

## Query transaction
```
./phenixcli query tx <TX HASH> --chain-id=phenix
```

## Query account
```
./phenixcli query account <address> --chain-id=phenix
```