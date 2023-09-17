package tuuid

import "github.com/google/uuid"

func GenerateUUID() string {
	uid, err := uuid.NewUUID()
	if err != nil {
		// 生成出错
	}
	return uid.String()
}
