package main

import (
	_ "embed"
	"encoding/csv"
	"io"
	"strings"

	"github.com/charmbracelet/log"
)

//go:embed Colleges.csv
var colleges []byte

// After colleges are parsed, place here.
// Allows for College list to be updated outside of Go
var records []CollegeRecord

type CollegeRecord struct {
	name   string
	domain string
	city   string
	state  string
}

func parseCSV() {
	list := csv.NewReader(strings.NewReader(string(colleges)))

	for {
		record, err := list.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Failed parsing CSV %v", err)
		}

		parsed := CollegeRecord{
			name:   record[0],
			domain: record[1],
			city:   record[2],
			state:  record[3],
		}

		records = append(records, parsed)
	}
}

// Given an domain (osu.edu), check for it's data in the records
func fetchRecord(domain string) *CollegeRecord {

	// use basic for loop instead of range to avoid copying each record
	for i := range records {
		if records[i].domain == domain {
			return &records[i]
		}
	}

	// Return pointer as it avoids copying (items in records don't change)
	return nil
}
