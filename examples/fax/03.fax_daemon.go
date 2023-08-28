// Retarus GmbH ©2023
// The Fax Daemon is a meticulously engineered, production-ready application. While it's primed for real-world deployment,
// it also exemplifies the myriad applications and capabilities accessible through our platform. The daemon demonstrates
// the integration prowess of the 'retarus-go' SDK with our services, empowering developers to create high-value products and solutions.

// If you're testing or considering using this application:
// Set environment variables 'retarus_username', 'retarus_password', and 'retarus_cuno' with your credentials.
// This ensures the SDK can authenticate with our cloud communication services.
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/retarus/retarus-go/common"
	"github.com/retarus/retarus-go/fax"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var watcher *fsnotify.Watcher

func prepareFax(filePath string) (string, error) {
	// Read the entire file into a byte slice
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Encode the file content to base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

func writeJobReport(jobReport *fax.Report, outdir string) {
	jsonData, err := json.MarshalIndent(jobReport, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}
	if !strings.HasSuffix(outdir, "/") {
		outdir = outdir + "/"
	}
	storeString := outdir + jobReport.JobID + ".json"

	// Ensure the directory exists
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.Mkdir(outdir, 0755)
	}

	// Write the JSON data to a file in the specified directory
	err = ioutil.WriteFile(storeString, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("JSON data written successfully!")
}

func pullFaxReportWorker(jobChannel chan string, faxClient fax.Client, outdir string) {
	id, ok := <-jobChannel

	time.Sleep(3 * time.Second)
	if !ok {
		log.Println("Channel closed, worker terminating")
		return
	}
	isFinished := false
	res, err := faxClient.GetReport(id)
	// If job wasn't finished processing and got a 404, reschedule to fax means -> store it into channel
	if err != nil {
		isFinished = true
		log.Fatalf("Could not find fax report for jobid: %c", id)
	}

	if !isFinished {
		jobChannel <- id
		return
	}

	writeJobReport(res, outdir)
}

func start(inDir string, outDir string) {

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	for i, arg := range os.Args {
		fmt.Printf("os.Args[%d]: %s\n", i, arg)
	}

	config := fax.NewConfigFromEnv(common.Europe)

	faxClient := fax.NewClient(config)
	jobChan := make(chan string)

	go pullFaxReportWorker(jobChan, faxClient, outDir)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Create) {
					if strings.HasSuffix(event.Name, ".pdf") {
						pathSplitted := strings.Split(event.Name, "/")
						filename := pathSplitted[len(pathSplitted)-1]
						splitted := strings.Split(filename, ".")
						if len(splitted) != 2 {
							log.Fatal("Fax could not be send, naming schema was wrong.")
						}
						number := splitted[0]
						fax_document_data, err := prepareFax(event.Name)
						if err != nil {
							panic("Could not read fax file from path")
						}

						job := fax.NewJob()
						job.AddRecipient(fax.Recipient{Number: number})
						job.AddDocument(fax.Document{Name: "test.pdf", Reference: "fax_daemon_job", Data: fax_document_data})

						res, err := faxClient.Send(job)
						if err != nil {
							panic("Could not send fax")
						}
						// send job id to fax report fetcher
						jobChan <- res
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := watcher.Add(inDir)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "out",
				Value: "",
				Usage: "Destination where the fax reports should be stored.",
			},
			&cli.StringFlag{
				Name:  "in",
				Value: "",
				Usage: "Which folder destination should be watched for new fax pdfs.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			log.Println("Length, ", len(cCtx.FlagNames()))
			if len(cCtx.FlagNames()) != 4 {
				log.Println("Missing flag, might be --in or --out or both.")
				return nil
			}

			start(cCtx.FlagNames()[1], cCtx.FlagNames()[2])
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println("Error")
	}
}
