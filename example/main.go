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

func (hello) Test02(s *gounit.Suite) {

	log.Println("test run ....")

	s.ASSERT(true)
}

func (hello) Test03(s *gounit.Suite) {

	log.Println("test run ....")

	s.ASSERT(true)
}

type world int

func (world) Test01(s *gounit.Suite) {

	log.Println("test run ....")

	s.ASSERT(true)
}

func (world) Test02(s *gounit.Suite) {

	log.Println("test run ....")

	s.ASSERT(true)
}

func (world) Test03(s *gounit.Suite) {

	log.Println("test run ....")

	s.ASSERT(true)
}

func main() {

	var a hello
	var b world

	gounit.AddSuite(&a)
	gounit.AddSuite(&b)

	gounit.RunAllTests()

}
