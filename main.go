package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	ydfanyi "github.com/hnmaonanbei/go-youdao-fanyi"
)

func main() {

	txt := read()

	tt := strings.Fields(txt)

	len := len(tt)

	arr := generate(len)

	b := make([]byte, 1)
	opts := ydfanyi.NewOptions("", "", "")
	opts.To = ydfanyi.ZH
	opts.From = ydfanyi.EN

	for i := range arr {
		fmt.Print(tt[arr[i]])
		res, _ := ydfanyi.Do(tt[arr[i]], opts)

		os.Stdin.Read(b)
		fmt.Println(res.String())
		os.Stdin.Read(b)
	}
	fmt.Println("_(:з」∠)_")
	os.Stdin.Read(b)
}

func read() string {
	fmt.Println("请输入你要背的词汇表的名称")
	var book string
	fmt.Scanln(&book)
	f, err := ioutil.ReadFile(book)
	if err != nil {
		fmt.Println("请看清文件名，加上后缀")
	}

	fmt.Println("ヾ(•ω•`)o")
	return string(f)
}
func generate(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
