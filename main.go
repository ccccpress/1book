package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Resp ...
type Resp struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func main() {

	txt := read()

	tt := strings.Fields(txt)

	len := len(tt)

	arr := generate(len)

	b := make([]byte, 1)

	for i := range arr {
		fmt.Print(tt[arr[i]])
		zh := baidu(tt[arr[i]])
		os.Stdin.Read(b)
		fmt.Println(zh)
		os.Stdin.Read(b)
	}
	fmt.Println("-----")
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

	fmt.Println("-----")
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

func baidu(q string) string {

	baidu := "https://fanyi-api.baidu.com/api/trans/vip/translate"
	//申请一个api填在下面
	appid := "********"
	passport := "***********"
	salt := strconv.FormatInt(time.Now().Unix(), 10)
	get := baidu + "?q=" + q + "&from=en&to=zh&appid=" + appid + "&salt=" + salt + "&sign=" + md5V(appid+q+salt+passport)

	client := &http.Client{}
	resp, err := client.Get(get)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(body))
	//fmt.Println(reflect.TypeOf(body))
	resp2 := Resp{}

	err = json.Unmarshal(body, &resp2)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(resp2)
	return resp2.TransResult[0].Dst
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
