package utils

import (
	"strconv"
	"strings"

	"github.com/guregu/null"
	"github.com/labstack/echo/v4"
)

type ReqPaging struct {
	Page   int         `default:"1"`
	Search string      `default:""`
	Limit  int         `default:"10"`
	Offset int         `default:"0"`
	Sort   string      `default:"ASC"`
	Order  string      `default:"id"`
	Custom interface{} `default:""`
}

type ResPaging struct {
	TotalData       int         `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
	Error           string      `json:"error"`
	Status          int         `default:"200" json:"status"`
	Messages        string      `default:"Success" json:"message"`
	Data            interface{} `default:"[]" json:"data"`
	Search          string      `default:"" json:"search"`
	Next            bool        `default:"false" json:"next"`
	Back            bool        `default:"false" json:"back"`
	Limit           int         `default:"10" json:"limit"`
	Offset          int         `default:"0" json:"offset"`
	TotalPage       int         `default:"0" json:"total_page"`
	CurrentPage     int         `default:"1" json:"current_page"`
	Sort            string      `default:"ASC" json:"sort"`
	Order           string      `default:"id" json:"order"`
	Summary         interface{} `json:"summary,omitempty"`
	LastUpdated     null.Time   `json:"last_updated"`
}

func PopulatePaging(c echo.Context, custom string) (param ReqPaging) {
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
	order := c.QueryParam("sort")
	if strings.ToLower(order) == "asc" {
		order = "ASC"
	} else {
		order = "DESC"
	}
	sort := c.QueryParam("order")
	if sort == "" {
		sort = "created_at " + order
	} else {
		sort = sort + " " + order + ", created_at " + order
	}
	param = ReqPaging{
		Search: c.QueryParam("search"),
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Custom: customval,
		Page:   page}
	return
}

func PopulateResPaging(param *ReqPaging, data interface{}, totalResult int64, totalFiltered int64, lastUpdated null.Time) (output ResPaging) {
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

	output = ResPaging{
		Status:          200,
		Data:            data,
		Search:          param.Search,
		Order:           param.Order,
		Limit:           param.Limit,
		Offset:          param.Offset,
		Sort:            param.Sort,
		Next:            next,
		Back:            back,
		TotalData:       int(totalResult),
		RecordsFiltered: int(totalFiltered),
		CurrentPage:     currentPage,
		TotalPage:       totalPages,
		LastUpdated:     lastUpdated,
	}
	return
}
