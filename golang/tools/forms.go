package tools

import (
	"turtle/interfaces"
)

func CheckOrCreateUid[T interfaces.UidProvider](ele T) string {
	if ele.GetUid() == "" {
		ele.SetUid(GetUUID4())
	}

	return ele.GetUid()
}

func CheckOrCreateShortUid[T interfaces.UidProvider](ele T, pool string) {
	if ele.GetUid() == "" {

		ele.SetUid(ShortUid(pool))
	}
}

func CheckOrCreateId[T interfaces.IdProvider](ele T) {
	if ele.GetId() == 0 {
		ele.SetId(GetTimeNowMillis())
	}

}
