package input

import (
	"mattwach/rpngo/common/rpn"
	"sort"
	"strings"
)

func (gl *getLine) tabComplete(r *rpn.RPN, idx int) int {
	if (idx <= 0) || (idx > len(gl.line)) {
		return idx
	}
	if (idx < len(gl.line)) && gl.line[idx] != ' ' {
		return idx
	}

	startIdx := gl.findStartOfWord(idx)
	if startIdx == idx {
		return idx
	}

	word := string(gl.line[startIdx:idx])
	newWord := gl.findNewWord(r, word)

	if len(newWord) == 0 {
		return idx
	}

	startLine := string(gl.line[:startIdx])
	endLine := string(gl.line[idx:])
	gl.line = gl.line[:0]
	for _, c := range startLine {
		gl.line = append(gl.line, byte(c))
	}
	for _, c := range newWord {
		gl.line = append(gl.line, byte(c))
	}
	for _, c := range endLine {
		gl.line = append(gl.line, byte(c))
	}

	// update the line
	gl.txtb.Shift(startIdx - idx)
	gl.txtb.PrintBytes(gl.line[startIdx:], true)
	numSpaces := len(word) - len(newWord)
	if numSpaces > 0 {
		for i := 0; i < numSpaces; i++ {
			gl.txtb.Write(' ', true)
		}
		gl.txtb.Shift(-numSpaces)
	}
	gl.txtb.Shift(-len(endLine))

	idx = idx + len(newWord) - len(word)
	return idx
}

func (gl *getLine) findStartOfWord(idx int) int {
	startIdx := idx - 1
	for {
		if startIdx <= 0 {
			return 0
		}
		lastChar := gl.line[startIdx]
		switch lastChar {
		case '@', '$', '{', '\'', '"':
			return startIdx
		case ' ':
			return startIdx + 1
		}
		startIdx--
	}
}

func (gl *getLine) findNewWord(r *rpn.RPN, word string) string {
	var wordList []string
	var varPrefix string
	switch word[0] {
	case '$':
		varPrefix = "$"
		word = word[1:]
		wordList = gl.allVariableNames(r)
	case '@':
		varPrefix = "@"
		word = word[1:]
		wordList = gl.allStringVariables(r)
	case '\'', '{', '"':
		varPrefix = word[:1]
		word = word[1:]
		wordList = gl.allStringVariables(r)
		gl.getFileList()
		wordList = gl.fileList
	default:
		wordList = r.AllFunctionNames()
	}

	// Look for an exact match of the word
	var newWord string
	for wordIdx := 0; wordIdx < len(wordList); wordIdx++ {
		if wordList[wordIdx] == word {
			newWord = wordList[(wordIdx+1)%len(wordList)]
			break
		}
	}

	if len(newWord) == 0 {
		// look for a partial match
		for wordIdx := 0; wordIdx < len(wordList); wordIdx++ {
			if strings.HasPrefix(wordList[wordIdx], word) {
				newWord = wordList[wordIdx]
				break
			}
		}
	}

	return varPrefix + newWord
}

func (gl *getLine) allVariableNames(r *rpn.RPN) []string {
	gl.names = r.AppendAllVariableNames(gl.names[:0])
	sort.Strings(gl.names)
	return gl.names
}

func (gl *getLine) allStringVariables(r *rpn.RPN) []string {
	var wordList []string
	fn := func(name string, vals []rpn.Frame) bool {
		if vals[len(vals)-1].IsString() {
			wordList = append(wordList, name)
		}
		return true
	}
	r.IterateAllVariables(fn)
	sort.Strings(wordList)
	return wordList
}

func (gl *getLine) getFileList() {
	if len(gl.fileList) > 0 {
		return
	}
	gl.fileList, _ = gl.fs.ListFiles(".", gl.fileList)
}
