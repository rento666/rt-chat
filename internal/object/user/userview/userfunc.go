package userview

import (
	"errors"
	"github.com/DanPlayer/randomname"
	"log"
	"rt-chat/internal/global"
	"rt-chat/internal/object/user"
	"rt-chat/internal/result"
	"rt-chat/pkg/tools/tcaptcha"
	"rt-chat/pkg/tools/temail"
	"rt-chat/pkg/tools/tencrypt"
	"rt-chat/pkg/tools/tjwt"
	"rt-chat/pkg/tools/tredis"
	"rt-chat/pkg/tools/ttrans"
	"rt-chat/pkg/tools/tuuid"
	"strconv"
	"time"
)

// JudgeExistFunc 根据类型、账号来判断用户是否存在
func JudgeExistFunc(vo user.Auth) (bool, user.Auth) {
	ua := user.Auth{}
	if vo.BasicId != 0 {
		ua.BasicId = vo.BasicId
	}
	if vo.IdentityType != "" {
		ua.IdentityType = vo.IdentityType
	}
	if vo.Identifier != "" {
		ua.Identifier = vo.Identifier
	}
	if vo.Credential != "" {
		ua.Credential = vo.Credential
	}
	return judgeExist(ua)
}

// IdentifierCodeFunc 向该账号发送验证码，账号的类型有手机号、邮箱......
func IdentifierCodeFunc(ic IdentifierCode) bool {
	return identifierCode(ic.IdentityType, ic.Identifier)
}

// RegisterFunc 注册用户函数,成功返回token，失败返回error
func RegisterFunc(form RegisterVo) (string, error) {
	// 解密数据并返回
	username, form, err := decryptRegisterData(form)
	if err != nil {
		return "", err
	}
	// 判断是否存在
	b, _ := JudgeExistFunc(user.Auth{IdentityType: form.IdentityType, Identifier: form.Identifier})
	if b {
		return "", errors.New("用户已存在，不能重复注册！")
	}
	// 判断验证码是否正确
	if !codeIsValid(form) {
		return "", errors.New("验证码错误，请重新输入！")
	}
	// 添加到数据库
	uid, err := addUserDB(username, form)
	if err != nil {
		return "", err
	}
	// 删除redis的code
	// 将redis中的验证码删除，防止多次使用
	_, err = tredis.DelByKey(global.REDIS, form.Identifier+"-register")
	if err != nil {
		return "", err
	}
	token, err := tjwt.GenerateToken(strconv.Itoa(int(uid)))
	if err != nil {
		return "", err
	}
	// 返回bool,err
	return token, nil
}

// LoginFunc 登录用户函数,成功返回token，失败返回error
func LoginFunc(form LoginVo) (string, error) {
	// 解密数据并返回
	rv, err := decryptLoginData(form)
	if err != nil {
		return "", err
	}
	// 判断是否存在
	b, u := JudgeExistFunc(rv)
	if !b {
		return "", errors.New("用户不存在，登录失败！")
	}
	token, err := tjwt.GenerateToken(strconv.Itoa(int(u.BasicId)))
	if err != nil {
		return "", err
	}
	// 返回bool,err
	return token, nil
}

// GetUserListFunc 获取用户列表函数，成功返回'用户切片'，失败返回error
func GetUserListFunc(p result.Page) ([]UserList, error) {
	return getList(p)
}

// ------------func分界线----------

// judgeExist 根据类型、账号来判断用户是否存在
func judgeExist(ua user.Auth) (bool, user.Auth) {
	u := user.Auth{}
	tx := global.DB.Where(&ua).Find(&u)
	if tx.RowsAffected == 0 {
		// 不存在
		return false, user.Auth{}
	} else {
		// 存在
		return true, u
	}
}

// identifierCode 发送验证码
func identifierCode(identityType, identifier string) bool {
	switch identityType {
	case user.EMAIL:
		co := ""
		var expire = time.Minute * 30
		// 如果redis中有值，那么直接取出就行，不用再次生成验证码了
		if s := tredis.GetString(global.REDIS, identifier+"-register"); s != "" && s != "key not exist" && s != "err" {
			co = s
			expire, _ = tredis.GetExpiration(global.REDIS, identifier+"-register")
		} else {
			co = tcaptcha.GenerateVerificationCode()
			tredis.SetString(global.REDIS, identifier+"-register", co, time.Minute*30)
		}
		go func() {
			err := temail.SendVerificationCode(identifier, co, expire)
			if err != nil {
				log.Println(err, "-验证码发送失败")
			}
		}()
	default:
		return false
	}
	return true
}

