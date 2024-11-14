package load

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/p2p/conn"
)

func TestMultipleConnections(t *testing.T) {
	cfg := config.DefaultP2PConfig()
	cfg.AllowDuplicateIP = true
	cfg.DialTimeout = 10 * time.Second
	mcfg := conn.DefaultMConnConfig()
	mcfg.SendRate = 5000000
	mcfg.RecvRate = 5000000
	mcfg.FlushThrottle = 100 * time.Millisecond

	peerCount := 20
	reactors := make([]*MockReactor, peerCount)
	nodes := make([]*node, peerCount)

	chainID := "base-30"

	for i := 0; i < peerCount; i++ {
		reactor := NewMockReactor(defaultTestChannels, defaultMsgSizes)
		node, err := newnode(*cfg, mcfg, chainID, reactor)
		require.NoError(t, err)

		err = node.start()
		require.NoError(t, err)
		defer node.stop()

		reactors[i] = reactor
		nodes[i] = node
		fmt.Println("added node", i, node.addr)
	}

	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	for i := 1; i < peerCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i, nodes[i].addr)
			err := nodes[0].sw.DialPeerWithAddress(nodes[i].addr)
			require.NoError(t, err)
		}(i)

	}

	wg.Wait()

	go func() {
		for i := 0; i < 90; i++ {
			for _, reactor := range reactors {
				reactor.PrintReceiveSpeed()
			}
			time.Sleep(5 * time.Second)
		}
	}()

	for _, reactor := range reactors {
		reactor.FloodAllPeers(&wg, time.Second*30,
			FirstChannel,
			SecondChannel,
			ThirdChannel,
			FourthChannel,
			FifthChannel,
			SixthChannel,
			SeventhChannel,
			EighthChannel,
			NinthChannel,
			TenthChannel,
		)
	}

	wg.Wait()

	reactors[0]
}
