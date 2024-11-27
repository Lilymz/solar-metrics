package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"solar-metrics/internal/model"
	"strconv"
)

type Metrics struct {
	solarTotalPower *prometheus.GaugeVec
}

var (
	metrics *Metrics
	reg     *prometheus.Registry
)

func init() {
	reg = prometheus.NewRegistry()
	metrics = newMetrics(reg)
}
func Registry() *prometheus.Registry {
	return reg
}
func newMetrics(reg *prometheus.Registry) *Metrics {
	// 初始化指标
	solarTotalPower := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solar_total_power",           // 指标名称
			Help: "Total solar power generated", // 指标描述
		},
		[]string{"collector_code", "project_id", "equipment_id"}, // 标签
	)
	// 注册指标
	reg.MustRegister(solarTotalPower)
	return &Metrics{
		solarTotalPower: solarTotalPower,
	}
}
func GetMetrics() *Metrics {
	return metrics
}

// IncSolarTotalPower 指标采集处理封装
func (m *Metrics) IncSolarTotalPower(power model.EquipmentOriginalPower) {
	// 更新 Counter 值
	value, _ := strconv.ParseFloat(power.OriginalPower, 64)
	m.solarTotalPower.With(prometheus.Labels{
		"collector_code": power.CollectorCode,
		"project_id":     power.ProjectId,
		"equipment_id":   power.EquipmentId,
	}).Set(value)
}
