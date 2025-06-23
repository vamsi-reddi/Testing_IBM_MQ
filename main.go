package main

import (
	"encoding/xml"
	"log"
	"os"
	ibmq "testing_ibmmq/IBMQ"
)

func main(){
	order := processXMLFile(os.Args[1])

	xmlBytes,_ := xml.Marshal(order)

	var ibmq = ibmq.IBMQ{}

	if !ibmq.ConnectToQueueManager(){
		log.Fatal("failed to Connect to Queue manager")
	}

	err := ibmq.PutMessageIntoQueue(xmlBytes)

	if err != nil{
		log.Fatal("failed to PUT : ", err)
	}

	msg, err := ibmq.GetMessageFromQueue()

	if err != nil {
		log.Fatal("failed to GET : ", err)
	}

	processMessage(msg)

}