## 操作文档

1. 第一步 创建并写入创世区块: 创世区块硬编码写死到代码中,钱包写死在目录data/wallet.dat_bck 中
 `
 curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.GenesisBlock","params":[],"id":1}' http://localhost:8080
 `
 2. 钱包相关命令: 钱包文件会保存在./data/wallet.dat文件中
*  生成钱包
 `
curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.CreateWallet","params":{},"id":1}' http://localhost:8080
`
* 显示所有本地钱包缓存
`
curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.ListWallets","params":{},"id":1}' http://localhost:8080
`

3. 区块相关命令
* 获取当前区块高度
`
 curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.GetBlockCount","params":{},"id":1}' http://localhost:8080
`

* 根据区块高度获取区块详情: 传参区块高度
`
 curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.GetBlockByNumber","params":{1},"id":1}' http://localhost:8080
`

* 根据区块hash获取区块详情:传参区块hash
`
 curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.GetBlockByHash","params":{"${BLOCKHASH}"},"id":1}' http://localhost:8080
`
* 获取最新确认的区块详情: 无参
`
 curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.LastBlock","params":{"${BLOCKHASH}"},"id":1}' http://localhost:8080
`

3 交易相关操作
* 发起一笔交易: 传入接收方和发送方地址(是钱包地址) 已经额度
`
 curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.SendTransaction","params":{"From":"1AeGtuczZ6aHoZRkWWHBWpUjeY3HxAe5ie","To": "1AeGtuczZ6aHoZRkWWHBWpUjeY3HxAe5ie","Value":2},"id":1}' http://localhost:8080
`
* 

4 账户相关操作
* 显示地址账户详情(UTXO): 参数: 钱包地址
`
curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.ListWallets","params":{},"id":1}' http://localhost:8080
`
* 显示账户余额: 参数: 钱包地址 
`
 curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.GetBalance","params":["$1"],"id":1}' http://localhost:8080
`
