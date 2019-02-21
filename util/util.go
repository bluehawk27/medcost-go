package util

import (
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/bluehawk27/medcost-go/store"
)

var years = [...]string{"2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018"}

func ParseDrg(drg string) *store.DRG {
	drgStruct := store.DRG{}
	stringSlice := strings.Split(drg, " - ")

	drgStruct.Code = StringToInt(stringSlice[0])
	drgStruct.Name = &stringSlice[1]

	return &drgStruct
}

func StringToInt(s string) *int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Error("Error converting string to int")
		return nil
	}

	return &i
}

func StringToFloat64(s string) *float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Error("Error converting string to float64")
		return nil
	}

	return &f
}

func DetectYear(s string) *int {
	defaultYear := 0000
	for _, year := range years {
		if strings.Contains(s, year) {
			yearInt := StringToInt(year)
			return yearInt
		}
	}

	return &defaultYear
}
