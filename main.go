package main

import (
	"encoding/xml"
	"log"
	"os"
	ibmq "testing_ibmmq/IBMQ"
)

func main() {
	// Validate command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: program <xml-file-path>")
	}

	// Process XML file
	order := processXMLFile(os.Args[1])
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
}
