package importer

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

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

// ImportInpatientData : imports the files for inpatient data
func (i *Importer) ImportInpatientData(ctx context.Context, f string) error {

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
	if len(lines) <= 1 || lines == nil {
		return errors.New("this csv has no lines to insert")
	}

	for index, line := range lines {
		impIn := ImportInpatient{}
		inp, err := i.lineToImportInpatient(ctx, line, &impIn)
		if err != nil {
			return err
		}

		if index == 0 {
			continue //skip header
		}
		// lat, long := util.GeocodeAddress(p)
		drg := util.ParseDrg(line[0])
		drgID, err := i.getOrInsertDRG(ctx, drg)
		if err != nil {
			return err
		}

		provID, err := i.getOrInsertProvider(ctx, inp.Provider)

		inp.Inpatient.ProviderID = provID
		inp.Inpatient.DrgID = drgID
		inpSvc, err := i.Store.InsertInpatientService(ctx, inp.Inpatient)
		if err != nil {
			log.Fatal("error inserting patient data", err)
			return err
		}

		fmt.Println("DRG Def: ", line[0])
		fmt.Println("Provoder_id: ", line[1])
		fmt.Println("Provider_name: ", line[2])
		fmt.Println("Provider_street_address: ", line[3])
		fmt.Println("Provider_city: ", line[4])
		fmt.Println("Provider_state: ", line[5])
		fmt.Println("Provider_zip_code: ", line[6])
		fmt.Println("Provider_HRR: ", line[7])
		fmt.Println("Provider_inpatient_total_discharges: ", line[8])
		fmt.Println("Provider_inpatient_avg_covered_charges: ", line[9])
		fmt.Println("Provider_inpatient_avg_total_payments: ", line[10])
		fmt.Println("Provider_inpatient_avg_medicare_payments: ", line[11])
		fmt.Println("#################################################")
	}

	return nil
}

func (i *Importer) getOrInsertDRG(ctx context.Context, drg *store.DRG) (*int, error) {
	drgsByID, err := i.Store.GetDrgByDRGID(ctx, drg.Code)
	if err != nil {
		return nil, err
	}

	if len(*drgsByID) < 1 {
		insertedDrg, err := i.Store.InsertDrg(ctx, drg)
		if err != nil {
			return nil, err
		}

		insDrg := *insertedDrg
		return &insDrg[0].ID, nil
	}

	returnedDrgs := *drgsByID
	return &returnedDrgs[0].ID, nil
}

func (i *Importer) getOrInsertProvider(ctx context.Context, prov *store.Provider) (*int, error) {
	drgsByID, err := i.Store.GetProviderByProviderID(ctx, prov.ProviderID)
	if err != nil {
		return nil, err
	}

	if len(*drgsByID) < 1 {
		lat, long, err := util.GeocodeAddress(*prov)
		if err != nil {
			log.Fatal("Error geocoding the address")
			return nil, err
		}

		prov.Latitude = &lat
		prov.Longitude = &long
		insertedDrg, err := i.Store.InsertProvider(ctx, prov)
		if err != nil {
			return nil, err
		}

		insDrg := *insertedDrg
		return &insDrg[0].ID, nil
	}

	returnedDrgs := *drgsByID
	return &returnedDrgs[0].ID, nil
}

func (i *Importer) lineToImportInpatient(ctx context.Context, line []string, impIn *ImportInpatient) (ImportInpatient, error) {
	impIn.Provider.ProviderID =
	impIn.Provider.Name =
	impIn.Provider.Street =
	i, err := strconv.Atoi(s)

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
