package paging4go

import (
	"math"
	"strconv"
	"strings"
)

// --------------------------------------------------------------------------------
type Pagination interface {
	GetKeywords() string
	GetLimit() int64
	GetPage() int64
	GetBeginTime() int64
	GetEndTime() int64
	GetOrderBy() []string
}

// --------------------------------------------------------------------------------
// ListForm 用于转换 http request 的分页参数
type ListForm struct {
	Keywords  string   `form:"keywords"        json:"keywords"`
	Limit     int64    `form:"limit"           json:"limit"`
	Page      int64    `form:"page"            json:"page"`
	OrderBy   []string `form:"order_by"        json:"order_by"`
	BeginTime int64    `form:"begin_time"      json:"begin_time"`
	EndTime   int64    `form:"end_time"        json:"end_time"`
}

func (this *ListForm) CleanedPage(p string) int64 {
	var page, _ = strconv.ParseInt(p, 10, 64)
	page = page - 1
	if page < 0 {
		return 0
	}
	return page
}

func (this *ListForm) CleanedKeywords(p string) string {
	return strings.TrimSpace(p)
}

func (this *ListForm) DefaultLimit() int64 {
	return 20
}

func (this *ListForm) DefaultPage() int64 {
	return 0
}

func (this *ListForm) GetKeywords() string {
	return strings.TrimSpace(this.Keywords)
}

func (this *ListForm) GetLimit() int64 {
	return this.Limit
}

func (this *ListForm) GetPage() int64 {
	return this.Page
}

func (this *ListForm) GetBeginTime() int64 {
	return this.BeginTime
}

func (this *ListForm) GetEndTime() int64 {
	return this.EndTime
}

func (this *ListForm) GetOrderBy() []string {
	return this.OrderBy
}

// --------------------------------------------------------------------------------
// ListData 用于返回给客户端
type ListData struct {
	Total int64       `json:"total"      sql:"total"` // 共有多少条数据
	Page  int64       `json:"page"       sql:"-"`     // 当前页码
	Limit int64       `json:"limit"      sql:"-"`     // 每页最大数据量
	Data  interface{} `json:"data"       sql:"-"`     // 当前页的数据
}

func (this *ListData) GetPageInfo() *PageInfo {
	return pageInfo(this.Total, this.Page, this.Limit)
}

type PageInfo struct {
	Total    int64   `json:"total"          sql:"-"`
	Page     int64   `json:"page"           sql:"-"`
	Limit    int64   `json:"limit"          sql:"-"`
	PrevPage int64   `json:"prev_page"      sql:"-"`
	NextPage int64   `json:"next_page"      sql:"-"`
	PageList []int64 `json:"page_list"      sql:"-"`
}

func pageInfo(total, currentPage, pageLimit int64) (page *PageInfo) {
	var maxPage int64 = 5
	var midPage = int64(math.Ceil(float64(maxPage) / 2.0))
	var totalPage = int64(math.Ceil(float64(total) / float64(pageLimit)))

	page = &PageInfo{}
	page.Total = totalPage
	page.Page = currentPage
	page.Limit = pageLimit

	if page.Page == 1 || page.Total <= 0 {
		page.PrevPage = -1
	} else {
		page.PrevPage = page.Page - 1
	}
	if page.Page == page.Total || page.Total <= 0 {
		page.NextPage = -1
	} else {
		page.NextPage = page.Page + 1
	}

	if totalPage <= maxPage {
		var i int64 = 0
		for i = 0; i < totalPage; i++ {
			page.PageList = append(page.PageList, i+1)
		}
	} else {
		var beginPage = int64(math.Max(float64(int(currentPage-midPage)), 0))
		var endPage = int64(math.Min(float64(beginPage)+float64(maxPage), float64(totalPage)))
		if endPage-beginPage < maxPage {
			beginPage = beginPage - (maxPage - (endPage - beginPage))
		}
		for i := beginPage; i < endPage; i++ {
			page.PageList = append(page.PageList, i+1)
		}
	}
	return page
}
