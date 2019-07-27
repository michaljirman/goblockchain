### goblockchain

Basic implementation of the blockchain / cryptocurrency in the Golang (PoC, WIP)


### Examples

#### Example 1
1. create two wallets so you can send and receive

```
DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go createwallet
```

```
{"level":"debug","time":"2019-07-27T15:26:36+01:00","message":"pub key: 5af839df8776d16621994bf8b8b5f61c5faa5a2eaed6079fa921333ed69ecfa4554d6d9d5aba260caebdefc1ba51d4b5348cc95bea9cb2152e3878b78e8e1bcf"}
{"level":"debug","time":"2019-07-27T15:26:36+01:00","message":"pub hash: cb68331e111d5baea6aed4feff8c8d6d32f67064"}
{"level":"debug","time":"2019-07-27T15:26:36+01:00","message":"address: 1KYWzs15ziCSkbcgduTyPN79YZD9du4ASd"}
New address is: 1KYWzs15ziCSkbcgduTyPN79YZD9du4ASd
```

```
DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go createwallet
```
```
{"level":"debug","time":"2019-07-27T15:26:52+01:00","message":"pub key: d9a0043ca8fb62371a65398d19d5203447f61ffc8ca0496fe72bb27b417fb0cdb72abae093d14ef0057d398ce419c48f5feacc437dc9e665b7d6301dffdababe"}
{"level":"debug","time":"2019-07-27T15:26:52+01:00","message":"pub hash: 601ca2cfca145260edcc0ed4078074d3df5a0b93"}
{"level":"debug","time":"2019-07-27T15:26:52+01:00","message":"address: 19mCB5ZmuyhMESRZt7KwCggoEupby8DZo8"}
New address is: 19mCB5ZmuyhMESRZt7KwCggoEupby8DZo8
```


2. create a blockchain with one of the new addresses
```
DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go createblockchain -address 1KYWzs15ziCSkbcgduTyPN79YZD9du4ASd
```

```
badger 2019/07/27 15:31:10 INFO: All 0 tables opened in 0s
00000d8068fc517cbe34a443e2b7cb909c2c0f59cf2aef0053de85e402701203
{"level":"debug","time":"2019-07-27T15:31:11+01:00","message":"Genesis created"}
badger 2019/07/27 15:31:11 DEBUG: Storing value log head: {Fid:0 Len:42 Offset:631}
badger 2019/07/27 15:31:11 INFO: Got compaction priority: {level:0 score:1.73 dropPrefix:[]}
badger 2019/07/27 15:31:11 INFO: Running for level: 0
badger 2019/07/27 15:31:11 DEBUG: LOG Compact. Added 3 keys. Skipped 0 keys. Iteration took: 100.181µs
badger 2019/07/27 15:31:11 DEBUG: Discard stats: map[]
badger 2019/07/27 15:31:11 INFO: LOG Compact 0->1, del 1 tables, add 1 tables, took 14.788318ms
badger 2019/07/27 15:31:11 INFO: Compaction for level: 0 DONE
badger 2019/07/27 15:31:11 INFO: Force compaction on level 0 done
Finished!
```

3. print blockchain, it will contains only one block, the Genesis block
``` 
DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go printchain
```

```
badger 2019/07/27 15:31:15 INFO: All 1 tables opened in 0s
badger 2019/07/27 15:31:15 INFO: Replaying file id: 0 at offset: 673
badger 2019/07/27 15:31:15 INFO: Replay took: 27.826µs
Prev. hash:
Hash: 00000d8068fc517cbe34a443e2b7cb909c2c0f59cf2aef0053de85e402701203
PoW: true
--- Transaction 67ed82db3f222b9289e0bac60447189cbce6e9ceef06f536e775bc25fc44a1d2:
Input 0:
TXID:
Out:       -1
Signature:
PubKey:    4669727374205472616e73616374696f6e2066726f6d2047656e65736973
Output 0:
Value:  100
Script: cb68331e111d5baea6aed4feff8c8d6d32f67064

badger 2019/07/27 15:31:15 INFO: Got compaction priority: {level:0 score:1.73 dropPrefix:[]}
```

