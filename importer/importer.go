package importer

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
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

// ImportInpatientData : imports the files for inpatient data
func (i *Importer) ImportInpatientData(ctx context.Context, f string) error {

	csvFile, err := os.Open(f)
	if err != nil {
		log.Error(err)
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
		if index == 0 {
			continue //skip header
		}

		impIn := ImportInpatient{}
		inp, err := i.lineToImportInpatient(ctx, line, &impIn)
		if err != nil {
			return err
		}

		drgID, err := i.getOrInsertDRG(ctx, &inp.Drg)
		if err != nil {
			return err
		}

		provID, err := i.getOrInsertProvider(ctx, &inp.Provider)
		if err != nil {
			log.Error(err)
			return err
		}

		inp.Inpatient.ProviderID = provID
		inp.Inpatient.DrgID = drgID
		inp.Inpatient.Year = util.DetectYear(f)
		inpSvc, err := i.getOrInsertInpatientSvc(ctx, inp)
		if err != nil {
			log.Fatal("error inserting Inpatient data", err)
			return err
		}
		log.Info(inpSvc)
		fmt.Println("#################################################")
	}

	return nil
}

func (i *Importer) getOrInsertDRG(ctx context.Context, drg *store.DRG) (*int, error) {
	drgsByID, err := i.Store.GetDrgByDRGID(ctx, drg.Code)
	if err != nil {
		log.Error("error Getting the DRG by ID", err)
		return nil, err
	}

	if len(*drgsByID) < 1 {
		insertedDrg, err := i.Store.InsertDrg(ctx, drg)
		if err != nil {
			log.Error("error Inserting the DRG", err, drg)
			return nil, err
		}

		insDrg := *insertedDrg
		return &insDrg[0].ID, nil
	}

	returnedDrgs := *drgsByID
	return &returnedDrgs[0].ID, nil
}

func (i *Importer) getOrInsertProvider(ctx context.Context, prov *store.Provider) (*int, error) {
	provByID, err := i.Store.GetProviderByProviderID(ctx, prov.ProviderID)
	if err != nil {
		log.Error("error Getting the providerby ID", err)
		return nil, err
	}

	if len(*provByID) < 1 {
		lat, long, err := util.GeocodeAddress(*prov)
		if err != nil {
			log.Error("Error geocoding the address", err)
		}

		prov.Latitude = &lat
		prov.Longitude = &long
		insertedProv, err := i.Store.InsertProvider(ctx, prov)
		if err != nil {
			log.Error("Error Inserting the provider", err, prov.ID)
			return nil, err
		}

		insProv := *insertedProv
		return &insProv[0].ID, nil
	}

	returnedProvs := *provByID
	return &returnedProvs[0].ID, nil
}

func (i *Importer) lineToImportInpatient(ctx context.Context, line []string, impIn *ImportInpatient) (*ImportInpatient, error) {
	drg := util.ParseDrg(line[0])
	impIn.Drg = *drg
	impIn.Provider.ProviderID = util.StringToInt(line[1])
	impIn.Provider.Name = &line[2]
	impIn.Provider.Street = &line[3]
	impIn.Provider.City = &line[4]
	impIn.Provider.State = &line[5]
	impIn.Provider.Zipcode = util.StringToInt(line[6])
	impIn.Provider.HrrDescription = &line[7]
	impIn.Inpatient.TotalDiscarges = util.StringToInt(line[8])
	impIn.Inpatient.AvgCoveredCharges = util.StringToFloat64(line[9])
	impIn.Inpatient.AvgTotalPayments = util.StringToFloat64(line[10])
	impIn.Inpatient.AvgMedicarePayment = util.StringToFloat64(line[11])

	return impIn, nil
}

func (i *Importer) getOrInsertInpatientSvc(ctx context.Context, inp *ImportInpatient) (*int, error) {
	inpSvc, err := i.Store.GetInpatientServiceByProvIDDrgIDYear(ctx, inp.Inpatient.ProviderID, inp.Inpatient.DrgID, inp.Inpatient.Year)
	if err != nil {
		log.Error("error Getting the inpatienservice", err)
		return nil, err
	}

	if len(*inpSvc) < 1 {
		insertedInpSvc, err := i.Store.InsertInpatientService(ctx, &inp.Inpatient)
		if err != nil {
			log.Fatal("error inserting Inpatient data", err)
			return nil, err
		}

		insInpatient := *insertedInpSvc
		log.Info("inpatient Service inserted:")
		return &insInpatient[0].ID, nil
	}

	log.Info("inpatient Service Already Existed: ")
	returnedInpatient := *inpSvc
	return &returnedInpatient[0].ID, nil
}
