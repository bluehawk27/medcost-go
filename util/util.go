package util

import (
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/bluehawk27/medcost-go/store"
)

func ParseDrg(drg string) *store.DRG {
	drgStruct := store.DRG{}
	stringSlice := strings.Split(drg, " - ")
	code, err := strconv.Atoi(stringSlice[0])
	if err != nil {
		log.Error("Error converting string to int")
		return nil
	}
	drgStruct.Code = &code
	drgStruct.Name = &stringSlice[1]

	return &drgStruct
}
