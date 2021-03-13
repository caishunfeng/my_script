package main

import (
	"strconv"
	"strings"
	"sync"
	"time"

	log "code.google.com/p/log4go"
	client "github.com/influxdata/influxdb/client/v2"
)

var influxClient client.Client

// var influxBatchPoints client.BatchPoints
var influxdbLock *sync.RWMutex = new(sync.RWMutex)

func InitInfluxdb(adds, user, pass string) error {
	if adds == "" {
		return nil
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     adds,
		Username: user,
		Password: pass,
	})
	if err == nil {
		influxClient = c
	} else {
		log.Error("influxdb 错误:", err)
	}

	// go func() {
	// 	for {
	// 		if influxBatchPoints != nil {
	// 			influxdbLock.Lock()

	// 			oldBp := influxBatchPoints
	// 			influxBatchPoints = nil

	// 			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
	// 				Database:  "interface",
	// 				Precision: "ns",
	// 			})
	// 			if err != nil {
	// 				log.Error("influxdb生成批量节点错误:", err)
	// 				continue
	// 			}
	// 			influxBatchPoints = bp

	// 			influxdbLock.Unlock()

	// 			SendInfluxDB(oldBp)
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	return nil
}

func CreateBatchPoints() (client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "interface",
		Precision: "ns",
	})
	if err != nil {
		log.Error("influxdb 错误:", err)
		return nil, err
	}
	return bp, nil
}

func AddInfluxQueue(bp client.BatchPoints, measurement, customerId, busiType, success, fail, timeStamp string) error {

	successInt, err := strconv.Atoi(success)
	if err != nil {
		log.Error("转换错误:", err)
		return err
	}
	failInt, err := strconv.Atoi(fail)
	if err != nil {
		log.Error("转换错误:", err)
		return err
	}
	//毫秒级别转换为纳秒
	stamp := strings.TrimSpace(timeStamp) + GetRandomString(6)

	stampLong, err := strconv.ParseInt(stamp, 0, 64)
	if err != nil {
		log.Error("转换错误:", err)
		return err
	}

	fields := map[string]interface{}{
		"success": successInt,
		"fail":    failInt,
	}

	tags := map[string]string{
		"customerId": customerId,
		"busiType":   busiType,
	}

	pt, err := client.NewPoint(
		measurement,
		tags,
		fields,
		time.Unix(0, stampLong),
	)

	if err != nil {
		log.Error("influxdb 错误:%v", err)
		return err
	}

	bp.AddPoint(pt)

	return nil
}

func SendInfluxDB(bp client.BatchPoints) error {

	if err := influxClient.Write(bp); err != nil {
		log.Error("批量写入influxdb 错误:%v", err)
		return err
	}

	log.Info("批量写入influxDB成功")

	bp = nil

	return nil
}
