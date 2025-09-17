package services

import (
	"strconv"
	"strings"
)

type Pagination struct {
	TotalRecords   int   `json:"total_records"`
	PerPage        int   `json:"per_page"`
	CurrentPage    int   `json:"current_page"`
	TotalPages     int   `json:"total_pages"`
	Offset         int   `json:"offset"`
	PerPageOptions []int `json:"per_page_options"`
}

func (p *Pagination) Build() string {
	if p.PerPage <= 0 {
		p.PerPage = 10 // Default per-page value
	}
	if p.CurrentPage <= 0 {
		p.CurrentPage = 1
	}

	// Calculate total pages
	if p.TotalRecords > 0 {
		p.TotalPages = (p.TotalRecords + p.PerPage - 1) / p.PerPage
	} else {
		p.TotalPages = 1
	}

	// Calculate offset
	p.Offset = (p.CurrentPage - 1) * p.PerPage

	// Set default PerPage options if not provided
	if len(p.PerPageOptions) == 0 {
		p.PerPageOptions = []int{10, 25, 50, 100, 200}
	}

	// Build HTML
	var html strings.Builder

	// Dropdown for rows per page
	html.WriteString(`<div class="d-flex align-items-center">`)
	// html.WriteString(`<label class="me-2">Rows per page:</label>`)
	html.WriteString(`<select id="perPageSelect" class="form-select form-select me-1">`)
	for _, option := range p.PerPageOptions {
		selected := ""
		if option == p.PerPage {
			selected = ` selected`
		}
		html.WriteString(`<option value="` + strconv.Itoa(option) + `"` + selected + `>` + strconv.Itoa(option) + `</option>`)
	}
	html.WriteString(`</select>`)
	html.WriteString(`</div>`)

	// Pagination buttons
	html.WriteString(`<div class="d-flex justify-content-center">`)
	html.WriteString(`<button class="btn btn-sm btn-falcon-default me-1" type="button" title="Previous" onclick="pagination(` + strconv.Itoa(p.CurrentPage-1) + `)"`)

	if p.CurrentPage == 1 {
		html.WriteString(` disabled`)
	}
	html.WriteString(`><span class="fas fa-chevron-left"></span></button>`)

	html.WriteString(`<ul class="pagination mb-0">`)
	startPage, endPage := calculatePageRange(p.CurrentPage, p.TotalPages)

	for i := startPage; i <= endPage; i++ {
		if i == p.CurrentPage {
			html.WriteString(`<li class="page-item active " ><a class="page-link" style="border:none;  outline:none; background-color:#871A83;">` + strconv.Itoa(i) + `</a></li>`)
		} else {
			html.WriteString(`<li class="page-item" onclick="pagination(` + strconv.Itoa(i) + `)"><a class="page-link">` + strconv.Itoa(i) + `</a></li>`)
		}
	}

	html.WriteString(`</ul>`)

	html.WriteString(`<button class="btn btn-sm btn-falcon-default ms-1" type="button" title="Next" onclick="pagination(` + strconv.Itoa(p.CurrentPage+1) + `)"`)

	if p.CurrentPage == p.TotalPages {
		html.WriteString(` disabled`)
	}
	html.WriteString(`><span class="fas fa-chevron-right"></span></button>`)
	html.WriteString(`</div>`)

	return html.String()
}

// calculatePageRange determines the range of pagination buttons to display
func calculatePageRange(currentPage, totalPages int) (int, int) {
	startPage := currentPage - 1
	endPage := currentPage + 1

	if currentPage == 1 {
		startPage = 1
		endPage = 3
	} else if currentPage == totalPages {
		startPage = totalPages - 2
		endPage = totalPages
	}

	if startPage < 1 {
		startPage = 1
	}
	if endPage > totalPages {
		endPage = totalPages
	}

	return startPage, endPage
}
