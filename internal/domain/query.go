package domain

const defaultPage = 1
const defaultLimit = 20

type PaginationQuery struct {
	Page  int `form:"page" binding:"omitempty"`
	Limit int `form:"limit" binding:"omitempty,max=50"`
}

func (p *PaginationQuery) NormalizePagination() {
	if p.Page <= 0 {
		p.Page = defaultPage
	}
	if p.Limit <= 0 {
		p.Limit = defaultLimit
	}
}
