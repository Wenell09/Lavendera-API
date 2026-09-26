package dto

type OutletFilter struct {
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

func (o *OutletFilter) SetDefault() {
	if o.Page < 1 {
		o.Page = 1
	}
	if o.Limit < 1 {
		o.Limit = 10
	}
	if o.Limit > 100 {
		o.Limit = 100
	}
}
