package util

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
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

func ParseApc(apc string) *store.APC {
	apcStruct := store.APC{}
	stringSlice := strings.Split(apc, " - ")

	apcStruct.Code = StringToInt(stringSlice[0])
	apcStruct.Name = &stringSlice[1]

	return &apcStruct
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
	s = normalizeMoneyFields(s)

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Error("Error converting string to float64 :", err)
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

func normalizeMoneyFields(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}

	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	if strings.Contains(s, "$") {
		s = strings.Trim(s, "$")

	}

	if strings.Contains(s, ",") {
		s = strings.Replace(s, ",", "", -1)
	}

	return s
}

func FileToLines(f string) ([][]string, error) {
	csvFile, err := os.Open(f)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	lines, err := reader.ReadAll()
	if err != nil {
		if err == io.EOF {
			log.Fatal("The file is empty")
		}
		return nil, err
	}

	if len(lines) <= 1 || lines == nil {
		return nil, errors.New("this csv has no lines to insert")
	}

	return lines, nil
}
