package influxdb

import (
	"sync"
	"time"

	"errors"
	"github.com/influxdata/influxdb/client/v2"
)

type InfluxCfg struct {
	adds string
	user string
	pass string
}

type InfluxClient struct {
	client client.Client
	bp     client.BatchPoints
	lock   *sync.RWMutex
}

func InitInfluxClient(cfg *InfluxCfg) (influxClient *InfluxClient, err error) {
	if cfg.adds == "" {
		return nil, errors.New("influxdb connect address is nil")
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.adds,
		Username: cfg.user,
		Password: cfg.pass,
	})
	if err != nil {
		return nil, err
	}
	influxClient = &InfluxClient{
		client: c,
	}

	return
}

func (this *InfluxClient) CreateBatchPoints(dataBase, precision string) (err error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dataBase,
		Precision: precision,
	})
	if err != nil {
		return err
	}
	this.lock.Lock()
	this.bp = bp
	this.lock.Unlock()
	return
}

func (this *InfluxClient) AddInfluxQueue(measurement, tagName, tagValue, fieldName string, fieldValue float64, ts time.Time) error {
	fields := map[string]interface{}{
		fieldName: fieldValue,
	}

	tags := map[string]string{
		tagName: tagValue,
	}

	pt, err := client.NewPoint(
		measurement,
		tags,
		fields,
		ts,
	)

	if err != nil {
		return err
	}

	this.bp.AddPoint(pt)

	return nil
}

func (this *InfluxClient) SendInfluxDB(bp client.BatchPoints) (err error) {
	if err := this.client.Write(bp); err != nil {
		return err
	}

	this.lock.Lock()
	bp = nil
	this.lock.Unlock()
	return
}
