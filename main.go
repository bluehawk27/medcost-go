package main

import (
	"context"
	"fmt"

	"github.com/bluehawk27/medcost-go/config"
	"github.com/bluehawk27/medcost-go/importer"
)

func main() {
	config.Init()
	i := importer.NewImporter()
	ctx := context.Background()
	fmt.Println(i.ImportInpatientData(ctx, "test-Medicare_Provider_Charge_Inpatient_DRG100_FY2012.csv"))
}
