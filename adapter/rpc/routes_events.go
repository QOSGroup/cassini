package rpc

// copy from tendermint/rpc/core/events.go

import (
	"context"

	"github.com/QOSGroup/cassini/log"
	"github.com/pkg/errors"

	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Subscribe 指定订阅条件，订阅交易事件。
func (s DefaultHandlerService) Subscribe(wsCtx rpctypes.WSRPCContext, query string) (*ctypes.ResultSubscribe, error) {
	addr := wsCtx.GetRemoteAddr()
	log.Info("Subscribe to query", "remote", addr, "query", query)

	q, err := tmquery.New(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse query")
	}

	ctx, cancel := context.WithTimeout(context.Background(), subscribeTimeout)
	defer cancel()
	ch := make(chan interface{})
	err = s.eventBusFor(wsCtx).Subscribe(ctx, addr, q, ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for event := range ch {
			tmResult := &ctypes.ResultEvent{Query: query, Data: event.(tmtypes.TMEventData)}
			wsCtx.TryWriteRPCResponse(rpctypes.NewRPCSuccessResponse(wsCtx.Codec(), wsCtx.Request.ID+"#event", tmResult))
		}
	}()

	return &ctypes.ResultSubscribe{}, nil
}

// Unsubscribe 根据具体订阅条件，取消交易事件的订阅。
func (s DefaultHandlerService) Unsubscribe(wsCtx rpctypes.WSRPCContext, query string) (*ctypes.ResultUnsubscribe, error) {
	addr := wsCtx.GetRemoteAddr()
	log.Info("Unsubscribe from query", "remote", addr, "query", query)
	q, err := tmquery.New(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse query")
	}
	err = s.eventBusFor(wsCtx).Unsubscribe(context.Background(), addr, q)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultUnsubscribe{}, nil
}

// UnsubscribeAll 取消全部交易事件订阅。
func (s DefaultHandlerService) UnsubscribeAll(wsCtx rpctypes.WSRPCContext) (*ctypes.ResultUnsubscribe, error) {
	addr := wsCtx.GetRemoteAddr()
	log.Info("Unsubscribe from all", "remote", addr)
	err := s.eventBusFor(wsCtx).UnsubscribeAll(context.Background(), addr)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultUnsubscribe{}, nil
}

func (s DefaultHandlerService) eventBusFor(wsCtx rpctypes.WSRPCContext) tmtypes.EventBusSubscriber {
	es := wsCtx.GetEventSubscriber()
	if es == nil {
		es = s.eventHub
	}
	return es
}
