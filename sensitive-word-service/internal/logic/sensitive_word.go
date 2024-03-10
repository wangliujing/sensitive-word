package logic

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/wangliujing/foundation-framework/err/biz"
	"github.com/wangliujing/sensitive-word/internal/core"
	"github.com/wangliujing/sensitive-word/internal/infra/repo"
	"github.com/wangliujing/sensitive-word/internal/pojo/dto"
	"github.com/wangliujing/sensitive-word/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync"
	"time"
)

type SensitiveWordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lock   sync.Mutex
}

func (s *SensitiveWordLogic) InitTrie() error {
	t := time.Now()
	if s.lock.TryLock() {
		defer s.lock.Unlock()
	} else {
		return biz.New("正在初始化，无须重复初始化")
	}

	logx.Info("初始化敏感词库...")
	trie := core.New()
	trie.SetSkip(core.NewSkipSpecialSymbols())
	trie.AddFilters(core.NewSingleEnglishWordFilter(), core.NewSpecialSymbolsFilter())
	sensitiveWordRepo := repo.NewSensitiveWordRepoImpl(s.svcCtx)
	page := 0
	pageSize := 100000
	for {
		page++
		words, err := sensitiveWordRepo.Load(page, pageSize)
		if err != nil {
			logx.Errorf("初始化敏感词库异常：%+v", err)
			return err
		}
		if words == nil || len(words) == 0 {
			logx.Info("初始化敏感词库结束..., 耗时：", time.Now().Sub(t))
			s.svcCtx.Trie = trie
			return nil
		}
		trie.AddKeywords(words...)
	}
}

func (s *SensitiveWordLogic) Detect(text *string) (*dto.DetectResult, error) {
	if text == nil || len(*text) == 0 {
		return &dto.DetectResult{
			Text: "",
		}, nil
	}
	emits := s.svcCtx.Trie.FindAll(*text)
	tokens := core.Tokenize(emits, *text)
	var builder strings.Builder
	keywordHashSet := hashset.New()
	for _, token := range tokens {
		if token.Emit != nil {
			builder.WriteString("<em>")
			builder.WriteString(token.Fragment)
			builder.WriteString("</em>")
			keywordHashSet.Add(token.Fragment)
		} else {
			builder.WriteString(token.Fragment)
		}
	}

	result := &dto.DetectResult{
		Text:    builder.String(),
		Keyword: keywordHashSet.Values(),
	}

	if keywordHashSet.Size() != 0 {
		result.ContainsSensitive = true
	}
	return result, nil
}

func NewSensitiveWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SensitiveWordLogic {
	return &SensitiveWordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
