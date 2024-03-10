package repo

import (
	"github.com/wangliujing/sensitive-word/internal/svc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SensitiveWordRepoImplTest struct {
	svcCtx *svc.ServiceContext
}

type KeyWord struct {
	Name string
}

func (s *SensitiveWordRepoImplTest) Load(page int, pageSize int) ([]string, error) {
	if page == 2 {
		return nil, nil
	}
	var arr []string
	// 连接数据库
	dsn := "jmy:Qwer1234@tcp(rm-bp1k0fj4ej13saps7.mysql.rds.aliyuncs.com:3306)/amz_sensitive_center?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var word = make([]KeyWord, 100000)
	db.Table("amz_sensitive").Find(&word)
	for _, keyWord := range word {
		arr = append(arr, keyWord.Name)
	}

	word = make([]KeyWord, 100000)
	db.Table("amz_sensitive_language").Find(&word)
	for _, keyWord := range word {
		arr = append(arr, keyWord.Name)
	}
	return arr, nil
}

func NewSensitiveWordRepoImplTest(svcCtx *svc.ServiceContext) SensitiveWordRepo {
	return &SensitiveWordRepoImplTest{}
}
