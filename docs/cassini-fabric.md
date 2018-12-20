# Notice: 跨链交易中的转账交易中，不需要转出方签名，只要联盟链签名即可

## RESTFUL Interface
### Fabric provides 4 RESTFUL interfaces:
//The Tx is the QCPTx define in QCP 协议：交易数据结构文档
//the github.com/QOSChain/qbase has go language structure, you could directly use it
```
int GetOutSeq()
int GetInSeq()
Tx GetOutTx(int Sequence)
  //Put the TX result of QOS to QStars
int Publish(Tx)
```
### Cassini provides 1 RESTFUL interfaces:
```
//The Tx is same to above
Notify(Tx) 
```
具体内容见下

# 中继跨链适配hyperledger fabric 区块链

Cassini 中继是基于QCP 协议实现跨链交易的服务中间件。

QCP 协议：交易数据结构文档

https://github.com/QOSGroup/qbase/blob/master/docs/quick_start.md

QCP 协议：交易处理流程文档

https://github.com/QOSGroup/cassini/blob/master/doc/cassini.md

### QOS 区块链技术体系原生QCP 协议实现与中继跨链适配的区别

　　QOS 原生实现：《交易处理流程文档》内容为QOS 区块链技术体系原生QCP 协议实现，接入链与公链均为基于QOS 区块链技
术体系实现。接入链与公链通过中继完成跨链交易的传输。接入链与中继以及公链与中继通过网络链接进行交易事件的订阅，以及交
易数据的查询和传输。

　　中继跨链适配：公链仍是QOS 区块链技术体系，但接入链为非QOS 区块链技术体系实现的区块链。由于网络和部署运营架构的
不可预期，为了减少网络环节可能造成的复杂性以及可能对跨链交易执行效率的影响，中继适配跨链交易时适配模块作为中继的插件
直接内嵌在中继中运行，中继适配模块与中继运行在同一个进程中，中继适配模块通过接入链技术体系的SDK和授权证书，访问接入
链，通过接入链提供的接口或接入链自身的协议标准适配查询获取跨链交易数据，以完成跨链交易。


### 中继适配接入链的基本要求

　　中继适配的目的是为了简化接入链的开发，因此中继适配定义标准化的RPC 接口实现适配。
以下列表为实现QCP 协议最基本的要求，具有以下列表的RPC 接口，才能完成基于QCP 协议进行跨链交易：

1. 查询传出交易序列号（ GetOutSeq ）：

fabric 需要对跨链的传出交易进行编号（sequence）；

编号需要顺序递增（每个交易递增１），不能中断（需要接收到确定的交易执行成功或失败的回传确认交易）；

交易编号查询结果为当前已执行完的交易序号，即：如果还没有执行过交易，查询结果应为0；

2. 查询传入交易序列号（ GetInSeq ）：

fabric 需要记录已执行完成的传入交易的编号（sequence）；

编号同样会顺序递增（每个交易递增１），不能中断（如果交易序号中断，应拒绝该交易）；

交易编号查询结果为当前已执行完的交易序号，即：如果还没有执行过交易，查询结果应为0，此时传入的交易序号应为1，否则应拒绝；

3. 查询传出交易数据（ GetOutTx(Sequence) ）：

fabric 需提供交易数据查询接口，中继适配查询传出交易时，会传入交易序列号；

接口根据交易序列号查询指定交易，并签名（参考：6. 交易签名）返回给中继适配。

4. 提交(传入)交易数据（ Publish(Tx) ）：

传入的交易会自带交易编号（交易序列号/sequence）；

编号会顺序递增（每个交易递增１），不能中断（如果交易序号中断，应拒绝该交易）；

如果还没有执行过交易，查询传入交易序列号结果应为0，此时传入的交易序号应为1，否则应拒绝；

如果传入交易序列号结果为５，则表示在等待交易序列号为６的交易传入，否则应拒绝；

5. 新交易通知机制（ Notify(Event)/RPC回调方法 ）：

中继适配提供交易通知RPC接口，以供接入链调用；

接入链调用通知接口，传送新交易序列号，中继适配接收到通知即可根据序列号查询交易数据；

6. 交易签名（ Sign(Tx) ）：

QOS 公链给fabric 签发证书，并提供交易签名方法提供给fabric-sdk 调用用于签名交易；

中继适配查询交易时，fabric-sdk 将sequence 指定的交易使用QOS 公链提供的证书和方法签名后返回给中继适配；

注：中继适配查询出交易后转交给中继进行共识，中继会在共识完成后使用中继的证书进行签名，然后转发给QOS 公链完成交易。

