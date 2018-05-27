package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type TrainLegs struct {
	XMLName xml.Name `xml:"TrainLegs"`
	//Trains map[string]Train `xml:"TrainLeg"`
	Trains []Train `xml:"TrainLeg"`
}

//func unmarshalTickets(byteValue []byte) (map[string]Train, error) {
func unmarshalTickets(byteValue []byte) ([]Train, error) {

	var tl TrainLegs
	err := xml.Unmarshal(byteValue, &tl)
	if err != nil {
		return nil, err
	}

	return tl.Trains, nil
}

func parseXmlFile(filename string) ([]byte, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	// read our opened xmlFile as a byte array.
	return ioutil.ReadAll(xmlFile)
}
