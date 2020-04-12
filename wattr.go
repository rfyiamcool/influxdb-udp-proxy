package main

import (
	"bytes"
	"context"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

func createAttrWorker(ctx context.Context, wg *sync.WaitGroup, Conn *net.UDPConn) *attrWorker {
	aw := attrWorker{conn: Conn, wg: wg, Context: ctx}
	go aw.loop()

	return &aw
}

type attrWorker struct {
	conn *net.UDPConn
	wg   *sync.WaitGroup

	context.Context
}

func (aw *attrWorker) loop() {
	aw.wg.Add(1)
	defer aw.wg.Done()

	var buf [64 * 1024]byte
	for {
		if aw.Err() == nil {
			return
		}

		l, addr, err := aw.conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Warnf("Read from UDP port error %s", err)
			continue
		}
		data := bytes.Replace(buf[:l], []byte("$IP"), []byte(addr.IP.String()), 1)

		db.save(data)
	}
}
