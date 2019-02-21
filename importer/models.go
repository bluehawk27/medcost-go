package importer

import "github.com/bluehawk27/medcost-go/store"

type ImportInpatient struct {
	Provider  store.Provider
	Drg       store.DRG
	Inpatient store.Inpatient
}

type ImportOutpatient struct {
	Provider   *store.Provider
	Apc        *store.APC
	Outpatient *store.Outpatient
}
