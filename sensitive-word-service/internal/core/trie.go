package core

import (
	"container/list"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Emit struct {
	Begin, SkipNum, End int
	Keyword             string
}

func (e *Emit) Length() int {
	return e.End - e.Begin
}

func (e *Emit) Equals(o *Emit) bool {
	return e.Begin == o.Begin && e.End == o.End && e.Keyword == o.Keyword
}

func (e *Emit) Overlaps(o *Emit) bool {
	return e.Begin < o.End && e.End > o.Begin
}

func (e *Emit) Contains(o *Emit) bool {
	return e.Begin <= o.Begin && e.End >= o.End
}

func (e *Emit) String() string {
	return fmt.Sprintf("%d:%d:%d=%s", e.Begin, e.SkipNum, e.End, e.Keyword)
}

type Token struct {
	Fragment string
	Emit     *Emit
}

func (t *Token) IsMatch() bool {
	return t.Emit != nil
}

func (t *Token) String() string {
	if t.Emit == nil {
		return t.Fragment
	} else {
		return fmt.Sprintf("%s(%v)", t.Fragment, t.Emit)
	}
}

type State struct {
	//current        string
	success        map[rune]*State
	failure        *State
	keywordLengths []uint8
}

func (s *State) NextState(c rune) *State {
	return s.GetState(c)

}

func (s *State) GetState(c rune) *State {
	if s.success == nil {
		return nil
	}
	// 忽略大小写
	c = s.ignoreCase(c)
	state, exists := s.success[c]
	if exists {
		return state
	}
	return nil
}

func (s *State) AddState(str string) *State {
	state := s
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		state = state.addState(runes[i])
	}
	return state
}

func (s *State) addState(c rune) *State {
	if s.success == nil {
		s.success = make(map[rune]*State)
	}
	// 忽略大小写
	c = s.ignoreCase(c)
	state, exists := s.success[c]
	if exists {
		return state
	}
	ns := &State{ /*current: string(c)*/ }
	s.success[c] = ns
	return ns
}

func (s *State) ignoreCase(c rune) rune {
	if c >= 'A' && c <= 'Z' {
		c = unicode.ToLower(c)
	}
	return c
}

func (s *State) HasKeyword(keywordLength int) bool {
	for _, length := range s.keywordLengths {
		if length == uint8(keywordLength) {
			return true
		}
	}
	return false
}

func (s *State) AddKeyword(keyword string) {
	if len(keyword) == 0 {
		return
	}
	s.ensureKeywords()
	length := utf8.RuneCountInString(keyword)
	if !s.HasKeyword(length) {
		s.keywordLengths = append(s.keywordLengths, uint8(length))
	}
}

func (s *State) AddKeywords(keywordLengths []uint8) {
	if len(keywordLengths) == 0 {
		return
	}
	s.ensureKeywords()
	for _, keywordLength := range keywordLengths {
		if !s.HasKeyword(int(keywordLength)) {
			s.keywordLengths = append(s.keywordLengths, keywordLength)
		}
	}
}

func (s *State) ensureKeywords() {
	if s.keywordLengths == nil {
		s.keywordLengths = make([]uint8, 0, 2)
	}
}

type Trie struct {
	root    *State
	filters []Filter
	skip    Skip
}

type Filter interface {
	Do(runes []rune, emit *Emit) bool
}

type Skip interface {
	Do(t *Trie, s *State, runes []rune, index int) (bool, *State, int, int)
}

func New(keywords ...string) *Trie {
	t := Trie{root: &State{}}
	if len(keywords) > 0 {
		t.AddKeywords(keywords...)
	}
	return &t
}

func (t *Trie) AddFilters(filters ...Filter) *Trie {
	t.filters = append(t.filters, filters...)
	return t
}

func (t *Trie) SetSkip(skip Skip) *Trie {
	t.skip = skip
	return t
}

func (t *Trie) AddKeywords(keywords ...string) *Trie {
	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if len(keyword) > 0 {
			t.root.AddState(keyword).AddKeyword(keyword)
		}
	}
	states := list.New()
	// 构建失败指针
	t.root.failure = t.root
	for _, state := range t.root.success {
		state.failure = t.root
		states.PushBack(state)
	}
	// 层次遍历
	for states.Len() > 0 {
		state := states.Remove(states.Front()).(*State)
		if state.success == nil {
			continue
		}
		for c, next := range state.success {
			f := state.failure
			fn := f.NextState(c)
			if fn == nil {
				fn = t.root.NextState(c)
				if fn == nil {
					fn = t.root
				}
			}
			next.failure = fn
			if fn.keywordLengths != nil {
				next.AddKeywords(fn.keywordLengths)
			}
			states.PushBack(next)
		}
	}
	return t
}

