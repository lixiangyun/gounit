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

type Test struct {
	sync.Mutex
	name    string
	fun     func(*Test)
	success int
	falied  int
	ret     RETURN
	err     string
	log     []string
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
		fun:  fun.(func(*Test)),
	}

	s.test = append(s.test, test)
	log.Printf("Suite %s Add [%s]\r\n", s.name, name)
}

func (s *Suite) runTest() {

	log.Println("Suite :", s.name)

	for _, t := range s.test {
		t.fun(t)

		if t.ret == FAILED {

			var loginfo string

			for _, v := range t.log {
				loginfo += v
			}

			loginfo += t.err

			if loginfo != "" {
				log.Printf("Test : %s run failed!\r\n(log: %s)\r\n", t.name, loginfo)
			} else {
				log.Printf("Test : %s run failed!\r\n", t.name)
			}

		} else {
			log.Printf("Test : %s run paas!\r\n", t.name)
		}
	}
}

func (s *Suite) statTest() Stat {

	s.Lock()
	defer s.Unlock()

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

	s.Lock()
	defer s.Unlock()

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

func (t *Test) ASSERT_LOG(b bool, log ...string) {

	t.Lock()
	defer t.Unlock()

	if b == false {
		t.falied++
		t.ret = FAILED

		if len(log) != 0 {
			t.log = log
		}

		if bReocrdCallStack {
			t.err = string(debug.Stack())
		}

	} else {
		t.success++
	}
}

func (t *Test) ASSERT(b bool) {
	t.ASSERT_LOG(b)
}

func (t *Test) ASSERT_STRING(str1, str2 string) {
	if str1 != str2 {
		t.ASSERT_LOG(false, fmt.Sprintf("EQUAL_FAIL!(%s,%s)\r\n", str1, str2))
	} else {
		t.ASSERT_LOG(true)
	}
}
