package importer

import "github.com/bluehawk27/medcost-go/store"

type ImportInpatient struct {
	Provider  store.Provider
	Drg       store.DRG
	Inpatient store.Inpatient
}

type ImportOutpatient struct {
	Provider   store.Provider
	Apc        store.APC
	Outpatient store.Outpatient
}

// https://data.cms.gov/resource/j94m-eswy.json
type OutpatientAPI struct {
	Apc                                   string `json:"apc"`
	ApcDescription                        string `json:"apc_description"`
	AverageEstimatedTotalSubmittedCharges string `json:"average_estimated_total_submitted_charges"`
	AverageMedicareAllowedAmount          string `json:"average_medicare_allowed_amount"`
	AverageMedicarePaymentAmount          string `json:"average_medicare_payment_amount"`
	ComprehensiveApcServices              string `json:"comprehensive_apc_services"`
	OutlierComprehensiveApcServices       string `json:"outlier_comprehensive_apc_services"`
	ProviderCity                          string `json:"provider_city"`
	ProviderHrr                           string `json:"provider_hrr"`
	ProviderID                            string `json:"provider_id"`
	ProviderName                          string `json:"provider_name"`
	ProviderState                         string `json:"provider_state"`
	ProviderStreetAddress                 string `json:"provider_street_address"`
	ProviderZipCode                       string `json:"provider_zip_code"`
}

type OutpatientAPISlice struct {
	OutpatientAPIsrvcs []OutpatientAPI
}
