package main

import (
	"log"

	"github.com/lixiangyun/gounit"
)

type hello int

func (hello) Test01(s *gounit.Suite) {

	log.Println("test run ....")

	s.ASSERT(true)
}

func main() {

	var a hello

	gounit.AddSuite(&a)

	gounit.RunAllTests()

}