// decryptRegisterData 解密注册数据并返回 username、registerVO、error
func decryptRegisterData(form RegisterVo) (string, RegisterVo, error) {
	/*
		{
			type: type
			identifier: RSA加密(username(aes加密) | phone/email(aes加密))
			credential: password
			code: register_code
		}
	*/
	// rsa解密，得到 username(aes加密后) | identifier(aes加密后)
	plaintext, err := tencrypt.RsaDecrypt(ttrans.String2Bytes(form.Identifier))
	if err != nil {
		return "", RegisterVo{}, err
	}
	res := ttrans.Bytes2String(plaintext)
	// 分隔符 |
	sep := "|"
	// a => [username(aes加密后),identifier(aes加密后)]
	arr := ttrans.SplitInString(res, sep)
	if len(arr) != 2 {
		return "", RegisterVo{}, errors.New("账号填写有误！")
	}
	// aes解密，得到 username []byte类型
	u1, err := tencrypt.AesDecrypt(ttrans.String2Bytes(arr[0]))
	if err != nil {
		return "", RegisterVo{}, errors.New("用户名解密失败！")
	}
	username := ttrans.Bytes2String(u1)
	// aes解密，得到identifier []byte类型
	i1, err := tencrypt.AesDecrypt(ttrans.String2Bytes(arr[1]))
	if err != nil {
		return "", RegisterVo{}, errors.New("账号解密失败！")
	}
	identifier := ttrans.Bytes2String(i1)
	form.Identifier = identifier

	return username, form, nil
}

// codeIsValid 判断验证码是否正确
func codeIsValid(form RegisterVo) bool {
	if s := tredis.GetString(global.REDIS, form.Identifier+"-register"); form.Code != s {
		return false
	}
	return true
}

// addUserDB 添加用户到数据库，成功返回uuid
func addUserDB(username string, form RegisterVo) (uint, error) {
	// md5加密
	credential := tencrypt.MD5Encrypt(form.Credential)
	// 登录类型：账密、QQ、WX、钉钉、github、微博等
	UUID := tuuid.GenerateUUID()
	// 昵称取随机的,用了个github的库(github.com/DanPlayer/。。。)
	nickname := randomname.GenerateName()
	// 向数据库添加该字段 添加两个授权、一个基础
	base := &user.Basic{Uuid: UUID, Nickname: nickname, Status: true}
	switch form.IdentityType {
	case user.EMAIL:
		base.Email = form.Identifier
	case user.PHONE:
		base.Phone = form.Identifier
	default:
		return 0, errors.New("未知的账号类型！")
	}
	// 添加base表
	global.DB.Create(base)

	auth := &user.Auth{BasicId: base.ID, IdentityType: form.IdentityType,
		Identifier: form.Identifier, Credential: credential}
	// 添加auth表
	global.DB.Create(auth)
	auth.IdentityType = user.NAME
	auth.Identifier = username
	// 添加auth表
	global.DB.Create(auth)

	return base.ID, nil
}

// decryptLoginData 解密登录数据并返回 user_auth error
func decryptLoginData(form LoginVo) (user.Auth, error) {
	// res 为rsa解密后的[]byte转的string
	decrypt, err := tencrypt.RsaDecrypt(ttrans.String2Bytes(form.Info))
	if err != nil {
		return user.Auth{}, nil
	}
	res := ttrans.Bytes2String(decrypt)
	// sep 为分隔符， 例如： |
	sep := "|"
	// fs 为分割后的res， 例如： identityType|identifier|credential    ->   1|1371546845|888miMa123456
	fs := ttrans.SplitInString(res, sep)
	if len(fs) != 3 {
		return user.Auth{}, errors.New("数据有误，解密失败！")
	}
	// type
	identityType := fs[0]
	// aes解密
	temp, _ := tencrypt.AesDecrypt(ttrans.String2Bytes(fs[1]))
	// identifier
	identifier := ttrans.Bytes2String(temp)
	// credential
	credential := tencrypt.MD5Encrypt(fs[2])
	return user.Auth{IdentityType: identityType, Identifier: identifier, Credential: credential}, nil
}

// getList 获取用户列表，返回list，error
func getList(page result.Page) ([]UserList, error) {
	var users []user.Basic

	// 构建GORM查询
	query := global.DB.Model(&user.Basic{})

	// 添加模糊搜索条件
	if page.Keyword != "" {
		// 构建LIKE语句
		keyword := "%" + page.Keyword + "%"
		query = query.Where("nickname LIKE ? OR phone LIKE ? OR email LIKE ? OR address LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	// 添加反向搜索条件
	if page.Desc {
		query = query.Order("created_at DESC")
	} else {
		query = query.Order("created_at ASC")
	}

	// 添加分页条件
	offset := (page.PageNum - 1) * page.PageSize
	query = query.Offset(offset).Limit(page.PageSize)

	// Unscoped() 方法是包括上已经软删除的用户
	// 执行查询并预加载 Auth 信息 , 这个AuthList是basic新添加的，目的是与auth表关联
	if err := query.Unscoped().Preload("AuthList").Find(&users).Error; err != nil {
		return nil, err
	}

	// 将查询结果转换为 UserList 结构
	userList := make([]UserList, len(users))
	for i, user0 := range users {
		identList := make([]Ident, len(user0.AuthList))
		for j, auth := range user0.AuthList {
			identList[j] = Ident{
				IdentityType: auth.IdentityType,
				Identifier:   auth.Identifier,
			}
		}

		userList[i] = UserList{
			Uuid:      user0.Uuid,
			Nickname:  user0.Nickname,
			Avatar:    user0.Avatar,
			Status:    user0.Status,
			Phone:     user0.Phone,
			Email:     user0.Email,
			Address:   user0.Address,
			AuthList:  identList,
			IsDeleted: user0.DeletedAt.Valid,
		}

	}

	return userList, nil
}
