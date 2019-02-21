package main

import (
	"context"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/bluehawk27/medcost-go/config"
	"github.com/bluehawk27/medcost-go/importer"
)

func main() {
	config.Init()
	done := make(chan string)
	log.SetReportCaller(true)
	ctx := context.Background()
	go importFromCSV(ctx, done)
	msg := <-done // Block until we receive a message on the channel
	log.Info(msg)

}

func importFromCSV(ctx context.Context, done chan string) {
	files, err := ioutil.ReadDir("./inpatient")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fname := "./inpatient/" + file.Name()
		log.Info("importing", fname)

		if err := importer.NewImporter().ImportInpatientData(ctx, fname); err != nil {
			log.Error(err)
			done <- "I finished importing " + fname + " with errors" + err.Error()
		}
		done <- "Successfully finished importing " + fname
	}
}
