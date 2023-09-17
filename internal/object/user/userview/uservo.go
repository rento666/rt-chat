package userview

// 接口接收到的用户信息，如登录、注册时的`账密验`组合

// Ident 账号类型、账号、密码
type Ident struct {
	IdentityType string `bind:"required" json:"identityType" form:"identityType" comment:"用户标识类型"`
	Identifier   string ` bind:"required" json:"identifier" form:"identifier" comment:"用户登录标识"`
	Credential   string ` bind:"required" json:"credential" form:"credential" comment:"用户密钥"`
}

// RegisterVo 注册信息
type RegisterVo struct {
	Ident
	Code string `bind:"required" json:"code" form:"code" comment:"注册验证码"`
}

// LoginVo 登录信息
type LoginVo struct {
	Info string `bind:"required" json:"info" form:"info" comment:"账号类型、账号、密码"`
}

// IdentifierCode 验证码接收人信息
type IdentifierCode struct {
	IdentityType string `bind:"required" json:"identityType" form:"identityType" comment:"注册标识"`
	Identifier   string `bind:"required" json:"identifier" form:"identifier" comment:"账号"`
}

// UserList 用户列表信息
type UserList struct {
	Uuid     string `bind:"required" json:"uuid" form:"uuid" comment:"用户唯一标识"`
	Nickname string `bind:"required" json:"nickname" form:"nickname" comment:"用户昵称"`
	Avatar   string `bind:"required" json:"avatar" form:"avatar" comment:"用户头像"`
	Status   bool   `bind:"required" json:"status" form:"status" comment:"用户状态"`
	Phone    string `bind:"required" json:"phone" form:"phone" comment:"用户手机号"`
	Email    string `bind:"required" json:"email" form:"email" comment:"用户邮箱"`
	Address  string `bind:"required" json:"address" form:"address" comment:"用户地址"`
	// ------- 因为用户有多种登录方式，所以这里存储账密的切片
	AuthList []Ident
	// 与deletedAt是否为空有关，为空是启用，不为空是删除——【软删除技术】
	IsDeleted bool `bind:"required" json:"isDeleted" form:"isDeleted" comment:"是否删除"`
}
