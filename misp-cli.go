/*
CLI to interact with your favorite MISP instance

Copyright (C) 2017  RawSec SARL (0xrawsec)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/0xrawsec/golang-misp/misp"
	"github.com/0xrawsec/golang-utils/fsutil"
	"github.com/0xrawsec/golang-utils/log"
)

const (
	// ExitSuccess RC
	ExitSuccess = 0
	// ExitFail RC
	ExitFail = 1
	// Version String
	Version   = "MISP CLI 1.0"
	Copyright = "MISP CLI  Copyright (C) 2017 RawSec SARL (@0xrawsec)"
	License   = `License GPLv3: This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain
conditions;`
)

var (
	reValidMispDate = regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}$")

	// mispConfig file
	mispConfig, _ = fsutil.AbsFromRelativeToBin("config.json")

	// debug flag
	debug   bool
	version bool

	attributeSearchFlag bool
	eventSearchFlag     bool
	// Arguments to pass to Misp Query can be used for both attributes and events
	qValue    string
	qType     string
	qCategory string
	qOrg      string
	qTags     string
	qFrom     string
	qTo       string
	qLast     string
	qEventID  string

	// MISP related structures
	mispCon   misp.MispCon
	mispQuery misp.MispQuery
)

func isFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.Mode().IsRegular()
}

func toJSON(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

func isValidMispDate(date string) bool {
	return reValidMispDate.MatchString(date)
}

func main() {
	flag.BoolVar(&debug, "d", debug, "Enable debugging")
	flag.BoolVar(&eventSearchFlag, "e", eventSearchFlag, "Flag to search for events")
	flag.BoolVar(&attributeSearchFlag, "a", attributeSearchFlag, "Flag to search for attributes")
	flag.BoolVar(&version, "version", version, "Print version information")
	flag.StringVar(&mispConfig, "c", mispConfig, "Configuration file to connect to MISP")
	flag.StringVar(&qValue, "v", qValue, "Value to search for")
	flag.StringVar(&qLast, "l", qLast, "Last event query")
	flag.StringVar(&qFrom, "from", qFrom, "Query events from date")
	flag.StringVar(&qTo, "to", qTo, "Query events until 'to' parameter")
	flag.StringVar(&qCategory, "cat", qCategory, "Category to query")
	flag.StringVar(&qType, "type", qType, "Type argument for query")
	flag.StringVar(&qTags, "tags", qTags, "Tags argument for query")
	flag.StringVar(&qOrg, "org", qOrg, "Organisation")
	flag.StringVar(&qEventID, "eventid", qEventID, "Event ID to look for")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: %[1]s [OPTIONS]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	// Initialization and input validation

	if debug {
		log.InitLogger(log.LDebug)
	}

	if version {
		fmt.Fprintf(os.Stderr, "%s\n%s\n%s\n", Version, Copyright, License)
		return
	}

	if attributeSearchFlag == false && eventSearchFlag == false {
		log.LogErrorAndExit(fmt.Errorf("Specify either attribute or event to search for"), ExitFail)
	}

	if attributeSearchFlag == true && eventSearchFlag == true {
		log.LogErrorAndExit(fmt.Errorf("Cannot query for both attributes and events"), ExitFail)
	}

	if !isFile(mispConfig) {
		log.LogErrorAndExit(fmt.Errorf("Configuration file not found: %s", mispConfig), ExitFail)
	}

	if qFrom != "" && !isValidMispDate(qFrom) {
		log.LogErrorAndExit(fmt.Errorf("From parameter expects format like YYYY-MM-DD"), ExitFail)
	}

	if qTo != "" && !isValidMispDate(qTo) {
		log.LogErrorAndExit(fmt.Errorf("To parameter expects format like YYYY-MM-DD"), ExitFail)
	}

	// Load the configuration file
	c := misp.LoadConfig(mispConfig)
	// Start a new MISP connection
	mispCon = misp.NewCon(c.Proto, c.Host, c.APIKey, c.APIURL)

	// Start with further processing

	////// If it is Attribute Search
	if attributeSearchFlag {
		mispQuery = misp.MispAttributeQuery{
			Value:    qValue,
			Type:     qType,
			Category: qCategory,
			Org:      qOrg,
			Tags:     qTags,
			From:     qFrom,
			To:       qTo,
			Last:     qLast,
			EventID:  qEventID,
		}
	}

	///// If it is Event Search
	if eventSearchFlag {
		mispQuery = misp.MispEventQuery{
			Value:    qValue,
			Type:     qType,
			Category: qCategory,
			Org:      qOrg,
			Tags:     qTags,
			From:     qFrom,
			To:       qTo,
			Last:     qLast,
			EventID:  qEventID,
		}
	}

	matches, err := mispCon.Search(mispQuery)
	if err != nil {
		log.LogErrorAndExit(err, ExitFail)
	}
	for match := range matches.Iter() {
		fmt.Println(string(toJSON(match)))
	}
}
