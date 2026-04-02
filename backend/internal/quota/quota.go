package quota

import "fmt"

// PlanConfig 套餐配置
type PlanConfig struct {
	Name         string
	MaxVMs       int
	MaxConns     int
	MaxMounts    int
	MaxSnapshots int
}

var plans = map[string]PlanConfig{
	"free":       {Name: "免费版", MaxVMs: 1, MaxConns: 1, MaxMounts: 2, MaxSnapshots: 3},
	"pro":        {Name: "专业版", MaxVMs: 5, MaxConns: 5, MaxMounts: 10, MaxSnapshots: 20},
	"enterprise": {Name: "企业版", MaxVMs: -1, MaxConns: -1, MaxMounts: -1, MaxSnapshots: -1},
}

// GetPlan returns the plan config for the given plan name.
func GetPlan(plan string) PlanConfig {
	if p, ok := plans[plan]; ok {
		return p
	}
	return plans["free"]
}

// CheckQuota checks if the current usage exceeds the plan limit.
func CheckQuota(plan string, resource string, current int) error {
	p := GetPlan(plan)
	var limit int
	switch resource {
	case "vms":
		limit = p.MaxVMs
	case "connections":
		limit = p.MaxConns
	case "mounts":
		limit = p.MaxMounts
	case "snapshots":
		limit = p.MaxSnapshots
	default:
		return nil
	}
	if limit == -1 {
		return nil
	}
	if current >= limit {
		return fmt.Errorf("%s 配额已满（%s 套餐限制 %d 个），请升级", resource, p.Name, limit)
	}
	return nil
}
