package store

// Provider Represents the medical providers
type Provider struct {
	ID             int      `json:"id" db:"id"`
	ProviderID     *int     `json:"provider_id" db:"provider_id"`
	Name           *string  `json:"name" db:"name"`
	Street         *string  `json:"street" db:"street"`
	City           *string  `json:"city" db:"city"`
	Zipcode        *int     `json:"zip_code" db:"zip_code"`
	State          *string  `json:"state" db:"state"`
	HrrDescription *int     `json:"HRR_description" db:"hrr_description"`
	CreatedAt      *string  `json:"created_at" db:"created_at"`
	UpdatedAt      *string  `json:"updated_at" db:"updated_at"`
	Latitude       *float64 `json:"latitude" db:"lat"`
	Longitude      *float64 `json:"longitude" db:"long"`
}

type DRG struct {
	ID   int     `json:"id" db:"id"`
	Code *int    `json:"code" db:"code"`
	Name *string ` json:"name" db:"name"`
}

type APC struct {
	ID   int     `json:"id" db:"id"`
	Code *int    `json:"code" db:"code"`
	Name *string `json:"name" db:"name"`
}

type Outpatient struct {
	ID                     int      `json:"id" db:"id"`
	ProviderID             *int     `json:"provider_id" db:"provider_id"`
	ApcID                  *int     `json:"ambulatory_payment_classification_id" db:"ambulatory_payment_classification_id"`
	ServicesCount          *int     `json:"services_count" db:"services_count"`
	AvgEstSubmittedCharges *float64 `json:"avg_est_submitted_charges" db:"avg_est_submitted_charges"`
	AvgTotalPayments       *float64 `json:"avg_total_payments" db:"avg_total_payments"`
	Year                   *int     `json:"year" db:"year"`
	CreatedAt              *string  `json:"created_at" db:"created_at"`
	UpdatedAt              *string  `json:"updated_at" db:"updated_at"`
}

type Inpatient struct {
	ID                 int      `json:"id" db:"id"`
	ProviderID         *int     `json:"provider_id" db:"provider_id"`
	DrgID              *int     `json:"diagnosis_related_group_id" db:"drg_id"`
	TotalDiscarges     *int     `json:"total_discharges" db:"total_discharges"`
	AvgCoveredCharges  *float64 `json:"avg_covered_discharges" db:"avg_covered_discharges"`
	AvgTotalPayments   *float64 `json:"avg_total_payments" db:"avg_total_payments"`
	AvgMedicarePayment *float64 `json:"avg_medicare_payments" db:"avg_medicare_payments"`
	Year               *int     `json:"year" db:"year"`
	CreatedAt          *string  `json:"created_at" db:"created_at"`
	UpdatedAt          *string  `json:"updated_at" db:"updated_at"`
}
