package services

type Rates struct {
}

func NewRates() *Rates {
	return &Rates{}
}

func (u Rates) GetRates(base string, target string) (float64, error) {

	return 0, nil
}
