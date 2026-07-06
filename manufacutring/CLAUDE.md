# Manufacturing module

Bill of materials + Advanced Planning (MRP) and Scheduling (APS).
Folder name is `manufacutring` (historical typo); the Go package is
`manufacturing` and is imported as `manufacturing "turtle/manufacutring"`.

All sub-modules are wired in `a.go` → `InitManufacturingApi(r)` (called from
`main.go`). Every sub-module follows the hand-written `themes/`-style pattern
(models + ctrl + api, org-scoped via `user.FillOrgQuery`).

## Sub-modules

| Folder | Entity / job | Mongo collection | Routes |
|--------|--------------|------------------|--------|
| `bom/`         | `Bom` + `[]BomComponent` (recipe) | `manufacturing_boms` | `/api/manufacturing/bom[s]` |
| `demand/`      | `DemandOrder` (independent demand) | `manufacturing_demand` | `/api/manufacturing/demand[s]` |
| `workcenters/` | `WorkCenter` (capacity resource) | `manufacturing_workcenters` | `/api/manufacturing/workcenter[s]` |
| `routing/`     | `Routing` + `[]Operation` | `manufacturing_routings` | `/api/manufacturing/routing[s]` |
| `mrp/`         | requirements planning engine | — (stateless) | `POST/GET /api/manufacturing/mrp/run` |
| `aps/`         | finite-capacity scheduler | — (stateless) | `POST/GET /api/manufacturing/aps/run` |

CRUD sub-modules expose list (`…s`), get (`…?uid=`), `POST` (create/update,
form field `data`), and `DELETE`. `mrp` and `aps` are pure computations over the
current master data (no persistence).

## Data flow

```
Demand ─┐
BOM ────┼─▶ MRP.RunMrp ─▶ planned orders (production + purchase) + netting + exceptions
Items ──┘                      │
                               ▼
Routings + WorkCenters ─▶ APS.RunAps ─▶ finite-capacity Gantt + work-center load
```

## How the engines work

### MRP (`mrp/mrp_ctrl.go` → `RunMrp`)
Aggregate (single-bucket) MRP:
1. Load items, active BOMs (prefers `status=active`), open demand.
2. **Low-level coding** (`computeLowLevelCodes`, depth-capped at 25 for cycle
   safety) so each item is netted only after all parents contributed dependent
   demand. The netting loop runs `level = 0 … maxLevel(llc)`.
3. Per item: `net = gross + safetyStock − onHand`; `applyLotSize` (min-lot then
   round up to `lotSize`); release date = `requiredDate − leadTimeDays`.
4. `make` items explode through their BOM: `perUnit = qty/baseQty × (1+scrap%)`;
   components inherit the parent's release as their required date.
5. Emits `PlannedOrder`s (production|purchase), a `RequirementRow` netting table,
   and `MrpException`s (expedite / missing BOM / unknown component).

### APS (`aps/aps_ctrl.go` → `RunAps`)
Finite-capacity **forward** scheduling:
1. Runs MRP, takes only `production` orders, dispatches earliest-due-first.
2. Each routing operation starts at `max(work-center-free, previous-op-end)` →
   no work center runs two operations at once (the capacity constraint).
   Duration = `(setup + run×qty) / (efficiency/100)`; shift start = 08:00.
3. Emits `ScheduledOperation`s (start/end, late flag) + per-work-center
   `WorkCenterLoad` (load vs capacity, utilization %).

Both are intentionally pragmatic (aggregate MRP, single-capacity forward APS),
not full time-phased / optimising planners.

## Where this is in SAP

SAP modules: **PP** (Production Planning) and **PP-DS / APO** (Advanced Planning).

| Turtle | SAP transaction | SAP table |
|--------|-----------------|-----------|
| BOM (`bom/`) | **CS01/CS02/CS03** (create/change/display BOM) | `MAST` (matl-BOM link), `STKO` (header), `STPO` (items) |
| Work center (`workcenters/`) | **CR01/CR02/CR03** | `CRHD` (header), `CRCA`/`KAKO` (capacity) |
| Routing (`routing/`) | **CA01/CA02/CA03** | `PLKO` (header), `PLPO` (operations), `PLAS` (sequence) |
| Demand (`demand/`) | **MD61** (planned independent reqs) / **VA01** (sales order) | `PBIM`/`PBED` (PIR), `VBAK`/`VBAP` (SO) |
| MRP run (`mrp/`) | **MD01** (total) / **MD02** (single-item, multi-level) / **MD03** | planned orders `PLAF`, purchase reqs `EBAN` |
| MRP result views | **MD04** (stock/reqs list), **MD05** (MRP list) | — |
| APS (`aps/`) | **PP-DS** detailed scheduling; planning board **/SAPAPO/CDPS0**; capacity leveling **CM01/CM21**, planning table **MF50** | `/SAPAPO/*`, order `AFKO`/`AFVC` |

A machine-readable version of the full mapping (all entities, fields, routes)
lives in **`sap_mapping.json`** in this folder — the Turtle→SAP correspondence
schema (CSP).

### Field-level mapping schema (Turtle → SAP)

BOM component:
- `bom.baseQuantity` → `STKO.BMENG` (base quantity)
- `component.quantity` → `STPO.MENGE`, `component.scrapPct` → `STPO.AUSCH` (component scrap)

Routing operation:
- `operation.workCenterUid` → `PLPO.ARBID` (→ `CRHD`)
- `operation.setupMinutes` → `PLPO.VGW01` (setup / Rüstzeit)
- `operation.runMinutesPerUnit` → `PLPO.VGW02` (machine/labour per unit)

MRP output:
- `PlannedOrder{orderType:production}` → planned order `PLAF` (later converted to
  production order `AUFK`/`AFKO` via **CO40/CO41**)
- `PlannedOrder{orderType:purchase}` → purchase requisition `EBAN` (→ PO `EKKO`)
- `MrpException` → SAP MRP **exception messages** (e.g. 10 "reschedule in",
  30 "plan process", 07 "safety stock shortfall")

See material-master field mapping in `../inventory/CLAUDE.md` and the UI in
`../../TurtleIntelligenceFrontend/src/TurtleManufacturing/CLAUDE.md`.