4. sends tokens from the first to the second wallet 
```
DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go send -from 1KYWzs15ziCSkbcgduTyPN79YZD9du4ASd -to 19mCB5ZmuyhMESRZt7KwCggoEupby8DZo8 -amount 30
```

```
badger 2019/07/27 15:34:12 INFO: All 1 tables opened in 0s
badger 2019/07/27 15:34:12 INFO: Replaying file id: 0 at offset: 673
badger 2019/07/27 15:34:12 INFO: Replay took: 83.756µs
00002b8800949692ee49ff507c1c0b1b390aea13aa73a8cd18f29537e8b522a9
Success!
badger 2019/07/27 15:34:14 DEBUG: Storing value log head: {Fid:0 Len:42 Offset:1496}
badger 2019/07/27 15:34:14 INFO: Got compaction priority: {level:0 score:1.73 dropPrefix:[]}
badger 2019/07/27 15:34:14 INFO: Running for level: 0
badger 2019/07/27 15:34:14 DEBUG: LOG Compact. Added 6 keys. Skipped 0 keys. Iteration took: 275.606µs
badger 2019/07/27 15:34:14 DEBUG: Discard stats: map[]
badger 2019/07/27 15:34:14 INFO: LOG Compact 0->1, del 2 tables, add 1 tables, took 16.555894ms
badger 2019/07/27 15:34:14 INFO: Compaction for level: 0 DONE
badger 2019/07/27 15:34:14 INFO: Force compaction on level 0 done
```

5. print blockchain, it will contains the genesis block and the second block with our new transaction
```
DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go printchain
```

```  
badger 2019/07/27 15:35:26 INFO: All 1 tables opened in 0s
badger 2019/07/27 15:35:26 INFO: Replaying file id: 0 at offset: 1538
badger 2019/07/27 15:35:26 INFO: Replay took: 37.609µs
Prev. hash: 00000d8068fc517cbe34a443e2b7cb909c2c0f59cf2aef0053de85e402701203
Hash: 00002b8800949692ee49ff507c1c0b1b390aea13aa73a8cd18f29537e8b522a9
PoW: true
--- Transaction 3497be7dd0429cd80304b2ef314d1b653fff1e4446f536116a70ef022b31c26e:
Input 0:
TXID:     67ed82db3f222b9289e0bac60447189cbce6e9ceef06f536e775bc25fc44a1d2
Out:       0
Signature: ade91cf2107c9e76312912ffb7abd57204b3f3a9480b0abdd9d9b9b226f2a94be44d003c163778e2699549932fbb762ec8bab82ad317b930ce7160cb52c5fa5d
PubKey:    5af839df8776d16621994bf8b8b5f61c5faa5a2eaed6079fa921333ed69ecfa4554d6d9d5aba260caebdefc1ba51d4b5348cc95bea9cb2152e3878b78e8e1bcf
Output 0:
Value:  30
Script: 601ca2cfca145260edcc0ed4078074d3df5a0b93
Output 1:
Value:  70
Script: cb68331e111d5baea6aed4feff8c8d6d32f67064

Prev. hash:
Hash: 00000d8068fc517cbe34a443e2b7cb909c2c0f59cf2aef0053de85e402701203
PoW: true
--- Transaction 67ed82db3f222b9289e0bac60447189cbce6e9ceef06f536e775bc25fc44a1d2:
Input 0:
TXID:
Out:       -1
Signature:
PubKey:    4669727374205472616e73616374696f6e2066726f6d2047656e65736973
Output 0:
Value:  100
Script: cb68331e111d5baea6aed4feff8c8d6d32f67064

badger 2019/07/27 15:35:26 INFO: Got compaction priority: {level:0 score:1.73 dropPrefix:[]}
```

#### References

Based on the Tensor's [tutorial](https://github.com/tensor-programming/golang-blockchain)\
Logging library [zerolog](https://github.com/rs/zerolog)\
Key-value DB [badger](https://github.com/dgraph-io/badger)