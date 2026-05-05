package tabcomplete

import (
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/rpn"
	"sort"
	"strings"
)

type TabComplete struct {
	fs    fileops.FileOpsDriver
	names []string
	// cached list of files for tab complete
	fileList []string
}

func (tc *TabComplete) Init(fs fileops.FileOpsDriver) {
	tc.fs = fs
	tc.fileList = make([]string, 0, 16)
	tc.names = make([]string, 0, 16)
}

func (tc *TabComplete) Clear() {
	tc.fileList = tc.fileList[:0]
}

func (tc *TabComplete) FindNewWord(r *rpn.RPN, word string) string {
	var wordList []string
	var varPrefix string
	switch word[0] {
	case '$':
		varPrefix = "$"
		word = word[1:]
		wordList = tc.allVariableNames(r)
	case '@':
		varPrefix = "@"
		word = word[1:]
		wordList = tc.allStringVariables(r)
	case '\'', '{', '"':
		varPrefix = word[:1]
		word = word[1:]
		wordList = tc.allStringVariables(r)
		tc.getFileList()
		wordList = tc.fileList
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

func (tc *TabComplete) allVariableNames(r *rpn.RPN) []string {
	tc.names = r.AppendAllVariableNames(tc.names[:0])
	sort.Strings(tc.names)
	return tc.names
}

func (tc *TabComplete) allStringVariables(r *rpn.RPN) []string {
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

func (tc *TabComplete) getFileList() {
	if len(tc.fileList) > 0 {
		return
	}
	tc.fileList, _ = tc.fs.ListFiles(".", tc.fileList)
}
