DB: 模块选用rocksdb 或者 leveldb 或者IPFS

## rpc api
### Blockchain Api
GetBlockCount
// curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"Client.Say","params":{"Who":"sang"},"id":1}' http://localhost:1234
### Block Api

TODO:实现以太坊类似的interactive console
interactive (console) --> https://github.com/motemen/gore
https://github.com/mkouhei/gosh

Go REPL
