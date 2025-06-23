package main

import (
	"encoding/xml"
	"log"
	"os"
)

type OrderStruct struct {
	Order  struct {
		ID string `xml:"id"`
		Name string `xml:"name"`
		Quantity int `xml:"quantity"`
	} 
}

func processXMLFile(filepath string) OrderStruct{
	data , err := os.ReadFile(filepath)

	var order OrderStruct

	if err != nil {
		log.Fatal("error reading file: ", err)
	}

	err = xml.Unmarshal(data, &order)

	if err != nil {
		log.Fatal("error parsing xml to object : ", err)
	}

	return order
}

func processMessage(msg string) {
    var order OrderStruct
    err := xml.Unmarshal([]byte(msg), &order)
    if err != nil {
        log.Println("Error unmarshaling:", err)
        return
    }
    log.Printf("Order received: %+v\n", order)
}