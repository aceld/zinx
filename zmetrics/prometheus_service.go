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
		// metricsService only starts one server. (metricsServic 只启动一个服务)
		go func() {
			http.Handle(METRICS_ROUTE, promhttp.Handler())
			// Multiple processes cannot listen on the same port.
			// (多个进程不可监听同一个端口)
			err = http.ListenAndServe(conf.PrometheusListen, nil)
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
			[]string{LABEL_ADDRESS, LABEL_NAME},
		)

		Metrics().taskTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: GANGEVEC_ZINX_TASK_TOTAL_NAME,
				Help: GANGEVEC_ZINX_TASK_TOTAL_HELP,
			},
			[]string{LABEL_ADDRESS, LABEL_NAME, LABEL_WORKER_ID},
		)

		Metrics().routerScheduleTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: GANGEVEC_ZINX_ROUTER_SCHEDULE_TOTAL_NAME,
				Help: GANGEVEC_ZINX_ROUTER_SCHEDULE_TOTAL_HELP,
			},
			[]string{LABEL_ADDRESS, LABEL_NAME, LABEL_WORKER_ID, LABEL_MSG_ID},
		)

		Metrics().routerScheduleDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    HISTOGRAM_ZINX_ROUTER_SCHEDULE_DURATION_NAME,
				Help:    HISTOGRAM_ZINX_ROUTER_SCHEDULE_DURATION_HELP,
				Buckets: []float64{0.005, 0.01, 0.03, 0.08, 0.1, 0.5, 1.0, 5.0, 10, 100, 1000, 5000, 30000}, //单位ms,最大半分钟
			},
			[]string{LABEL_ADDRESS, LABEL_NAME, LABEL_WORKER_ID, LABEL_MSG_ID},
		)

		//Register
		prometheus.MustRegister(Metrics().connTotal)
		prometheus.MustRegister(Metrics().taskTotal)
		prometheus.MustRegister(Metrics().routerScheduleTotal)
		prometheus.MustRegister(Metrics().routerScheduleDuration)
	})

}
