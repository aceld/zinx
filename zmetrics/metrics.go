package zmetrics

import (
	"github.com/aceld/zinx/zconf"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

var _metrics *zinxMetrics
var _metricsOnce sync.Once

type zinxMetrics struct {
	// 链接总数
	connTotal *prometheus.GaugeVec //[address, name]:ConnTotal
	// 处理的任务总数
	taskTotal *prometheus.GaugeVec //[address, name, workerID]:TaskTotal
	// 路由Router调度的Handler次数
	routerScheduleTotal *prometheus.GaugeVec //[address, name, workerID, MsgID]:RouterScheduleTotal
	// 路由Router调度的Handler耗时
	routerScheduleDuration *prometheus.HistogramVec //[address, name, workerID, MsgID]:RouterScheduleDuration
}

// Metrics 获取单例
func Metrics() *zinxMetrics {
	_metricsOnce.Do(func() {
		_metrics = new(zinxMetrics)
	})
	return _metrics
}

func (m *zinxMetrics) IsEnable() bool {
	return zconf.GlobalObject.PrometheusMetricsEnable
}

// Zinx的链接数量累加
func (m *zinxMetrics) IncConn(serverAddress, serverName string) {
	if zconf.GlobalObject.PrometheusMetricsEnable {
		m.connTotal.WithLabelValues(serverAddress, serverName).Inc()
	}
}

// Zinx的链接数量累减
func (m *zinxMetrics) DecConn(serverAddress, serverName string) {
	if zconf.GlobalObject.PrometheusMetricsEnable {
		m.connTotal.WithLabelValues(serverAddress, serverName).Dec()
	}
}

// Zinx的任务数量累加
func (m *zinxMetrics) IncTask(address, name, workerID string) {
	if zconf.GlobalObject.PrometheusMetricsEnable {
		m.taskTotal.WithLabelValues(address, name, workerID).Inc()
	}
}

func (m *zinxMetrics) IncRouterSchedule(address, name, workerID, msgID string) {
	if zconf.GlobalObject.PrometheusMetricsEnable {
		m.routerScheduleTotal.WithLabelValues(address, name, workerID, msgID).Inc()
	}
}

func (m *zinxMetrics) ObserveRouterScheduleDuration(address, name, workerID, msgID string, duration time.Duration) {
	if zconf.GlobalObject.PrometheusMetricsEnable {
		m.routerScheduleDuration.With(
			prometheus.Labels{
				LABEL_ADDRESS:   address,
				LABEL_NAME:      name,
				LABEL_WORKER_ID: workerID,
				LABEL_MSG_ID:    msgID,
			}).Observe(duration.Seconds() * 1000)
	}
}
