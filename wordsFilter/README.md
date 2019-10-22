# wordsFilter

单词过滤器（支持忽略指定符号）
```
  ignore := "!@#$%^&*()_+/*-="
	words := "毛泽东;毛主席;狗屎;臭狗屎;fuck;fuck you;"
```

test：
```
  fmt.Println(wordsFilter.Replace("fuc++k"), wordsFilter.Check("fu-+ck"))
	fmt.Println(wordsFilter.Replace("fcu++k"), wordsFilter.Check("fc-+uk"))
	fmt.Println(wordsFilter.Replace("fuc++k you"), wordsFilter.Check("fu-+ck you"))
	fmt.Println(wordsFilter.Replace("你是==狗屎++"), wordsFilter.Check("你是狗+-/屎"))
	fmt.Println(wordsFilter.Replace("你是==狗屎//臭狗屎"), wordsFilter.Check("你是--狗+-/屎臭狗屎"))
```
output：
```
  ****** true
  fcu++k false
  ********** true
  你是==**++ true
  你是==**//*** true
```
