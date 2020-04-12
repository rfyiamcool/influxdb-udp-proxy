package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	db *influxDB
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

func main() {
	// init cfg
	initConfig()

	// process
	db = createInfluxDB(cfg.DB)
	conn := openSocket(cfg.Port)
	createAttrWorker(ctx, &wg, conn)

	// graceful exit
	waitStop()
}

func waitStop() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for {
		sig := <-sc
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT:
			cancel()
			wg.Wait()
			log.Info("exit: bye :-).")
			os.Exit(0)
		default:
			continue
		}
	}
}

func openSocket(port int) *net.UDPConn {
	udpConn, err := net.ListenUDP("udp",
		&net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: port},
	)

	if err != nil {
		log.Errorf("Open UDP port error %s:", err)
		os.Exit(-1)
	}

	err = udpConn.SetReadBuffer(64 * 1024)
	if err != nil {
		log.Errorf("Change read buffer size error %s", err)
		os.Exit(-1)
	}

	log.Infof("udp proxy start...")
	return udpConn
}
