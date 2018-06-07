package gounit

type Stat struct {
	Paased int
	Failed int
}

func (s1 *Stat) Add(s2 Stat) {
	s1.Failed += s2.Failed
	s1.Paased += s2.Paased
}

func (s *Stat) Total() int {
	return s.Paased + s.Failed
}
