package jsonrpc

import (
	"context"
	"github.com/wangliujing/sensitive-word/internal/logic"
	"github.com/wangliujing/sensitive-word/internal/pojo/dto"
	"github.com/wangliujing/sensitive-word/internal/svc"
)

type Text struct {
	Content string
}

type SensitiveWordRpcService struct {
	serviceContext *svc.ServiceContext
}

func NewSensitiveWordRpcService(serviceContext *svc.ServiceContext) *SensitiveWordRpcService {
	return &SensitiveWordRpcService{
		serviceContext: serviceContext,
	}
}

func (s *SensitiveWordRpcService) Detect(ctx context.Context, text *Text, reply *dto.DetectResult) error {
	wordLogic := logic.NewSensitiveWordLogic(ctx, s.serviceContext)
	result, err := wordLogic.Detect(&text.Content)
	*reply = *result
	return err
}
