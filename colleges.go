package main

import (
	_ "embed"
	"encoding/csv"
	"io"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
)

//go:embed colleges.csv
var colleges []byte

// Ivory League Colleges
var IvoryLeague = []string{
	"brown.edu",
	"columbia.edu",
	"cornell.edu",
	"dartmouth.edu",
	"harvard.edu",
	"upenn.edu",
	"princeton.edu",
	"yale.edu",
}

// State Abbr to name
var states = map[string]string{
	"AL": "Alabama",
	"AK": "Alaska",
	"AZ": "Arizona",
	"AR": "Arkansas",
	"CA": "California",
	"CO": "Colorado",
	"CT": "Connecticut",
	"DE": "Delaware",
	"FL": "Florida",
	"GA": "Georgia",
	"HI": "Hawaii",
	"ID": "Idaho",
	"IL": "Illinois",
	"IN": "Indiana",
	"IA": "Iowa",
	"KS": "Kansas",
	"KY": "Kentucky",
	"LA": "Louisiana",
	"ME": "Maine",
	"MD": "Maryland",
	"MA": "Massachusetts",
	"MI": "Michigan",
	"MN": "Minnesota",
	"MS": "Mississippi",
	"MO": "Missouri",
	"MT": "Montana",
	"NE": "Nebraska",
	"NV": "Nevada",
	"NH": "New Hampshire",
	"NJ": "New Jersey",
	"NM": "New Mexico",
	"NY": "New York",
	"NC": "North Carolina",
	"ND": "North Dakota",
	"OH": "Ohio",
	"OK": "Oklahoma",
	"OR": "Oregon",
	"PA": "Pennsylvania",
	"RI": "Rhode Island",
	"SC": "South Carolina",
	"SD": "South Dakota",
	"TN": "Tennessee",
	"TX": "Texas",
	"UT": "Utah",
	"VT": "Vermont",
	"VA": "Virginia",
	"WA": "Washington",
	"WV": "West Virginia",
	"WI": "Wisconsin",
	"WY": "Wyoming",
}

// After colleges are parsed, place here.
// Allows for College list to be updated outside of Go
var records []CollegeRecord

type CollegeRecord struct {
	name      string
	domain    string
	stateAbbr string
	isIvory   bool
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
			name:      record[0],
			domain:    record[1],
			stateAbbr: record[2],
			isIvory:   false,
		}

		if slices.Contains(IvoryLeague, parsed.domain) {
			parsed.isIvory = true
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
