package utils

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type PagingRequest struct {
	Page   int         `default:"1"`
	Search string      `default:""`
	Limit  int         `default:"10"`
	Offset int         `default:"0"`
	Order  string      `default:"id"`
	Custom interface{} `default:""`
}

type PagingResponse struct {
	Total         int         `json:"total"`
	TotalFiltered int         `json:"total_filtered"`
	Error         string      `json:"error"`
	Status        int         `default:"200" json:"status"`
	Messages      string      `default:"Success" json:"message"`
	Data          interface{} `default:"[]" json:"data"`
	Search        string      `default:"" json:"search"`
	Next          bool        `default:"false" json:"next"`
	Back          bool        `default:"false" json:"back"`
	Limit         int         `default:"10" json:"limit"`
	Offset        int         `default:"0" json:"offset"`
	TotalPage     int         `default:"0" json:"total_page"`
	CurrentPage   int         `default:"1" json:"current_page"`
	Order         string      `default:"id" json:"order"`
}

func PopulatePaging(c echo.Context, custom string) (param PagingRequest) {
	customval := c.QueryParam(custom)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 && offset == 0 {
		page = 1
		offset = 0
	}
	if page >= 1 && offset == 0 {
		offset = (page - 1) * limit
	}
	draw, _ := strconv.Atoi(c.QueryParam("draw"))
	if draw == 0 {
		draw = 1
	}
	sort := c.QueryParam("sort")
	if strings.ToLower(sort) == "asc" {
		sort = "ASC"
	} else {
		sort = "DESC"
	}
	order := c.QueryParam("order")
	if order == "" {
		order = "created_at " + order
	} else {
		order = order + " " + sort + ", created_at " + sort
	}
	param = PagingRequest{
		Search: c.QueryParam("search"),
		Limit:  limit,
		Offset: offset,
		Order:  order,
		Custom: customval,
		Page:   page}
	return
}

func PopulateResPaging(param *PagingRequest, data interface{}, totalResult int64, totalFiltered int64) (output PagingResponse) {
	totalPages := int(totalFiltered) / param.Limit
	if int(totalFiltered)%param.Limit > 0 {
		totalPages++
	}

	currentPage := param.Offset/param.Limit + 1
	next := false
	back := false
	if currentPage < totalPages {
		next = true
	}
	if currentPage <= totalPages && currentPage > 1 {
		back = true
	}

	output = PagingResponse{
		Status:        200,
		Data:          data,
		Search:        param.Search,
		Order:         param.Order,
		Limit:         param.Limit,
		Offset:        param.Offset,
		Next:          next,
		Back:          back,
		Total:         int(totalResult),
		TotalFiltered: int(totalFiltered),
		CurrentPage:   currentPage,
		TotalPage:     totalPages,
	}
	return
}
