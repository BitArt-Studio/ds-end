package request

type PageReq struct {
	PageNo   int      `json:"pageNo" valid:"pageNo" form:"pageNo"`
	PageSize int      `json:"pageSize" valid:"pageSize" form:"pageSize"`
	Orders   []string `json:"orders" valid:"orders" form:"orders"`
	Fields   []string `json:"fields" valid:"fields" form:"fields"`
}

type PageResp[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"` // nil 表示查询出错， [] 表示数据为空
}
