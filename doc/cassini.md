# 中继约定

版本:
v0.1

日期:
2018年09月26日

目录
[TOC]

## 简介：

QCP跨链协议中中继（relay）作为链之间连接的纽带，使跨链交易能够完成。前提是中继和链需要符合本文的约定。

## 体系结构：

![framework](https://github.com/QOSGroup/static/blob/master/relay.png?raw=true)

- BlockchainA和BlockchainB连接共同的中继Relay；
- 中继订阅链上跨链交易事件，获取交易摘要，进行2/3共识，随机找一个诚实节点获取跨链交易；
- 当BlockchainA新块内有跨链交易，将跨链交易放入outbox,并按顺序递增编号，当前最大编号叫sequcnce。区块链保证该编号的连续性；
- BlockchainB将通过中继收到的交易存入inbox；
- 中继查询BlockchainB inbox 的sequence 记作seq1,中继顺序取BlockchainA 的 outbox 的编号大于seq1的交易，一次可以取一条或多条；
- 中继对交易进行验签，2/3共识等处理后路由到目标链；
- 为防止单个中继节点作恶，未来会在中继节点之间进行拜占庭共识。

## Event 事件发布订阅以及Tx 交易查询与处理：

- 1、	中继通过websocket 订阅链上跨链event。过滤条件为 “tm.event = 'Tx' AND qcp.to = 'xxx'”；
- 2、	链将跨链的交易在结构体 ResponseDeliverTx成员变量Tags中增加值对 “qcp.to = 'xxx'“ 。XXX表示目标链的名称；
- 3、	中继收到事件后，会进行2/3 的共识校验；
- 4、	通过2/3 共识校验后，中继会进一步调用restful API（ABCI Query）查询交易数据；
- 5、	中继查询到交易数据，会通过调用restful API（ABCI BroadcastTxAsync 或ABCI BroadcastTxSync）向目标链提交交易，完成交易的处理；
- 6、	跨链交易结果返回过程同1,2,3,4,5步。

### Event数据结构

```
EventDataTx{
    Height
    Index
    Tx
    Result.Data
    Result.Tags    {
        {"qcp.from":     string  }, //qsc name 或 qos
        {"qcp.to":       string  }, //qsc name 或 qos
        {"qcp.sequence": int64   },
        {"qcp.hash":     []byte  },  //TxQcp 做 sha256
        {"qcp.height":   int64   },  //区块高度
    }
}

```

## Rest Service API需求：

- get       取 outbox 最大sequence
- get       取 inbox 最大sequence
- get	    取 outbox 中给定sequence编号的交易（TxQcp）
- post	    接收交易（TxQcp）
- get	    批量取 outbox 中给定sequence编号之后的交易（TxQcp）
- post	    批量接收交易（TxQcp）

### Rest Service 数据结构

```
type TxQcp struct {
	Payload  	[]byte 	//交易体
	From     	string 	//源链ID
	To       	string 	//目标链ID
	Sequence 	int64
	IsResult 	bool   	//是否为result
	Height     	int64   //块高度  
	Sign     	[]byte 	//对本交易的签名，公链不需要签名     
}

```