package importer

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/bluehawk27/medcost-go/util"

	log "github.com/Sirupsen/logrus"
	"github.com/bluehawk27/medcost-go/store"
)

// Importer importing stuff
type Importer struct {
	Store store.StoreType
}

// NewImporter : initiates connection with db to begin imports
func NewImporter() *Importer {
	store, err := store.NewStore()
	if err != nil {
		return nil
	}

	imp := &Importer{
		Store: store,
	}
	return imp
}

func (i *Importer) ImportInpatientData(f string) error {

	csvFile, err := os.Open(f)
	if err != nil {
		return err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	lines, err := reader.ReadAll()
	if err != nil {
		if err == io.EOF {
			log.Fatal("The file is empty")
		}
		return err
	}

	for i, line := range lines {
		if i == 0 {
			continue //skip header
		}

		drg := util.ParseDrg(line[0])
		fmt.Println(*drg.Code)
		fmt.Println(*drg.Name)
		fmt.Println("DRG Def: ", line[0])
		fmt.Println("Provoder_id: ", line[1])
		fmt.Println("Provider_name: ", line[2])
		fmt.Println("Provider_street_address: ", line[3])
		fmt.Println("Provider_city: ", line[4])
		fmt.Println("Provider_state: ", line[5])
		fmt.Println("Provider_zip_code: ", line[6])
		fmt.Println("Provider_HRR: ", line[7])
		fmt.Println("Provider_total_discharges: ", line[8])
		fmt.Println("Provider_avg_covered_charges: ", line[9])
		fmt.Println("Provider_avg_total_payments: ", line[10])
		fmt.Println("Provider_avg_medicare_payments: ", line[11])
		fmt.Println("#################################################")
		// p := store.Provider{}

		// i.Store.InsertProvider()
	}

	return nil
}

// var people []Person
// for {
// 	line, error := reader.Read()
// 	if error == io.EOF {
// 		break
// 	} else if error != nil {
// 		log.Fatal(error)
// 	}
// 	people = append(people, Person{
// 		Firstname: line[0],
// 		Lastname:  line[1],
// 		Address: &Address{
// 			City:  line[2],
// 			State: line[3],
// 		},
// 	})
// }
// peopleJson, _ := json.Marshal(people)
// fmt.Println(string(peopleJson))
