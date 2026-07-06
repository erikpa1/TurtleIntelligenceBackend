# Inventory module

Material master (products / parts) used by both stock keeping and manufacturing.
Modelled after the SAP **material master (MM)**.

## What is here

| File | Purpose |
|------|---------|
| `items/items_models.go` | `Item` entity (material master record) |
| `items/items_ctrl.go`   | Org-scoped CRUD (`CT_ITEMS = "inventory_items"`) |
| `items/items_api.go`    | REST routes |
| `a.go`                  | `InitInventoryApi(r)` — wired from `main.go` |

The API follows the same hand-written pattern as `themes/` (models + ctrl + api,
`ObjFromJsonPtr` + `_id,omitempty`, `user.FillOrgQuery`).

### Routes
- `GET  /api/inventory/items` — list (org scoped)
- `GET  /api/inventory/item?uid=` — one
- `POST /api/inventory/item` — create/update (form field `data` = JSON)
- `DELETE /api/inventory/item?uid=`

### `Item` fields
Basic: `sku, name, description, category (raw|semiFinished|finished|trading),
type, uom, unitPrice, currency, qtyOnHand, reorderPoint, warehouse, active`.

**MRP / planning view** (drives `manufacutring/mrp`):
`procurementType (make|buy)`, `leadTimeDays`, `safetyStock`, `lotSize`
(0 = lot-for-lot), `minLotSize`. Empty `procurementType` is inferred at planning
time (has a BOM → make, else buy).

## Where this is in SAP

Transactions **MM01 / MM02 / MM03** (create / change / display material).
This entity is a simplified flattening of several SAP "views" into one document.

| Turtle field | SAP view | SAP table.field |
|---|---|---|
| `sku` | Basic Data | `MARA.MATNR` |
| `name` / `description` | Basic Data | `MAKT.MAKTX` |
| `category` / `type` | Basic Data | `MARA.MTART` (material type) |
| `uom` | Basic Data | `MARA.MEINS` (base UoM) |
| `unitPrice` / `currency` | Accounting | `MBEW.STPRS` / `VERPR` |
| `qtyOnHand` | Storage | `MARD.LABST` |
| `warehouse` | Plant/Sloc | `MARC.WERKS` / `MARD.LGORT` |
| `reorderPoint` | MRP 1 | `MARC.MINBE` |
| `procurementType` | MRP 2 | `MARC.BESKZ` (E=in-house, F=external) |
| `leadTimeDays` | MRP 2 | `MARC.DZEIT` / `PLIFZ` |
| `safetyStock` | MRP 2 | `MARC.EISBE` |
| `lotSize` / `minLotSize` | MRP 1 | `MARC.DISLS` / `BSTMI` |

See `../manufacutring/CLAUDE.md` for the planning engines that consume this data.
