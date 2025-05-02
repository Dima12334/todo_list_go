package domain

import "strings"

const defaultPage = 1
const defaultLimit = 20

type PaginationQuery struct {
	Page   int `form:"page" binding:"omitempty"`
	Limit  int `form:"limit" binding:"omitempty,max=50"`
	Offset int
}

func (p *PaginationQuery) NormalizePagination() {
	if p.Page <= 0 {
		p.Page = defaultPage
	}
	if p.Limit <= 0 {
		p.Limit = defaultLimit
	}

	p.Offset = (p.Page - 1) * p.Limit
}

type TaskFiltersQuery struct {
	CreatedAtDateFrom string   `form:"createdAtDateFrom"`
	CreatedAtDateTo   string   `form:"createdAtDateTo"`
	Completed         *bool    `form:"completed"`
	CategoryIDs       []string `form:"categoryIds"`
}

func (f *TaskFiltersQuery) NormalizeFilters() {
	if len(f.CategoryIDs) == 1 {
		raw := strings.Split(f.CategoryIDs[0], ",")
		f.CategoryIDs = make([]string, 0, len(raw))
		for i := range raw {
			trimmed := strings.TrimSpace(raw[i])
			if trimmed != "" {
				f.CategoryIDs = append(f.CategoryIDs, trimmed)
			}
		}
	}
}

type GetTasksQuery struct {
	PaginationQuery
	TaskFiltersQuery
}
