package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type OrderStruct struct {
	Book []struct {
		ID        string `xml:"id,attr"`
		Title     string `xml:"title"`
		Author    string `xml:"author"`
		Genre     string `xml:"genre"`
		Published string `xml:"published"`
		Price     struct {
			Currency string `xml:"currency,attr"`
			Value    string `xml:",chardata"`
		} `xml:"price"`
	} `xml:"book"`
}

func processXMLFile(filepath string) OrderStruct {
	data, err := os.ReadFile(filepath)

	var order OrderStruct

	if err != nil {
		log.Fatal("error reading file: ", err)
	}

	err = xml.Unmarshal(data, &order)

	fmt.Println(order)

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
