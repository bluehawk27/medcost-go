package importer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bluehawk27/medcost-go/util"

	log "github.com/Sirupsen/logrus"
	"github.com/bluehawk27/medcost-go/store"
)

type ImporterType interface {
	ImportInpatientData(ctx context.Context, f string) error
	ImportOutpatientData(ctx context.Context, f string) error
	ImportOutpatientAPIData(ctx context.Context, url string) error
}

// Importer importing stuff
type Importer struct {
	Store store.StoreType
}

// NewImporter : initiates connection with db to begin imports
func NewImporter() ImporterType {
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
	lines, err := util.FileToLines(f)
	if err != nil {
		return err
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

// --------------------------Outpatient Stuff

// ImportOutpatientData imports data from a csv
func (i *Importer) ImportOutpatientData(ctx context.Context, f string) error {
	lines, err := util.FileToLines(f)
	if err != nil {
		return err
	}

	for index, line := range lines {
		if index == 0 {
			continue //skip header
		}

		impOut := ImportOutpatient{}
		importOut, err := i.lineToImportOutpatient(ctx, line, &impOut)
		if err != nil {
			return err
		}

		apcID, err := i.getOrInsertApc(ctx, &importOut.Apc)
		if err != nil {
			return err
		}

		provID, err := i.getOrInsertProvider(ctx, &importOut.Provider)
		if err != nil {
			log.Error(err)
			return err
		}

		importOut.Outpatient.ProviderID = provID
		importOut.Outpatient.ApcID = apcID
		importOut.Outpatient.Year = util.DetectYear(f)

		_, err = i.getOrInsertOutpatientSvc(ctx, importOut)
		if err != nil {
			log.Fatal("error inserting Outpatient data", err)
			return err
		}
		fmt.Println("#################################################")
	}

	return nil
}

func (i *Importer) lineToImportOutpatient(ctx context.Context, line []string, impOut *ImportOutpatient) (*ImportOutpatient, error) {
	apc := util.ParseApc(line[0])
	impOut.Apc = *apc
	impOut.Provider.ProviderID = util.StringToInt(line[1])
	impOut.Provider.Name = &line[2]
	impOut.Provider.Street = &line[3]
	impOut.Provider.City = &line[4]
	impOut.Provider.State = &line[5]
	impOut.Provider.Zipcode = util.StringToInt(line[6])
	impOut.Provider.HrrDescription = &line[7]
	impOut.Outpatient.ServicesCount = util.StringToInt(line[8])
	impOut.Outpatient.AvgEstSubmittedCharges = util.StringToFloat64(line[9])
	impOut.Outpatient.AvgTotalPayments = util.StringToFloat64(line[10])

	return impOut, nil
}

func (i *Importer) getOrInsertApc(ctx context.Context, apc *store.APC) (*int, error) {
	apcsByID, err := i.Store.GetApcByAPCID(ctx, apc.Code)
	if err != nil {
		log.Error("error Getting the DRG by ID", err)
		return nil, err
	}

	if len(*apcsByID) < 1 {
		insertedApc, err := i.Store.InsertApc(ctx, apc)
		if err != nil {
			log.Error("error Inserting the DRG", err, apc)
			return nil, err
		}

		insApc := *insertedApc
		return &insApc[0].ID, nil
	}

	returnedApcs := *apcsByID
	return &returnedApcs[0].ID, nil
}

func (i *Importer) getOrInsertOutpatientSvc(ctx context.Context, out *ImportOutpatient) (*int, error) {
	outSvc, err := i.Store.GetOutpatientServiceByProvIDApcIDYear(ctx, out.Outpatient.ProviderID, out.Outpatient.ApcID, out.Outpatient.Year)
	if err != nil {
		log.Error("error Getting the Outpatient service", err)
		return nil, err
	}

	if len(*outSvc) < 1 {
		insertedOutSvc, err := i.Store.InsertOutpatientService(ctx, &out.Outpatient)
		if err != nil {
			log.Fatal("error inserting Outpatient data", err)
			return nil, err
		}

		insOutpatient := *insertedOutSvc
		log.Info("Outpatient Service inserted:")
		return &insOutpatient[0].ID, nil
	}

	log.Info("Outpatient Service Already Existed: ")
	returnedOutpatient := *outSvc
	return &returnedOutpatient[0].ID, nil
}

func (i *Importer) ImportOutpatientAPIData(ctx context.Context, url string) error {
	outAPI := OutpatientAPISlice{}
	res, err := util.MakeRequest(url)
	if err != nil {
		log.Info(err)
		return err
	}

	json.NewDecoder(res).Decode(&outAPI.OutpatientAPIsrvcs)
	log.Info(len(outAPI.OutpatientAPIsrvcs))
	for _, object := range outAPI.OutpatientAPIsrvcs {
		impOut, err := i.outpatientAPIToImporterOutpatient(object)
		if err != nil {
			return err
		}

		apcID, err := i.getOrInsertApc(ctx, &impOut.Apc)
		if err != nil {
			return err
		}

		provID, err := i.getOrInsertProvider(ctx, &impOut.Provider)
		if err != nil {
			log.Error(err)
			return err
		}

		impOut.Outpatient.ProviderID = provID
		impOut.Outpatient.ApcID = apcID

		_, err = i.getOrInsertOutpatientSvc(ctx, &impOut)
		if err != nil {
			log.Fatal("error inserting Outpatient data", err)
			return err
		}
		// log.Info(impOutSvc)
		// fmt.Println("#################################################")

		// log.Info(impOut)
		// still need to figure out what to do with this...... ? update DB?
		// log.Info("AverageMedicareAllowedAmount--", object.AverageMedicareAllowedAmount)
		// log.Info("OutlierComprehensiveApcServices--", object.OutlierComprehensiveApcServices)
		// log.Info("################################################################################")
	}
	return nil
}

func (i *Importer) outpatientAPIToImporterOutpatient(outAPI OutpatientAPI) (ImportOutpatient, error) {
	// 2016 is the only year with this format for Outpatient
	year := 2016
	impOut := ImportOutpatient{}
	impOut.Apc.Code = util.StringToInt(outAPI.Apc)
	impOut.Apc.Name = &outAPI.ApcDescription
	impOut.Provider.Name = &outAPI.ProviderName
	impOut.Provider.Street = &outAPI.ProviderStreetAddress
	impOut.Provider.City = &outAPI.ProviderCity
	impOut.Provider.ProviderID = util.StringToInt(outAPI.ProviderID)
	impOut.Provider.State = &outAPI.ProviderState
	impOut.Provider.Zipcode = util.StringToInt(outAPI.ProviderZipCode)
	impOut.Provider.HrrDescription = &outAPI.ProviderHrr
	impOut.Outpatient.ServicesCount = util.StringToInt(outAPI.ComprehensiveApcServices)
	impOut.Outpatient.AvgEstSubmittedCharges = util.StringToFloat64(outAPI.AverageEstimatedTotalSubmittedCharges)
	impOut.Outpatient.AvgTotalPayments = util.StringToFloat64(outAPI.AverageMedicarePaymentAmount)
	impOut.Outpatient.Year = &year

	return impOut, nil
}
