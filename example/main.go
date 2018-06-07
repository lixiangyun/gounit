package main

import (
	"github.com/lixiangyun/gounit"
)

type hello int

func (hello) Test01(s *gounit.Test) {
	s.ASSERT(true)
}

func (hello) Test02(s *gounit.Test) {
	s.ASSERT(true)
}

type world int

func (world) Test01(s *gounit.Test) {
	for i := 0; i < 100; i++ {
		s.ASSERT(true)
	}
}

func (world) Test02(s *gounit.Test) {
	s.ASSERT(false)
	s.ASSERT(false)
}

func main() {

	var a hello
	var b world

	gounit.AddSuite(&a)
	gounit.AddSuite(&b)

	gounit.RunAllTests()

}
