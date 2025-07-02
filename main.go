package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing_ibmmq/config"
	ibmq "testing_ibmmq/ibmq"
)

func main() {
	// Validate command line arguments
	if len(os.Args) < 3 {
		log.Fatal("Usage: program <xml-file-path>")
	}

	//Load config
	configFile := os.Args[2]

	var secretData string
	var err error

	if filepath.Ext(configFile) == ".json" {
		secretData, err = config.ReadConfigFile(configFile)

		if err != nil {
			return
		}
	} else {
		fmt.Println("invalid file format")
		return
	}

	if !config.LoadConfig(secretData) {
		return
	}

	fmt.Println(os.Args[1])

	// Process XML file
	order := processXMLFile(os.Args[1])

	fmt.Println(order)

	xmlBytes, err := xml.Marshal(order)
	if err != nil {
		log.Fatal("Failed to marshal XML: ", err)
	}

	// Create IBM MQ instance
	var ibmq = ibmq.IBMQ{}

	// Ensure cleanup happens when function exits
	defer func() {
		if err := ibmq.Close(); err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
	}()

	// Connect to queue manager
	if !ibmq.ConnectToQueueManager() {
		log.Fatal("Failed to connect to queue manager")
	}

	// Put message into queue
	if err := ibmq.PutMessageIntoQueue(xmlBytes); err != nil {
		log.Fatal("Failed to PUT message: ", err)
	}
	log.Println("Message successfully put into queue")

	// Get message from queue
	msg, err := ibmq.GetMessageFromQueue()
	if err != nil {
		log.Fatal("Failed to GET message: ", err)
	}

	// Process the received message
	processMessage(msg)

	// if err := ibmq.Close(); err != nil {
	// 	log.Printf("Error during cleanup: %v", err)
	// }
}
