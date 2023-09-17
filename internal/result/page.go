package result

// Page 分页查询
type Page struct {
	PageNum  int    `json:"pageNum" form:"pageNum" comment:"页码数量"`
	PageSize int    `json:"pageSize" form:"pageSize" comment:"分页大小"`
	Keyword  string `json:"keyword" form:"keyword" comment:"关键字 => 模糊搜索"`
	Desc     bool   `json:"desc" form:"desc" comment:"是否反向搜索"`
}
