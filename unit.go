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
	success int
	falied  int
	total   int
}

type Test struct {
	name string
	fun  func(*Suite)
	ret  RETURN
	err  string
}

type Suite struct {
	sync.Mutex
	name string
	cur  *Test
	test []*Test
}

func (s *Suite) addTest(name string, fun interface{}) {

	test := &Test{name: name, fun: fun.(func(*Suite))}
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

func (s *Suite) stat() Result {
	var result Result

	for _, t := range s.test {
		if t.ret == SUCCESS {
			result.success++
		} else {
			result.falied++
		}
		result.total++
	}
	return result
}

func (r *Result) Add(a Result) {
	r.falied += a.falied
	r.success += a.success
	r.total += a.total
}

func (r *Result) Printf() string {

	s := fmt.Sprintf("Run Summary : Total   Passed    Failed\r\n")
	s += fmt.Sprintf("              %5d   %5d   %5d\r\n", r.total, r.success, r.falied)

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
	if b == true {
		s.cur.ret = SUCCESS
	} else {
		s.cur.ret = FAILED
		s.cur.err = string(debug.Stack())
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
		result.Add(s.stat())
	}

	fmt.Println(result.Printf())

	return result
}
