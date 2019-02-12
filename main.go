package main

import (
	"fmt"

	"github.com/bluehawk27/medcost-go/config"
	"github.com/bluehawk27/medcost-go/importer"
)

func main() {
	config.Init()
	i := importer.NewImporter()
	fmt.Println(i.ImportInpatientData("test-Medicare_Provider_Charge_Inpatient_DRG100_FY2012.csv"))
}
