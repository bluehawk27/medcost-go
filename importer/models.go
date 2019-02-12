package importer

import "github.com/bluehawk27/medcost-go/store"

type ImportInpatient struct {
	ProviderID int
	DrgID      int
	Inpatient  *store.Inpatient
}

type ImportOutpatient struct {
	ProviderID int
	ApcID      int
	Inpatient  *store.Inpatient
}
