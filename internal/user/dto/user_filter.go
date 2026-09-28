package dto

type UserFilter struct {
	Search string `json:"search"`
	Role   string `json:"role"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

func (f *UserFilter) SetDefault() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 10
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
}
