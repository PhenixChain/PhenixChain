# Luanch Blockchain
## Generate genesis file
```
./phenix init
```
Add both accounts, with coins to the genesis file
```
./phenix add-genesis-account <address> 10000000mycoin,666666coin1
```
## Start up the blockchain
```
./phenix start
```
## Reset the blockchain data
```
./phenix unsafe-reset-all
```