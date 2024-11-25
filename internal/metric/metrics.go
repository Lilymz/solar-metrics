package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"solar-metrics/internal/model"
)

type Metrics struct {
	solarTotalPower *prometheus.GaugeVec
}

func NewMetrics(reg *prometheus.Registry) *Metrics {
	// 初始化指标
	solarTotalPower := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solar_total_power",           // 指标名称
			Help: "Total solar power generated", // 指标描述
		},
		[]string{"collector_code", "project_id", "equipment_id", "original_power"}, // 标签
	)
	// 注册指标
	reg.MustRegister(solarTotalPower)
	return &Metrics{
		solarTotalPower: solarTotalPower,
	}
}

// IncSolarTotalPower 指标采集处理封装
func (m *Metrics) IncSolarTotalPower(power model.EquipmentOriginalPower) {
	// 更新 Counter 值
	m.solarTotalPower.With(prometheus.Labels{
		"collector_code": power.CollectorCode,
		"project_id":     power.ProjectId,
		"equipment_id":   power.EquipmentId,
	}).Set(power.OriginalPower)
}
