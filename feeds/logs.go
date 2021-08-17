package feeds

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type EthClient interface {
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	Close()
}

type logFeed struct {
	ctx    context.Context
	wssUrl string
	query  ethereum.FilterQuery
}

func (l *logFeed) ForEachLog(handler func(logEntry types.Log) error) error {
	logs := make(chan types.Log)
	client, err := ethclient.Dial(l.wssUrl)
	if err != nil {
		return err
	}

	sub, err := client.SubscribeFilterLogs(l.ctx, l.query, logs)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(l.ctx)
	ticker := time.NewTicker(500 * time.Millisecond)

	eg.Go(func() error {
		defer client.Close()
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-sub.Err():
				return err
			}
		}
	})

	eg.Go(func() error {
		for ethLog := range logs {
			log.Infof("received an event from agent registry (tx=%s)", ethLog.TxHash.Hex())
			if err := handler(ethLog); err != nil {
				return err
			}
		}
		return nil
	})
	log.Info("subscribed to agent registry updates")
	defer func() {
		log.Warn("stopped subscribing to agent registry updates")
	}()
	return eg.Wait()
}

func NewLogFeed(ctx context.Context, wssUrl string, contractAddrs []string) (*logFeed, error) {
	addrs := make([]common.Address, 0, len(contractAddrs))
	for _, addr := range contractAddrs {
		addrs = append(addrs, common.HexToAddress(addr))
	}

	return &logFeed{
		ctx:    ctx,
		wssUrl: wssUrl,
		query: ethereum.FilterQuery{
			Addresses: addrs,
		},
	}, nil
}