package core

type SkipSpecialSymbols struct {
}

func NewSkipSpecialSymbols() Skip {
	return &SkipSpecialSymbols{}
}

func (s *SkipSpecialSymbols) Do(t *Trie, state *State, runes []rune, index int) (bool, *State, int, int) {
	if t.root == state {
		return false, state, index, 0
	}
	// 查看前面是否是英文
	for i := index - 1; i >= 0; i-- {
		character := runes[i]
		if IsSpecialSymbols(character) {
			continue
		}

		if IsEnglish(character) {
			break
		}
		return false, state, index, 0
	}

	newIndex := index
	// 包含特殊符号且前后除特殊符号外是英文
	for IsSpecialSymbols(runes[newIndex]) {
		// 第一次进来的时候判断敏感词是否有空格，没有空格则不匹配
		if newIndex == index {
			// 判断前面字符是否是空格，如果不是则判断敏感词是否有空格
			if runes[newIndex-1] != ' ' {
				newState := state.NextState(' ')
				if newState == nil {
					return false, state, index, 0
				}
			}
		}
		newIndex++
		if newIndex == len(runes) {
			return false, state, index, 0
		}
	}
	// 不包含特殊符号直接返回
	if newIndex == index {
		return false, state, index, 0
	}
	// 查看后面是否是英文
	if IsEnglish(runes[newIndex]) {
		/*
			*	如果这个时候state还包含空格则证明原字符串字母后面没有空格与之匹配
				text := "d&**%b 中 Nick"
				trie := New("d b")
				如果不处理如上代码就匹配不上
		*/

		newState := state.NextState(' ')
		if newState != nil {
			return true, newState, newIndex - 1, 1
		}

		return true, state, newIndex - 1, 0
	}
	return false, state, index, 0
}

type SingleEnglishWordFilter struct {
}

func NewSingleEnglishWordFilter() Filter {
	return &SingleEnglishWordFilter{}
}

func (f *SingleEnglishWordFilter) Do(runes []rune, emit *Emit) bool {
	words := []rune(emit.Keyword)
	for _, word := range words {
		// 如果不是英文则直接返回true
		if !IsEnglish(word) {
			return true
		}
	}
	beginWord := ' '
	endWord := ' '
	// 查看前后是否都是数字或字母，如果是则返回false
	if emit.Begin > 0 {
		beginWord = runes[emit.Begin-1]
	}
	if emit.End < len(runes) {
		endWord = runes[emit.End]
	}
	if IsEnglish(beginWord) || IsEnglish(endWord) {
		return false
	}
	return true
}

type SpecialSymbolsFilter struct {
}

func NewSpecialSymbolsFilter() Filter {
	return &SpecialSymbolsFilter{}
}

func (f *SpecialSymbolsFilter) Do(runes []rune, emit *Emit) bool {
	if emit.SkipNum == 0 {
		return true
	}
	runes = []rune(emit.Keyword)
	/*	// 如果开头是特殊字符直接切
		if IsSpecialSymbols(runes[0]) {
			emit.Begin = emit.Begin + emit.SkipNum
			emit.Keyword = string(runes[emit.SkipNum:])
			emit.SkipNum = 0
		}*/

	newRunes := runes[emit.SkipNum:]
	// 如果切掉后第一个不是特殊字符直接切
	if !IsSpecialSymbols(newRunes[0]) {
		emit.Begin = emit.Begin + emit.SkipNum
		emit.Keyword = string(newRunes)
		emit.SkipNum = 0
	}

	return true
}
