package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var ShortUidFunc func(string) string = nil

func ShortUid(pool string) string {

	if ShortUidFunc != nil {
		return ShortUidFunc(pool)
	} else {
		newUUID, _ := uuid.NewRandom()
		return newUUID.String()
	}

}

func GetUUID4() string {
	//newUUID, _ := uuid.NewRandom()
	return GetUUID4_Shorter()
}

func GetUUID4_Shorter() string {
	newUUID, _ := uuid.NewRandom()
	return strings.Replace(newUUID.String(), "-", "", -1)
}

func ShortenUUIDs(strings []string) map[string]string {
	uidMap := make(map[string]string)
	uidSet := make(map[string]bool) // To track used UIDs

	for _, s := range strings {
		hash := md5.Sum([]byte(s))
		shortUID := hex.EncodeToString(hash[:])[:4]

		counter := 1
		uniqueUID := shortUID
		for uidSet[uniqueUID] {
			uniqueUID = fmt.Sprintf("%s%d", shortUID, counter)
			counter++
		}

		uidMap[s] = uniqueUID
		uidSet[uniqueUID] = true
	}

	return uidMap
}
