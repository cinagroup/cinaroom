package repository

import (
	"fmt"

	"github.com/cinagroup/cinaseek/backend/internal/model"
)

// MetricRepo provides operations for the vm_metrics table.
type MetricRepo struct {
}

// NewMetricRepo creates a new MetricRepo.
func NewMetricRepo() *MetricRepo {
	return &MetricRepo{}
}

// Create inserts a new metric sample.
func (r *MetricRepo) Create(metric *model.VMMetric) error {
	db := GetDB()
	if err := db.Create(metric).Error; err != nil {
		return fmt.Errorf("create metric: %w", err)
	}
	return nil
}

// ListByVM returns the most recent metrics for a VM, limited to n.
func (r *MetricRepo) ListByVM(vmID uint, limit int) ([]model.VMMetric, error) {
	db := GetDB()
	var metrics []model.VMMetric
	q := db.Where("vm_id = ?", vmID).Order("timestamp DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&metrics).Error; err != nil {
		return nil, fmt.Errorf("list metrics for vm %d: %w", vmID, err)
	}
	return metrics, nil
}
