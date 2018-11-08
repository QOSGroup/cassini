package types

import (
	"errors"

	"github.com/tendermint/tendermint/libs/common"
)

type CassiniEventDataTx struct {
	From      string `json:"from"` //qsc name 或 qos
	To        string `json:"to"`   //qsc name 或 qos
	Height    int64  `json:"height"`
	Sequence  int64  `json:"sequence"`
	HashBytes []byte `json:"hashBytes"` //TxQcp 做 sha256
}

type Event struct {
	NodeAddress        string               `json:"node"` //event 源地址
	CassiniEventDataTx `json:"eventDataTx"` //event 事件
}

func (c *CassiniEventDataTx) ConstructFromTags(tags []common.KVPair) (err error) {

	if tags == nil || len(tags) == 0 {
		return errors.New("empty tags")
	}
	for _, tag := range tags {
		if string(tag.Key) == "qcp.from" {
			c.From = string(tag.Value)
		}
		if string(tag.Key) == "qcp.to" {
			c.To = string(tag.Value)
		}
		if string(tag.Key) == "qcp.hash" {
			c.HashBytes = tag.Value
		}
		if string(tag.Key) == "qcp.sequence" {
			//c.Sequence, err = BytesInt64(tag.Value)
			//c.Sequence, err = strconv.ParseInt(string(tag.Value), 10, 64)
			c.Sequence, err = ParseSequence(tag.Value)
			if err != nil {
				return err
			}
		}
	}

	return
}
