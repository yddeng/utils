package wordsFilter_test

import (
	"fmt"
	"github.com/tagDong/dutil/wordsFilter"
)

func main() {
	fmt.Println(wordsFilter.Replace("fuc++k"), wordsFilter.Check("fu-+ck"))
	fmt.Println(wordsFilter.Replace("fcu++k"), wordsFilter.Check("fc-+uk"))
	fmt.Println(wordsFilter.Replace("fuc++k you"), wordsFilter.Check("fu-+ck you"))
	fmt.Println(wordsFilter.Replace("你是==狗屎++"), wordsFilter.Check("你是狗+-/屎"))
	fmt.Println(wordsFilter.Replace("你是==狗屎//臭狗屎"), wordsFilter.Check("你是--狗+-/屎臭狗屎"))
}
