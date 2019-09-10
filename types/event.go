package types

import (
	"errors"
	"strings"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/tendermint/tendermint/libs/common"
)

// CassiniEventDataTx holds tx event tags
type CassiniEventDataTx struct {
	From      string `json:"from"` //qsc name 或 qos
	To        string `json:"to"`   //qsc name 或 qos
	Height    int64  `json:"height"`
	Sequence  int64  `json:"sequence"`
	HashBytes []byte `json:"hashBytes"` //TxQcp 做 sha256
}

// Event cache tx event tags and node info
type Event struct {
	NodeAddress        string               `json:"node"` //event 源地址
	CassiniEventDataTx `json:"eventDataTx"` //event 事件
}

// ConstructFromTags parse tx event tags
func (c *CassiniEventDataTx) ConstructFromTags(tags map[string]string) (err error) {

	if tags == nil || len(tags) == 0 {
		err = errors.New("empty tags")
		return
	}
	for key, val := range tags {
		log.Debug("event.tag: ", key, "; ", val)
		if strings.EqualFold(key, "tx.height") {
			c.Height, err = ParseHeight(val)
			if err != nil {
				return err
			}
		}
		if strings.EqualFold(key, qcp.From) {
			c.From = val
		}
		if strings.EqualFold(key, qcp.To) {
			c.To = val
		}
		if strings.EqualFold(key, qcp.Hash) {
			c.HashBytes = []byte(val)
		}
		if strings.EqualFold(key, qcp.Sequence) {
			c.Sequence, err = ParseSequence(val)
			if err != nil {
				return err
			}
		}
	}

	return
}

// KV2map returns map
func KV2map(kvs []common.KVPair) (
	tags map[string]string, err error) {
	tags = make(map[string]string)
	if kvs == nil || len(kvs) == 0 {
		return tags, errors.New("empty tags")
	}
	for _, tag := range kvs {
		if strings.EqualFold(string(tag.Key), qcp.From) {
			tags[qcp.From] = string(tag.Value)
		}
		if strings.EqualFold(string(tag.Key), qcp.To) {
			tags[qcp.To] = string(tag.Value)
		}
		if strings.EqualFold(string(tag.Key), qcp.Hash) {
			tags[qcp.Hash] = string(tag.Value)
		}
		if strings.EqualFold(string(tag.Key), qcp.Sequence) {
			tags[qcp.Sequence] = string(tag.Value)
		}
	}

	return
}
