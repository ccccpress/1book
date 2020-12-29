package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	ydfanyi "github.com/hnmaonanbei/go-youdao-fanyi"

	"github.com/mattn/go-runewidth"
)

func main() {

	//如果啥也没加，就提示用法
	if len(os.Args) == 1 {
		fmt.Println("用法Ⅰ：.\\main.exe + 词汇表.txt")
		fmt.Println("用法Ⅱ：.\\main.exe + 英语单词")
		os.Exit(1)
	}

	//如果过短或者不以.txt结尾，那么就是查询单词
	if len(os.Args[1]) < 5 || os.Args[1][(len(os.Args[1])-4):] != ".txt" {
		fmt.Println(trans(os.Args[1]))
		os.Exit(2)
	}

	//读取词汇表&打乱设置顺序
	tt := strings.Fields(read(os.Args[1]))
	len := len(tt)
	arr := generate(len)

	//像PPT一样显示
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

	displayTxt(s, "开始背单词", 0)
	time.Sleep(time.Second)
	i := 0
	j := 0
	zh := []string{}

	//开始循环单词表
	for {
		if i == len {
			displayTxt(s, "放映结束，按任意键退出", 0)
		}
		if j == 0 && i != len {
			displayTxt(s, tt[arr[i]], arr[i]+1) //显示单词
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
					displayMore(s, zh, arr[i]+1) //显示解释
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

//生成乱序
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

//来自tcell的demo
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
func displayTxt(s tcell.Screen, str string, number int) {
	w, h := s.Size()
	s.Clear()
	emitStr(s, w/2-len(str)/2, h/2-1, tcell.StyleDefault, str)
	if number != 0 {
		emitStr(s, 1, 1, tcell.StyleDefault, strconv.Itoa(number))
	}
	s.Show()
}
func displayMore(s tcell.Screen, str []string, number int) {
	_, h := s.Size()
	s.Clear()
	for i, n := range str {
		emitStr(s, 1, h/2-1+i, tcell.StyleDefault, n)
	}
	if number != 0 {
		emitStr(s, 1, 1, tcell.StyleDefault, strconv.Itoa(number))
	}
	s.Show()
}

//有道翻译
func trans(str string) []string {
	opts := ydfanyi.NewOptions("", "", "")
	res, _ := ydfanyi.Do(str, opts)
	// zh := strings.Join(res.SmartResult.Entries, "")
	zh := res.SmartResult.Entries
	return zh
}
