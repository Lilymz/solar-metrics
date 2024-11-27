package model

type EquipmentOriginalPower struct {
	CollectorCode string `json:"collectorCode"`
	EquipmentId   string `json:"equipmentId"`
	ProjectId     string `json:"projectId"`
	OriginalPower string `json:"originalPower"`
}
