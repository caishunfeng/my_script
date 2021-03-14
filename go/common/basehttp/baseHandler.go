package basehttp

import (
	"bytes"
	log "code.google.com/p/log4go"
	"csf.com/sunshine/client/models/monitor"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

type Filter interface {
	BeforeServeHTTP(http.ResponseWriter, *http.Request) bool
	AfterServeHTTP(http.ResponseWriter, *http.Request) bool
}

type BaseHandler struct {
	Filters []Filter
	Handle  func(r *http.Request) (string, error)
}

// ServeHTTP http处理入口
func (b BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		rerr := recover()
		if rerr == nil {
			log.Info(r.URL.Path)
			return
		}

		var stack string
		var buf bytes.Buffer
		buf.Write(debug.Stack())
		stack = buf.String()
		w.WriteHeader(403)

		log.Error(r.URL.Path, rerr, stack)
		monitor.RecodeError(fmt.Sprintf("%s", rerr), stack)
		return
	}()

	startTime := time.Now()
	defer func() {
		useTime := time.Since(startTime)
		monitor.RecordUseTime(useTime.Nanoseconds() / 1000)
	}()

	monitor.AddQPS()

	if b.Filters != nil {
		for _, filter := range b.Filters {
			if filter != nil && filter.BeforeServeHTTP(w, r) == false {
				return
			}
		}
	}

	// 业务逻辑处理
	result, err := b.Handle(r)
	if err != nil {
		log.Error(err)
		return
	}

	if b.Filters != nil {
		for _, filter := range b.Filters {
			if filter != nil && filter.AfterServeHTTP(w, r) == false {
				return
			}
		}
	}

	w.Write([]byte(result))
}
