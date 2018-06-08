package gounit

type Stat struct {
	Paased int `json:"Paased"`
	Failed int `json:"Failed"`
}

func (s1 *Stat) Add(s2 Stat) {
	s1.Failed += s2.Failed
	s1.Paased += s2.Paased
}

func (s *Stat) Total() int {
	return s.Paased + s.Failed
}
