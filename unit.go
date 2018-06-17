package gounit

import (
	"fmt"
	"log"
)

type Result struct {
	Tests Stat `json:"Tests"`
	Suite Stat `json:"Suite"`
	Cases Stat `json:"Cases"`
}

func (r *Result) Printf() string {
	s := fmt.Sprintf("Run Summary : Total   Passed    Failed\r\n")
	s += fmt.Sprintf("      Suite   %5d   %5d   %5d\r\n", r.Suite.Total(), r.Suite.Paased, r.Suite.Failed)
	s += fmt.Sprintf("      Tests   %5d   %5d   %5d\r\n", r.Tests.Total(), r.Tests.Paased, r.Tests.Failed)
	s += fmt.Sprintf("      Cases   %5d   %5d   %5d\r\n", r.Cases.Total(), r.Cases.Paased, r.Cases.Failed)
	return s
}

var bReocrdCallStack = false

var gSuite []*Suite

func init() {
	gSuite = make([]*Suite, 0)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

func RecordCallStack(enable bool) {
	bReocrdCallStack = enable
}

func AddSuite(this interface{}) {
	s := newSuite(this)
	if s != nil {
		gSuite = append(gSuite, s)
	}
}

func RunAllTests() Result {

	var result Result

	for _, s := range gSuite {
		s.runTest()
	}

	for _, s := range gSuite {

		ts := s.statTest()
		cs := s.statCase()

		result.Tests.Add(ts)
		result.Cases.Add(cs)

		if ts.Failed > 0 {
			result.Suite.Failed++
		} else {
			result.Suite.Paased++
		}
	}

	fmt.Println(result.Printf())
	return result
}
