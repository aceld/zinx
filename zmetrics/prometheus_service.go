package zmetrics

import (
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/zlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

const (
	METRICS_ROUTE string = "/metrics"
)

var metricsServiceOnce sync.Once
var metricsInitOnce sync.Once

func RunMetricsService(conf *zconf.Config) (err error) {

	metricsServiceOnce.Do(func() {
		// metricsServic 只启动一个服务
		go func() {
			http.Handle(METRICS_ROUTE, promhttp.Handler())
			err = http.ListenAndServe(conf.PrometheusListen, nil) //多个进程不可监听同一个端口
			if err != nil {
				zlog.Ins().ErrorF("RunMetricsService err = %s\n", err)
			}
		}()
	})

	zlog.Ins().InfoF("RunMetricsService %s:%s success", METRICS_ROUTE, conf.PrometheusListen)

	return err
}

func InitZinxMetrics() {

	metricsInitOnce.Do(func() {
		Metrics().connTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: GANGEVEC_ZINX_CONNECTION_TOTAL_NAME,
				Help: GANGEVEC_ZINX_CONNECTION_TOTAL_HELP,
			},
			[]string{LABEL_SERVER_ADDRESS, LABEL_SERVER_NAME},
		)

		Metrics().taskTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: GANGEVEC_ZINX_TASK_TOTAL_NAME,
				Help: GANGEVEC_ZINX_TASK_TOTAL_HELP,
			},
			[]string{LABEL_WORKER_ID},
		)

		//Register
		prometheus.MustRegister(Metrics().connTotal)
		prometheus.MustRegister(Metrics().taskTotal)
	})

}
