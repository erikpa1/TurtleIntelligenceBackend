package tools

import (
	"turtle/interfaces"
)

func CheckOrCreateUid[T interfaces.UidProvider](ele T) {
	if ele.GetUid() == "" {
		ele.SetUid(GetUUID4())
	}
}

func CheckOrCreateId[T interfaces.IdProvider](ele T) {
	if ele.GetId() == 0 {
		ele.SetId(GetTimeNowMillis())
	}

}
