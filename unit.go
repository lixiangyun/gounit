package gounit

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
)

type RETURN int

const (
	_ RETURN = iota
	SUCCESS
	FAILED
)

type Result struct {
	Tests Stat
	Suite Stat
	Cases Stat
}

type Test struct {
	name    string
	fun     func(*Suite)
	success int
	falied  int
	ret     RETURN
	err     []string
}

type Suite struct {
	sync.Mutex
	name string
	cur  *Test
	test []*Test
}

func (s *Suite) addTest(name string, fun interface{}) {

	test := &Test{name: name, ret: SUCCESS, err: make([]string, 0), fun: fun.(func(*Suite))}

	s.test = append(s.test, test)

	log.Println("Suite "+s.name+" add test ", name, fun)
}

func (s *Suite) run() {

	log.Println("Suite :", s.name)

	for _, t := range s.test {
		s.cur = t
		t.fun(s)
		if t.ret == FAILED {
			log.Printf("Test : %s run failed!(%s)\r\n", t.name, t.err)
		} else {
			log.Printf("Test : %s run paas!\r\n", t.name)
		}
	}
}

func (s *Suite) statTest() Stat {
	var stat Stat
	for _, t := range s.test {
		if t.ret == SUCCESS {
			stat.Paased++
		} else {
			stat.Failed++
		}
	}
	return stat
}

func (s *Suite) statCase() Stat {
	var stat Stat
	for _, t := range s.test {
		stat.Paased = stat.Paased + t.success
		stat.Failed = stat.Failed + t.falied
	}
	return stat
}

func (r *Result) Printf() string {
	s := fmt.Sprintf("Run Summary : Total   Passed    Failed\r\n")
	s += fmt.Sprintf("      Suite   %5d   %5d   %5d\r\n", r.Suite.Total(), r.Suite.Paased, r.Suite.Failed)
	s += fmt.Sprintf("      Tests   %5d   %5d   %5d\r\n", r.Tests.Total(), r.Tests.Paased, r.Tests.Failed)
	s += fmt.Sprintf("      Cases   %5d   %5d   %5d\r\n", r.Cases.Total(), r.Cases.Paased, r.Cases.Failed)
	return s
}

func newSuite(this interface{}) *Suite {

	vfun := reflect.ValueOf(this)
	vtype := vfun.Type()
	num := vfun.NumMethod()

	s := &Suite{name: vtype.String(), test: make([]*Test, 0)}

	for i := 0; i < num; i++ {

		funvalue := vfun.Method(i)
		funtype := funvalue.Type()

		if funtype.NumIn() != 1 || funtype.NumOut() != 0 {
			continue
		}
		if funtype.In(0).String() != "*gounit.Suite" {
			continue
		}

		s.addTest(vtype.Method(i).Name, funvalue.Interface())
	}

	return s
}

func (s *Suite) ASSERT(b bool) {
	if b == false {
		s.cur.falied++
		s.cur.ret = FAILED
		s.cur.err = append(s.cur.err, string(debug.Stack()))
	} else {
		s.cur.success++
	}
}

var gSuite []*Suite

func init() {
	gSuite = make([]*Suite, 0)
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
		s.run()
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
