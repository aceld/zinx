package zmetrics

import (
	"github.com/aceld/zinx/zconf"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var _metrics *zinxMetrics
var _metricsOnce sync.Once

type zinxMetrics struct {
	// 链接总数
	connTotal *prometheus.GaugeVec
	// 处理的任务总数
	taskTotal *prometheus.GaugeVec //workerID:TaskTotal

	// TODO 路由Router调度的Handler次数
	// TODO 路由Router调度的Handler耗时
	// TODO 拦截器处理数据的次数
	// TODO 拦截器处理数据的耗时
	// TODO Handler调度错误
}

// Metrics 获取单例
func Metrics() *zinxMetrics {
	_metricsOnce.Do(func() {
		_metrics = new(zinxMetrics)
	})
	return _metrics
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
func (m *zinxMetrics) IncTask(workerID string) {
	if zconf.GlobalObject.PrometheusMetricsEnable {
		m.taskTotal.WithLabelValues(workerID).Inc()
	}
}
