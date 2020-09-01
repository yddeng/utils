package wordsFilter

import (
	"fmt"
)

func main() {
	fmt.Println(Replace("fuc++k"), Check("fu-+ck"))
	fmt.Println(Replace("fcu++k"), Check("fc-+uk"))
	fmt.Println(Replace("fuc++k you"), Check("fu-+ck you"))
	fmt.Println(Replace("你是==狗屎++"), Check("你是狗+-/屎"))
	fmt.Println(Replace("你是==狗屎//臭狗屎"), Check("你是--狗+-/屎臭狗屎"))
}
