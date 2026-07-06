package aps

import (
	"math"
	"sort"
	"time"

	"turtle/core/users"
	"turtle/manufacutring/mrp"
	"turtle/manufacutring/routing"
	"turtle/manufacutring/workcenters"
)

const dateLayout = "2006-01-02"

// RunAps produces a finite-capacity forward schedule:
//
//  1. Run MRP to obtain the production orders that must be made in-house.
//  2. Sort them earliest-due-date first (a simple, stable dispatching rule).
//  3. Walk each order's routing operation by operation. Each operation starts
//     at the later of (its work center becoming free) and (the previous
//     operation of the same order finishing), so no work center ever runs two
//     operations at once — that is the finite capacity constraint.
func RunAps(user *users.User) *ApsResult {
	base := startOfNextWorkday(time.Now())

	mrpResult := mrp.RunMrp(user)

	routingByProduct := map[string]routing.Routing{}
	for _, rt := range routing.ListRoutings(user) {
		if rt.ProductUid != "" {
			routingByProduct[rt.ProductUid] = rt
		}
	}

	wcByUid := map[string]workcenters.WorkCenter{}
	for _, w := range workcenters.ListWorkCenters(user) {
		wcByUid[w.Uid.Hex()] = w
	}

	// Production orders only, earliest due date first.
	var orders []mrp.PlannedOrder
	for _, po := range mrpResult.PlannedOrders {
		if po.OrderType == "production" {
			orders = append(orders, po)
		}
	}
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].DueDate != orders[j].DueDate {
			return orders[i].DueDate < orders[j].DueDate
		}
		return orders[i].BomLevel > orders[j].BomLevel // deeper items first
	})

	result := &ApsResult{
		GeneratedAt:     time.Now().Format(time.RFC3339),
		Operations:      []ScheduledOperation{},
		WorkCenterLoads: []WorkCenterLoad{},
		Unscheduled:     []string{},
	}

	wcFree := map[string]time.Time{}  // next free moment per work center
	loadHours := map[string]float64{} // scheduled hours per work center
	opCount := map[string]int{}       //
	var minStart, maxEnd time.Time
	haveSpan := false

	for _, po := range orders {
		rt, ok := routingByProduct[po.ItemUid]
		if !ok || len(rt.Operations) == 0 {
			result.Unscheduled = append(result.Unscheduled, po.Sku)
			continue
		}

		ops := append([]routing.Operation(nil), rt.Operations...)
		sort.SliceStable(ops, func(i, j int) bool { return ops[i].Sequence < ops[j].Sequence })

		// The order cannot start before its planned release date.
		prevEnd := base
		if rel, err := time.Parse(dateLayout, po.ReleaseDate); err == nil {
			relStart := startOfNextWorkday(rel)
			if relStart.After(prevEnd) {
				prevEnd = relStart
			}
		}
		due := parseDueEnd(po.DueDate)

		for _, op := range ops {
			wc := wcByUid[op.WorkCenterUid]
			wcName := wc.Name
			if wcName == "" {
				wcName = op.WorkCenterName
			}

			free := base
			if f, ok := wcFree[op.WorkCenterUid]; ok {
				free = f
			}
			start := laterOf(free, prevEnd)

			minutes := op.SetupMinutes + op.RunMinutesPerUnit*po.Quantity
			eff := wc.Efficiency
			if eff <= 0 {
				eff = 100
			}
			minutes = minutes / (eff / 100.0)
			durationHours := minutes / 60.0
			end := start.Add(time.Duration(minutes * float64(time.Minute)))

			wcFree[op.WorkCenterUid] = end
			prevEnd = end
			loadHours[op.WorkCenterUid] += durationHours
			opCount[op.WorkCenterUid]++

			late := !due.IsZero() && end.After(due)

			result.Operations = append(result.Operations, ScheduledOperation{
				OrderRef:       po.Sku,
				ItemUid:        po.ItemUid,
				Sku:            po.Sku,
				OperationName:  op.Name,
				Sequence:       op.Sequence,
				WorkCenterUid:  op.WorkCenterUid,
				WorkCenterName: wcName,
				Quantity:       po.Quantity,
				Start:          start.Format(time.RFC3339),
				End:            end.Format(time.RFC3339),
				DurationHours:  round(durationHours),
				DueDate:        po.DueDate,
				Late:           late,
			})

			if !haveSpan || start.Before(minStart) {
				minStart = start
			}
			if !haveSpan || end.After(maxEnd) {
				maxEnd = end
			}
			haveSpan = true
		}
	}

	if haveSpan {
		result.HorizonStart = minStart.Format(time.RFC3339)
		result.HorizonEnd = maxEnd.Format(time.RFC3339)
	}

	horizonDays := 1.0
	if haveSpan {
		horizonDays = math.Max(1, math.Ceil(maxEnd.Sub(minStart).Hours()/24.0))
	}

	for uid, load := range loadHours {
		wc := wcByUid[uid]
		capPerDay := wc.CapacityHoursPerDay
		if capPerDay <= 0 {
			capPerDay = 8
		}
		capacity := capPerDay * horizonDays
		util := 0.0
		if capacity > 0 {
			util = load / capacity * 100
		}
		name := wc.Name
		if name == "" {
			name = uid
		}
		result.WorkCenterLoads = append(result.WorkCenterLoads, WorkCenterLoad{
			WorkCenterUid:  uid,
			WorkCenterName: name,
			LoadHours:      round(load),
			CapacityHours:  round(capacity),
			Utilization:    round(util),
			Operations:     opCount[uid],
		})
	}
	sort.SliceStable(result.WorkCenterLoads, func(i, j int) bool {
		return result.WorkCenterLoads[i].Utilization > result.WorkCenterLoads[j].Utilization
	})

	return result
}

// startOfNextWorkday returns the given day at 08:00 (a simple shift start).
func startOfNextWorkday(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 8, 0, 0, 0, t.Location())
}

func parseDueEnd(s string) time.Time {
	if t, err := time.Parse(dateLayout, s); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 17, 0, 0, 0, t.Location())
	}
	return time.Time{}
}

func laterOf(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func round(v float64) float64 {
	return math.Round(v*100) / 100
}
