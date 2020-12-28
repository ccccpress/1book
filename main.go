package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	ydfanyi "github.com/hnmaonanbei/go-youdao-fanyi"

	"github.com/mattn/go-runewidth"
)

func main() {

	if len(os.Args) == 1 {
		os.Exit(1)
	}
	//读取词汇表
	tt := strings.Fields(read(os.Args[1]))
	//打乱设置顺序
	len := len(tt)
	arr := generate(len)
	//显示
	encoding.Register()
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	displayHelloWorld(s, "开始背单词~")
	time.Sleep(time.Second)
	i := 0
	j := 0
	zh := ""

	//开始循环单词表
	for {
		if i == len {
			displayHelloWorld(s, "放映结束，按任意键退出")
		}
		if j == 0 && i != len {
			displayHelloWorld(s, tt[arr[i]])
			zh = trans(tt[arr[i]])
		}

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || i == len {
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEnter {
				if j == 0 {
					displayHelloWorld(s, zh)
					j++
				} else {
					i++
					j--
				}
			}
		}
	}
}

func read(str string) string {
	f, err := ioutil.ReadFile(str)
	if err != nil {
		fmt.Println("请看清文件名，加上后缀")
	}
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
func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 { //这里应该是为了避免零宽字符
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}
func displayHelloWorld(s tcell.Screen, str string) {
	w, h := s.Size()
	s.Clear()
	emitStr(s, w/2-len(str)/2, h/2-1, tcell.StyleDefault, str)
	s.Show()
}
func trans(str string) string {
	opts := ydfanyi.NewOptions("", "", "")
	res, _ := ydfanyi.Do(str, opts)
	zh := strings.Join(res.SmartResult.Entries, "")
	return zh
}
