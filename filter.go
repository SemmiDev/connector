package g_learning_connector

import (
	"strings"
)

const (
	AscOrder           string = "ASC"
	DescOrder          string = "DESC"
	DefaultPageLimit   int64  = 50
	DefaultCurrentPage int64  = 1
	DefaultColumnOrder string = AscOrder
)

type Filter struct {
	CurrentPage int64  `json:"current_page" form:"current_page" query:"current_page"` // page (berpindah-pindah halaman)
	PerPage     int64  `json:"per_page" form:"per_page" query:"per_page"`             // limit (batas data yang ditampilkan)
	Keyword     string `json:"keyword" form:"keyword" query:"keyword"`                // search keyword (keyword pencarian)
	SortBy      string `json:"sort_by" form:"sort_by" query:"sort_by"`                // column name to sort
	Order       string `json:"order" form:"order" query:"order"`                      // asc or desc order
}

func NewFilterPagination() Filter {
	return Filter{
		CurrentPage: DefaultCurrentPage,
		PerPage:     DefaultPageLimit,
		Order:       DefaultColumnOrder,
	}
}

func (f *Filter) GetLimit() int64 {
	return f.PerPage
}

func (f *Filter) GetOffset() int64 {
	offset := (f.CurrentPage - 1) * f.PerPage // example: (1 - 1) * 10 = 0, (2 - 1) * 10 = 10
	return offset
}

func (f *Filter) HasKeyword() bool {
	return f.Keyword != ""
}

func (f *Filter) HasSort() bool {
	return f.SortBy != ""
}

func (f *Filter) IsDesc() bool {
	return strings.ToUpper(f.Order) == DescOrder
}
