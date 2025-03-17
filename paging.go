package g_learning_connector

import (
	"fmt"
)

// PageInfo: struct untuk menyimpan informasi halaman yang sedang ditampilkan.
// UI example:
// Menampilkan 10 dari 50 data < 1 2 (3) 4 5 > ➡️ Menampilkan [TotalDataInCurrentPage] dari [TotalData] data
// Menampilkan 20 sampai 30 dari 50 data < 1 2 (3) 4 5 > ➡️ Menampilkan [from] sampai [to] dari [total_data] data

type PageInfo struct {
	// flag to check if there is a previous page
	HasPreviousPage bool `json:"has_previous_page"`

	// flag to check if there is a next page
	HasNextPage bool `json:"has_next_page"`

	// the current page
	CurrentPage int64 `json:"current_page"`

	// the total data for each page
	PerPage int64 `json:"per_page"`

	// the total data
	TotalData int64 `json:"total_data"`

	// the total page (we can use this to last page number)
	LastPage int64 `json:"last_page"`

	// first number in this page
	From int64 `json:"from"`

	// last number in this page
	To int64 `json:"to"`

	// total data in this page
	TotalDataInCurrentPage int64 `json:"total_data_in_current_page"`

	// ui styles
	Style1 string `json:"style1"`
	Style2 string `json:"style2"`
}

// NewPageInfo membuat objek PageInfo baru berdasarkan informasi yang diberikan.
func NewPageInfo(
	currentPage,
	perPage,
	offset,
	totalData int64,
) (*PageInfo, error) {
	lastPage := totalData / perPage

	// pastikan ketika totalData tidak habis dibagi perPage maka perlu ditambah 1
	if totalData%perPage != 0 {
		lastPage++
	}

	// atur nilai from dan to
	to := offset + perPage
	if to > totalData {
		to = totalData
	}

	from := offset + 1

	// hitung total data dari from hingga to
	totalDataInCurrentPage := to - offset

	if totalDataInCurrentPage == 0 {
		from = 0
	}

	// pastikan currentPage tidak lebih besar dari lastPage
	if currentPage > lastPage {
		currentPage = lastPage
	}

	// generate style
	style1 := fmt.Sprintf("Menampilkan %d dari %d data", from, totalData)
	style2 := fmt.Sprintf("Menampilkan %d sampai %d dari %d data", from, to, totalData)

	return &PageInfo{
		HasPreviousPage:        currentPage > 1,
		HasNextPage:            currentPage < lastPage,
		CurrentPage:            currentPage,
		PerPage:                perPage,
		TotalData:              totalData,
		LastPage:               lastPage,
		From:                   from,
		To:                     to,
		TotalDataInCurrentPage: totalDataInCurrentPage,
		Style1:                 style1,
		Style2:                 style2,
	}, nil
}
