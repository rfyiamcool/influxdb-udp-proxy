package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

const buffMergeSize = 10 * 1024

type influxDB struct {
	url   string
	msgCh chan []byte
}

func createInfluxDB(cfg *InfluxDBConfig) *influxDB {
	db := &influxDB{
		msgCh: make(chan []byte, 10000),
		url:   fmt.Sprintf("http://%s/write?db=%s&u=%s&p=%s", cfg.Addr, cfg.Name, cfg.User, cfg.Pwd),
	}

	go db.process()
	return db
}

func (db *influxDB) process() {
	r, err := regexp.Compile("l1=(mem|cpu|heap)")
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	for {
		buf.Reset()
		loop := true
		et := time.Now().Add(time.Second)

		for loop {
			select {
			case data := <-db.msgCh:
				// 过滤系统监控
				attrs := bytes.Split(data, []byte("\n"))
				if len(attrs) == 0 {
					continue
				} else {
					for _, attr := range attrs {
						if len(attr) == 0 {
							continue
						}
						if r.Match(attr) {
							continue
						} else {
							buf.Write(attr)
							buf.WriteByte('\n')
						}
					}
				}

				if buf.Len() > buffMergeSize {
					loop = false
				}

			case t := <-time.After(time.Millisecond * 100):
				if t.After(et) {
					loop = false
				}
			}

		}

		if buf.Len() == 0 {
			continue
		}

		rsp, err := http.Post(db.url, "application/x-www-form-urlencoded", bytes.NewReader(buf.Bytes()))
		if err != nil {
			log.Warnf("http err:%s", err)
		}
		if rsp == nil {
			continue
		}
		if rsp.StatusCode != 204 {
			robots, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				log.Warnf("http rsp read body err:%s", err)
			}
			log.Infof("set influxdb rsp status: %s, body:%s", rsp.Status, string(robots))
		}

		if rsp.Body != nil {
			rsp.Body.Close()
		}
	}
}

func (db *influxDB) save(buf []byte) {
	select {
	case db.msgCh <- buf:
	default:
		log.Warnf("db msg chan full")
	}
}