func (t *Trie) FindAll(text string) []*Emit {
	emits := make([]*Emit, 0, 10)
	state := t.root
	runes := []rune(text)
	skipNum := 0
	for i := 0; i < len(runes); i++ {
		// 目标串字符
		r := runes[i]
		lastState := state
		state = state.NextState(r)
		if state == nil {
			// 匹配不成功尝试跳过逻辑
			if t.skip != nil {
				b, newState, index, subLength := t.skip.Do(t, lastState, runes, i)
				if b {
					skipNum += index + 1 - i - subLength
					if newState == nil {
						state = lastState
					} else {
						state = newState
					}
					i = index
					continue
				}
				//skipNum = 0
			}
			newState, b := t.nextState(lastState.failure, runes[i])
			state = newState
			if b {
				skipNum = 0
			}
			if state == nil {
				state = t.root
				continue
			}
		}
		if state.keywordLengths == nil {
			continue
		}

	LoopTag:
		for j := 0; j < len(state.keywordLengths); j++ {
			kwl := int(state.keywordLengths[j])
			begin := i + 1 - kwl - skipNum
			end := i + 1
			emit := &Emit{begin, skipNum, end, string(runes[begin:end])}
			if t.filters != nil {
				for _, filter := range t.filters {
					// 还要有一个过滤器不通过就忽略该敏感词
					if !filter.Do(runes, emit) {
						continue LoopTag
					}
				}

			}
			emits = append(emits, emit)
		}
	}
	emits = RemoveOverlaps(emits)
	emits = RemoveContains(emits)
	return emits
}

func IsEnglish(char rune) bool {
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
		return true
	}
	return false
}

func IsNumber(char rune) bool {
	if char >= '0' && char <= '9' {
		return true
	}
	return false
}

var SpecialSymbols = hashset.New('~', '`', '@', '#', '$', '%', '^', '&', '*', '-', '_', '+', '=', ':', '\\', '|', '/', ' ', '\r', '\n', '·', '.')

func IsSpecialSymbols(char rune) bool {
	return SpecialSymbols.Contains(char)
}

func (t *Trie) nextState(state *State, c rune) (*State, bool) {
	for {
		next := state.NextState(c)
		if state == t.root {
			// 从root重新匹配skip标记就需要重置
			// true 重置skipNum标记
			return next, true
		}
		if next == nil {
			state = state.failure
		} else {
			return next, false
		}
	}
}

func Tokenize(emits []*Emit, source string) []*Token {
	emits = RemoveContains(emits)
	el := len(emits)
	if el == 0 {
		return []*Token{{source, nil}}
	}
	index := 0
	runes := []rune(source)
	tokens := make([]*Token, 0, el*2+1)
	for i := 0; i < el; i++ {
		emit := emits[i]
		if index < emit.Begin {
			tokens = append(tokens, &Token{string(runes[index:emit.Begin]), nil})
		}
		tokens = append(tokens, &Token{string(runes[emit.Begin:emit.End]), emit})
		index = emit.End
	}
	last := emits[el-1]
	if last.End < utf8.RuneCountInString(source) {
		tokens = append(tokens, &Token{string(runes[last.End:]), nil})
	}
	return tokens
}

func Replace(emits []*Emit, source string, replacement string) string {
	emits = RemoveContains(emits)
	el := len(emits)
	if el == 0 {
		return source
	}
	index := 0
	runes := []rune(source)
	masks := []rune(replacement)
	ml := len(masks)
	for i := 0; i < el; i++ {
		emit := emits[i]
		if index < emit.Begin {
			index = emit.Begin
		}
		for j := emit.Begin; j < emit.End; j++ {
			runes[j] = masks[j%ml]
		}
		index = emit.End
	}
	return string(runes)
}

func RemoveOverlaps(emits []*Emit) []*Emit {
	return removeEmits(emits, func(a, b *Emit) bool {
		return a.Overlaps(b)
	})
}

func RemoveContains(emits []*Emit) []*Emit {
	return removeEmits(emits, func(a, b *Emit) bool {
		return a.Contains(b)
	})
}

func removeEmits(emits []*Emit, predicate func(a, b *Emit) bool) []*Emit {
	el := len(emits)
	if el < 1 {
		return nil
	} else if el == 1 {
		return []*Emit{emits[0]}
	}
	replica := make([]*Emit, el)
	copy(replica, emits)
	sortEmits(replica)
	emit := replica[0]
	sorted := make([]*Emit, 0, el)
	sorted = append(sorted, emit)
	for i := 1; i < el; i++ {
		next := replica[i]
		if !predicate(emit, next) {
			sorted = append(sorted, next)
			emit = next
		}
	}
	return sorted
}

func sortEmits(emits []*Emit) {
	sort.Slice(emits, func(i, j int) bool {
		a, b := emits[i], emits[j]
		if a.Begin != b.Begin {
			return a.Begin < b.Begin
		} else {
			return a.End > b.End
		}
	})
}
