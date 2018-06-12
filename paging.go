package paging4go

import (
	"github.com/smartwalle/time4go"
	"math"
	"strconv"
	"strings"
)

const (
	K_DATE_TIME_FORMAT = "2006-01-02T15:04:05-07:00"
)

// --------------------------------------------------------------------------------
type Pagination interface {
	GetKeywords() string
	GetLimit() uint64
	GetPage() uint64
	GetBeginTime() *time4go.Time
	GetEndTime() *time4go.Time
}

// --------------------------------------------------------------------------------
// ListForm 用于转换 http request 的分页参数
type ListForm struct {
	Keywords  string        `form:"keywords"`
	Limit     uint64        `form:"limit"`
	Page      uint64        `form:"pageInfo"`
	BeginTime *time4go.Time `form:"begin_time"`
	EndTime   *time4go.Time `form:"end_time"`
}

func (this *ListForm) CleanedPage(p string) uint64 {
	var page, _ = strconv.ParseUint(p, 10, 64)
	page = page - 1
	if page < 0 {
		return 0
	}
	return page
}

func (this *ListForm) CleanedBeginTime(p string) *time4go.Time {
	var t, err = time4go.Parse(K_DATE_TIME_FORMAT, p)
	if err == nil {
		t = t.Local()
		return t
	}
	return nil
}

func (this *ListForm) CleanedEndTime(p string) *time4go.Time {
	var t, err = time4go.Parse(K_DATE_TIME_FORMAT, p)
	if err == nil {
		t = t.Local()
		return t
	}
	return nil
}

func (this *ListForm) CleanedKeywords(p string) string {
	return strings.TrimSpace(p)
}

func (this *ListForm) DefaultLimit() uint64 {
	return 20
}

func (this *ListForm) DefaultPage() uint64 {
	return 0
}

func (this *ListForm) GetKeywords() string {
	return strings.TrimSpace(this.Keywords)
}

func (this *ListForm) GetLimit() uint64 {
	return this.Limit
}

func (this *ListForm) GetPage() uint64 {
	return this.Page
}

func (this *ListForm) GetBeginTime() *time4go.Time {
	return this.BeginTime
}

func (this *ListForm) GetEndTime() *time4go.Time {
	return this.EndTime
}

// --------------------------------------------------------------------------------
// ListData 用于返回给客户端
type ListData struct {
	Total uint64      `json:"total"      sql:"total"` // 共有多少条数据
	Page  uint64      `json:"page"       sql:"-"`     // 当前页码
	Limit uint64      `json:"limit"      sql:"-"`     // 每页最大数据量
	Data  interface{} `json:"data"       sql:"-"`     // 当前页的数据
}

func (this *ListData) GetPageInfo() *PageInfo {
	return pageInfo(this.Total, this.Page, this.Limit)
}

type PageInfo struct {
	Total    int   `json:"total"          sql:"-"`
	Page     int   `json:"page"           sql:"-"`
	Limit    int   `json:"limit"          sql:"-"`
	PrevPage int   `json:"prev_page"      sql:"-"`
	NextPage int   `json:"next_page"      sql:"-"`
	PageList []int `json:"page_list"      sql:"-"`
}

func pageInfo(total, currentPage, pageLimit uint64) (page *PageInfo) {
	var maxPage = 5
	var midPage = uint64(math.Ceil(float64(maxPage) / 2.0))
	var totalPage = int(math.Ceil(float64(total) / float64(pageLimit)))

	page = &PageInfo{}
	page.Total = totalPage
	page.Page = int(currentPage)
	page.Limit = int(pageLimit)

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
		for i := 0; i < totalPage; i++ {
			page.PageList = append(page.PageList, i+1)
		}
	} else {
		var beginPage = int(math.Max(float64(int(currentPage-midPage)), 0))
		var endPage = int(math.Min(float64(beginPage)+float64(maxPage), float64(totalPage)))
		if endPage-beginPage < maxPage {
			beginPage = beginPage - (maxPage - (endPage - beginPage))
		}
		for i := beginPage; i < endPage; i++ {
			page.PageList = append(page.PageList, i+1)
		}
	}
	return page
}
