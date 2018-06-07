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
	sync.Mutex
	name    string
	fun     func(*Test)
	success int
	falied  int
	ret     RETURN
	err     []string
}

type Suite struct {
	sync.Mutex
	name string
	test []*Test
}

func (s *Suite) addTest(name string, fun interface{}) {
	test := &Test{
		name: name,
		ret:  SUCCESS,
		err:  make([]string, 0),
		fun:  fun.(func(*Test)),
	}

	s.test = append(s.test, test)
	log.Printf("Suite %s Add [%s]\r\n", s.name, name)
}

func (s *Suite) run() {
	log.Println("Suite :", s.name)

	for _, t := range s.test {
		t.fun(t)

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
		if funtype.In(0).String() != "*gounit.Test" {
			continue
		}

		s.addTest(vtype.Method(i).Name, funvalue.Interface())
	}

	return s
}

func (t *Test) ASSERT(b bool) {
	if b == false {
		t.falied++
		t.ret = FAILED
		t.err = append(t.err, string(debug.Stack()))
	} else {
		t.success++
	}
}

func (r *Result) Printf() string {
	s := fmt.Sprintf("Run Summary : Total   Passed    Failed\r\n")
	s += fmt.Sprintf("      Suite   %5d   %5d   %5d\r\n", r.Suite.Total(), r.Suite.Paased, r.Suite.Failed)
	s += fmt.Sprintf("      Tests   %5d   %5d   %5d\r\n", r.Tests.Total(), r.Tests.Paased, r.Tests.Failed)
	s += fmt.Sprintf("      Cases   %5d   %5d   %5d\r\n", r.Cases.Total(), r.Cases.Paased, r.Cases.Failed)
	return s
}

var gSuite []*Suite

func init() {
	gSuite = make([]*Suite, 0)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
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
