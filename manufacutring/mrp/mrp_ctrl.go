package mrp

import (
	"math"
	"sort"
	"time"

	"turtle/core/users"
	"turtle/inventory/items"
	"turtle/manufacutring/bom"
	"turtle/manufacutring/demand"
)

const dateLayout = "2006-01-02"
const maxBomDepth = 25

// planItem is the mutable per-item accumulator used during the run.
type planItem struct {
	item     items.Item
	level    int
	gross    float64
	required time.Time // earliest date the quantity is needed
	hasDate  bool
}

// RunMrp executes a single, aggregate (single-bucket) material requirements
// planning run for the caller's organisation:
//
//  1. Load items, active BOMs and open demand.
//  2. Compute low-level codes so every item is netted only after all its
//     parents have contributed their dependent demand.
//  3. Net each item (gross + safety - on hand), apply lot sizing and offset by
//     lead time to obtain a planned order release date.
//  4. Explode "make" items through their BOM into dependent demand.
func RunMrp(user *users.User) *MrpResult {
	now := truncateDay(time.Now())

	itemList := items.ListItems(user)
	bomList := bom.ListBoms(user)
	demandList := demand.ListDemand(user)

	itemByUid := map[string]items.Item{}
	for _, it := range itemList {
		itemByUid[it.Uid.Hex()] = it
	}

	// Pick one BOM per product, preferring an "active" one.
	bomByProduct := map[string]bom.Bom{}
	for _, b := range bomList {
		if b.ProductUid == "" {
			continue
		}
		existing, ok := bomByProduct[b.ProductUid]
		if !ok || (b.Status == "active" && existing.Status != "active") {
			bomByProduct[b.ProductUid] = b
		}
	}

	llc := computeLowLevelCodes(itemByUid, bomByProduct)

	// The netting loop must reach the deepest component level, not just the
	// level of the demanded (finished) items, otherwise dependent demand would
	// never be exploded.
	maxLevel := 0
	for _, l := range llc {
		if l > maxLevel {
			maxLevel = l
		}
	}

	plan := map[string]*planItem{}
	ensure := func(uid string) *planItem {
		if p, ok := plan[uid]; ok {
			return p
		}
		p := &planItem{item: itemByUid[uid], level: llc[uid]}
		plan[uid] = p
		return p
	}

	result := &MrpResult{
		GeneratedAt:   time.Now().Format(time.RFC3339),
		PlannedOrders: []PlannedOrder{},
		Requirements:  []RequirementRow{},
		Exceptions:    []MrpException{},
	}

	// Independent demand seeds the top level items.
	for _, d := range demandList {
		if d.Status == "closed" || d.ProductUid == "" {
			continue
		}
		if _, ok := itemByUid[d.ProductUid]; !ok {
			result.Exceptions = append(result.Exceptions, MrpException{
				Severity: "error",
				ItemUid:  d.ProductUid,
				Sku:      d.ProductSku,
				Message:  "Demand references an unknown material",
			})
			continue
		}
		p := ensure(d.ProductUid)
		p.gross += d.Quantity
		addRequiredDate(p, parseDate(d.DueDate, now))
	}

	// Net level by level; exploding a parent injects demand into deeper levels.
	for level := 0; level <= maxLevel; level++ {
		// Snapshot current uids at this level (map grows as we explode).
		var atLevel []string
		for uid, p := range plan {
			if p.level == level {
				atLevel = append(atLevel, uid)
			}
		}
		sort.Strings(atLevel)

		for _, uid := range atLevel {
			p := plan[uid]
			it := p.item

			procurement := procurementType(it, bomByProduct)
			required := p.required
			if !p.hasDate {
				required = now
			}
			release := required.AddDate(0, 0, -int(it.LeadTimeDays))

			net := p.gross + it.SafetyStock - it.QtyOnHand
			plannedQty := 0.0
			if net > 0 {
				plannedQty = applyLotSize(net, it)
			}

			result.Requirements = append(result.Requirements, RequirementRow{
				ItemUid:          uid,
				Sku:              it.Sku,
				Name:             it.Name,
				BomLevel:         p.level,
				GrossRequirement: round(p.gross),
				OnHand:           it.QtyOnHand,
				SafetyStock:      it.SafetyStock,
				NetRequirement:   round(math.Max(0, net)),
				PlannedOrder:     round(plannedQty),
				ProcurementType:  procurement,
				RequiredDate:     required.Format(dateLayout),
			})

			if plannedQty <= 0 {
				continue
			}

			orderType := "purchase"
			if procurement == "make" {
				orderType = "production"
			}

			result.PlannedOrders = append(result.PlannedOrders, PlannedOrder{
				ItemUid:     uid,
				Sku:         it.Sku,
				Name:        it.Name,
				OrderType:   orderType,
				Quantity:    round(plannedQty),
				Uom:         it.Uom,
				ReleaseDate: release.Format(dateLayout),
				DueDate:     required.Format(dateLayout),
				BomLevel:    p.level,
			})

			if release.Before(now) {
				result.Exceptions = append(result.Exceptions, MrpException{
					Severity: "warning",
					ItemUid:  uid,
					Sku:      it.Sku,
					Message:  "Planned release date is in the past — expedite",
				})
			}

			// Explode production orders into component demand.
			if orderType == "production" {
				b, ok := bomByProduct[uid]
				if !ok {
					result.Exceptions = append(result.Exceptions, MrpException{
						Severity: "error",
						ItemUid:  uid,
						Sku:      it.Sku,
						Message:  "Make item has no BOM — cannot explode requirements",
					})
					continue
				}
				base := b.BaseQuantity
				if base <= 0 {
					base = 1
				}
				for _, comp := range b.Components {
					if comp.ItemUid == "" {
						continue
					}
					if _, ok := itemByUid[comp.ItemUid]; !ok {
						result.Exceptions = append(result.Exceptions, MrpException{
							Severity: "error",
							ItemUid:  comp.ItemUid,
							Sku:      comp.Sku,
							Message:  "BOM component is not in the material master",
						})
						continue
					}
					perUnit := (comp.Quantity / base) * (1 + comp.ScrapPct/100)
					child := ensure(comp.ItemUid)
					child.gross += plannedQty * perUnit
					// The component is needed when the parent order is released.
					addRequiredDate(child, release)
				}
			}
		}
	}

	sortResult(result)
	return result
}

