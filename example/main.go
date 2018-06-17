package main

import (
	"github.com/lixiangyun/gounit"
)

type hello int

func (hello) Test01(s *gounit.Test) {
	s.ASSERT(true)
}

func (hello) Test02(s *gounit.Test) {
	s.ASSERT(false)
}

type world int

func (world) Test01(s *gounit.Test) {
	for i := 0; i < 100; i++ {
		s.ASSERT(true)
	}
}

func (world) Test02(s *gounit.Test) {
	s.ASSERT(true)
}

func (world) Test03(s *gounit.Test) {
	s.ASSERT_LOG(false, "Test03 run failed!")
}

func (world) Test04(s *gounit.Test) {
	s.ASSERT_STRING("abc", "123")
}

func main() {

	var a hello
	var b world

	gounit.RecordCallStack(true)

	gounit.AddSuite(&a)
	gounit.AddSuite(&b)

	gounit.RunAllTests()

}
