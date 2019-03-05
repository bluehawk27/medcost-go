package main

import (
	"context"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bluehawk27/medcost-go/config"
	"github.com/bluehawk27/medcost-go/importer"
)

func main() {
	config.Init()
	// done := make(chan string)
	log.SetReportCaller(true)
	ctx := context.Background()
	importFromCSV(ctx)
	importFromApi(ctx)
	// msg := <-done // Block until we receive a message on the channel
	// log.Info(msg)

}

func importFromCSV(ctx context.Context) {
	imp := importer.NewImporter()
	inpatientFiles, err := ioutil.ReadDir("./inpatient")
	if err != nil {
		log.Fatal(err)
	}

	outpatientFiles, err := ioutil.ReadDir("./outpatient")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range inpatientFiles {
		fname := "./inpatient/" + file.Name()
		log.Info("importing", fname)

		if err := imp.ImportInpatientData(ctx, fname); err != nil {
			log.Error(err)
		}
	}

	for _, file := range outpatientFiles {
		fname := "./outpatient/" + file.Name()
		log.Info("importing", fname)

		if err := imp.ImportOutpatientData(ctx, fname); err != nil {
			log.Error(err)
		}
	}

}

func importFromApi(ctx context.Context) {
	imp := importer.NewImporter()
	key := os.Getenv("CMS_KEY")

	url := "https://data.cms.gov/resource/di27-xeyq.json?$limit=50000&$$app_token=" + key
	imp.ImportOutpatientAPIData(ctx, url)
}
