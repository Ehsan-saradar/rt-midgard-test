package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pascaldekloe/metrics/gostat"

	"gitlab.com/thorchain/midgard/chain"
	"gitlab.com/thorchain/midgard/event"
	"gitlab.com/thorchain/midgard/internal/api"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/timeseries"
	"gitlab.com/thorchain/midgard/internal/timeseries/stat"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	log.Print("daemon launch as ", strings.Join(os.Args, " "))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// include Go runtime metrics
	gostat.CaptureEvery(5 * time.Second)

	// read configuration
	var c config.Configuration
	switch len(os.Args) {
	case 1:
		break // refer to defaults
	case 2:
		// TODO(pascaldekloe): Move configuration to main, as it belongs
		// to the command/invokation/daemon. We don't need multiple format
		// support either. Drop the dependencies.
		// BUG(pascaldekloe): Configuration now silently ignores unknown
		// parameters, which may lead to settings not applied due typos.
		p, err := config.LoadConfiguration(os.Args[1])
		if err != nil {
			log.Fatal("exit on configuration unavailable: ", err)
		}
		c = *p
	default:
		log.Fatal("one optional configuration file argument only—no flags")
	}

	// apply configuration
	SetupDatabase(c.TimeScale)
	blocks := SetupBlockchain(c.ThorChain)
	if c.ListenPort == 0 {
		c.ListenPort = 8080
		log.Printf("default HTTP server listen port to %d", c.ListenPort)
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 2 * time.Second
		log.Printf("default HTTP server read timeout to %s", c.ReadTimeout)
	}
	if c.WriteTimeout == 0 {
		c.ReadTimeout = 2 * time.Second
		log.Printf("default HTTP server write timeout to %s", c.ReadTimeout)
	}
	srv := &http.Server{
		Handler:      api.CORS(api.Handler),
		Addr:         fmt.Sprintf(":%d", c.ListenPort),
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}

	// launch HTTP server
	go func() {
		err := srv.ListenAndServe()
		log.Print("HTTP stopped on ", err)
		signals <- syscall.SIGABRT
	}()

	// launch blockchain reading
	go func() {
		m := event.Demux{Listener: timeseries.EventListener}
		for block := range blocks {
			m.Block(block)
			err := timeseries.CommitBlock(block.Height, block.Time, block.Hash)
			if err != nil {
				log.Print("timeseries feed stopped on ", err)
				signals <- syscall.SIGABRT
				return
			}
		}
		log.Print("timeseries feed stopped")
		signals <- syscall.SIGABRT
	}()

	signal := <-signals
	log.Print("HTTP shutdown initiated with timeout in ", c.ShutdownTimeout)
	ctx, _ := context.WithTimeout(context.Background(), c.ShutdownTimeout)
	if err := srv.Shutdown(ctx); err != nil {
		log.Print("HTTP shutdown: ", err)
	}

	log.Fatal("exit on signal ", signal)
}

func SetupDatabase(c config.TimeScaleConfiguration) {
	db, err := sql.Open("pgx", fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%d", c.UserName, c.Database, c.Sslmode, c.Password, c.Host, c.Port))
	if err != nil {
		log.Fatal("exit on PostgreSQL client instantiation: ", err)
	}

	stat.DBQuery = db.Query
	timeseries.DBExec = db.Exec
	timeseries.DBQuery = db.Query
}

// SetupBlockchain launches the synchronisation routine.
func SetupBlockchain(c config.ThorChainConfiguration) <-chan chain.Block {
	// normalize configuration
	if c.Scheme == "" {
		c.Scheme = "http"
		log.Printf("default Tendermint RPC scheme to %q", c.Scheme)
	}
	if c.Host == "" {
		c.Host = "localhost:26657"
		log.Print("default Tendermint RPC host to %q", c.Host)
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 2 * time.Second
		log.Print("default Tendermint read timeout to %s", c.ReadTimeout)
	}
	// BUG(pascaldekloe): NoEventsBackoff is a misnommer
	// as chains without events do not cause any backoff.
	if c.NoEventsBackoff == 0 {
		c.NoEventsBackoff = 5 * time.Second
		log.Printf("default Tendermint no events backoff to %s", c.NoEventsBackoff)
	}

	// instantiate client
	endpoint := &url.URL{Scheme: c.Scheme, Host: c.RPCHost, Path: "/websocket"}
	log.Print("Tendermint enpoint set to ", endpoint.String())
	client, err := chain.NewClient(endpoint, c.ReadTimeout)
	if err != nil {
		// error check does not include network connectivity
		log.Fatal("exit on Tendermint RPC client instantiation: ", err)
	}

	// fetch current position (from commit log)
	offset, _, _, err := timeseries.Setup()
	if err != nil {
		// no point in running without a database
		log.Fatal("exit on RDB unavailable: ", err)
	}
	if offset != 0 {
		offset++
		log.Print("starting with previous blockchain height ", offset)
	}

	var lastNoData atomic.Value
	api.InSync = func() bool {
		return time.Since(lastNoData.Load().(time.Time)) < 2*c.NoEventsBackoff
	}

	// launch read routine
	ch := make(chan chain.Block, 99)
	go func() {
		backoff := time.NewTicker(c.NoEventsBackoff)
		defer backoff.Stop()

		// TODO(pascaldekloe): Could use a limited number of
		// retries with skip block logic perhaps?
		for {
			offset, err = client.Follow(ch, offset, nil)
			switch err {
			case chain.ErrNoData:
				lastNoData.Store(time.Now())
			default:
				log.Print("follow blockchain retry on ", err)
			}
			<-backoff.C
		}
	}()

	return ch
}