func procurementType(it items.Item, bomByProduct map[string]bom.Bom) string {
	if it.ProcurementType == "make" || it.ProcurementType == "buy" {
		return it.ProcurementType
	}
	if _, ok := bomByProduct[it.Uid.Hex()]; ok {
		return "make"
	}
	return "buy"
}

func applyLotSize(net float64, it items.Item) float64 {
	qty := net
	if it.MinLotSize > 0 && qty < it.MinLotSize {
		qty = it.MinLotSize
	}
	if it.LotSize > 0 {
		qty = math.Ceil(qty/it.LotSize) * it.LotSize
	}
	return qty
}

// computeLowLevelCodes assigns each item the deepest level at which it appears
// in any BOM. Depth is bounded by maxBomDepth so a cyclic BOM cannot recurse
// forever.
func computeLowLevelCodes(itemByUid map[string]items.Item, bomByProduct map[string]bom.Bom) map[string]int {
	llc := map[string]int{}
	for uid := range itemByUid {
		llc[uid] = 0
	}

	var visit func(uid string, depth int)
	visit = func(uid string, depth int) {
		if depth > llc[uid] {
			llc[uid] = depth
		}
		if depth >= maxBomDepth {
			return
		}
		b, ok := bomByProduct[uid]
		if !ok {
			return
		}
		for _, comp := range b.Components {
			if comp.ItemUid != "" {
				visit(comp.ItemUid, depth+1)
			}
		}
	}

	for uid := range itemByUid {
		visit(uid, 0)
	}
	return llc
}

func addRequiredDate(p *planItem, d time.Time) {
	if !p.hasDate || d.Before(p.required) {
		p.required = d
		p.hasDate = true
	}
}

func parseDate(s string, fallback time.Time) time.Time {
	if s == "" {
		return fallback
	}
	if t, err := time.Parse(dateLayout, s); err == nil {
		return truncateDay(t)
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return truncateDay(t)
	}
	return fallback
}

func truncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func round(v float64) float64 {
	return math.Round(v*1000) / 1000
}

func sortResult(result *MrpResult) {
	sort.SliceStable(result.PlannedOrders, func(i, j int) bool {
		a, b := result.PlannedOrders[i], result.PlannedOrders[j]
		if a.BomLevel != b.BomLevel {
			return a.BomLevel < b.BomLevel
		}
		if a.ReleaseDate != b.ReleaseDate {
			return a.ReleaseDate < b.ReleaseDate
		}
		return a.Sku < b.Sku
	})
	sort.SliceStable(result.Requirements, func(i, j int) bool {
		a, b := result.Requirements[i], result.Requirements[j]
		if a.BomLevel != b.BomLevel {
			return a.BomLevel < b.BomLevel
		}
		return a.Sku < b.Sku
	})
}
