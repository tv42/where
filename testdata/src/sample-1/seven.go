package sample

type seven struct {
	a sevenMore
}

func (s *seven) method() {}

type sevenMore struct {
}

func (s *sevenMore) method() {}
