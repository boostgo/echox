package echox

import "github.com/boostgo/pagex"

type PageParams struct {
	Page int   `query:"page" form:"page" default:"1"`
	Size int64 `query:"page-size" form:"page-size" default:"20"`
}

func (p PageParams) Pagination() pagex.Pagination {
	return pagex.Pagination{
		Page: p.Page,
		Size: p.Size,
	}
}

func (p PageParams) MaxPages(count int64) int64 {
	return pagex.MaxPages(p.Size, count)
}
