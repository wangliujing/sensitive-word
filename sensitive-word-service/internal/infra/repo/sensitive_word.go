package repo

import (
	"context"
	"github.com/wangliujing/sensitive-word/internal/svc"
)

type LoadParam struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
}

type LoadResult struct {
	List []Word `json:"list"`
}

type Word struct {
	Name string `json:"name"`
}

type SensitiveWordRepo interface {
	Load(page int, pageSize int) ([]string, error)
}

type SensitiveWordRepoImpl struct {
	svcCtx *svc.ServiceContext
}

func (s *SensitiveWordRepoImpl) Load(page int, pageSize int) ([]string, error) {
	result := &LoadResult{}
	err := s.svcCtx.JsonRpcClient.Call(context.Background(), "AmzSensitiveService",
		"/amz_sensitive/sensitiveLangListPage", LoadParam{CurrentPage: page, PageSize: pageSize}, result)
	if err != nil {
		return nil, err
	}
	var strList []string
	if result.List != nil && len(result.List) != 0 {
		for _, word := range result.List {
			strList = append(strList, word.Name)
		}
	}
	return strList, nil
}

func NewSensitiveWordRepoImpl(svcCtx *svc.ServiceContext) SensitiveWordRepo {
	return &SensitiveWordRepoImpl{
		svcCtx: svcCtx,
	}
}
