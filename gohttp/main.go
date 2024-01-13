package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func f651() {
	listener, _ := net.Listen("tcp", "localhost:8888")
	fmt.Println("Server is running")
	for {
		conn, _ := listener.Accept()
		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			request, _ := http.ReadRequest(bufio.NewReader(conn))
			dump, _ := httputil.DumpRequest(request, true)
			fmt.Println(string(dump))
			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body:       io.NopCloser(strings.NewReader("Hello world\n")),
			}
			response.Write(conn)
			conn.Close()
		}()
	}
}

func f652() {
	conn, _ := net.Dial("tcp", "localhost:8888")
	request, _ := http.NewRequest("GET", "localhost:8888", nil)
	request.Write(conn)
	response, _ := http.ReadResponse(bufio.NewReader(conn), request)
	dump, _ := httputil.DumpResponse(response, true)
	fmt.Println(string(dump))
}

func f661() {
	listener, _ := net.Listen("tcp", "localhost:8888")
	fmt.Println("Server is running")
	for {
		conn, _ := listener.Accept()
		go func() {
			defer conn.Close()
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			for {
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					neterr, ok := err.(net.Error)
					if ok && neterr.Timeout() {
						fmt.Println("タイムアウト")
						break
					} else if err == io.EOF {
						break
					}
					panic(err)
				}
				dump, _ := httputil.DumpRequest(request, true)
				fmt.Println(string(dump))
				content := "Hello world"
				response := http.Response{
					StatusCode:    200,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)),
					Body:          io.NopCloser(strings.NewReader(content)),
				}
				response.Write(conn)
			}
		}()
	}
}

func f662() {
	sendMessages := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}
	current := 0
	var conn net.Conn = nil
	for {
		var err error
		if conn == nil {
			conn, _ = net.Dial("tcp", "localhost:8888")
			fmt.Printf("Access: %d\n", current)
		}
		request, _ := http.NewRequest(
			"POST",
			"http://localhost:8888",
			strings.NewReader(sendMessages[current]),
		)
		err = request.Write(conn)
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			fmt.Println("リトライ")
			conn = nil
			continue
		}
		dump, _ := httputil.DumpResponse(response, true)
		fmt.Println(string(dump))
		current++
		if current == len(sendMessages) {
			break
		}
	}
	conn.Close()
}

var contents = []string{
	"これは、私わたしが小さいときに、村の茂平もへいというおじいさんからきいたお話です。",
	"村の茂平は、今から四十年ぐらいまえのことです。",
	"茂平は、そのころ、まだ若い男でした。",
	"ある日、茂平は、山へしばかりにいきました。",
	"山へしばかりにいくと、山の中の小さな池のほとりに、",
	"小さな小さな家がありました。",
	"茂平は、その家の戸をたたきました。",
}

func processSession681(conn net.Conn) {
	defer conn.Close()
	for {
		req, _ := http.ReadRequest(bufio.NewReader(conn))
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
		fmt.Fprintf(conn, strings.Join([]string{
			"HTTP/1.1 200 OK",
			"Content-Type: text/plain; charset=utf-8",
			"Transfer-Encoding: chunked",
			"", "",
		}, "\r\n"))
		for _, content := range contents {
			bytes := []byte(content)
			fmt.Fprintf(conn, "%x\r\n%s\r\n", len(bytes), content)
		}
		fmt.Fprintf(conn, "0\r\n\r\n")
	}
}

func f681() {
	listener, _ := net.Listen("tcp", "localhost:8888")
	fmt.Println("Server is running")
	for {
		conn, _ := listener.Accept()
		go processSession681(conn)
	}
}

// 料理配膳ロボット
func writeToConn(sessionResponses chan chan *http.Response, conn net.Conn) {
	defer conn.Close()
	//注文票管理システムを常に見張っている
	for sessionResponse := range sessionResponses {
		//料理が完成していたら
		response := <-sessionResponse
		//客に届けにいく
		response.Write(conn)
		close(sessionResponse)
	}
}

// 料理を準備する
func handleRequest691(request *http.Request, resultReceiver chan *http.Response) {
	dump, _ := httputil.DumpRequest(request, true)
	fmt.Println(string(dump))
	content := "Hello world"
	response := http.Response{
		StatusCode:    200,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: int64(len(content)),
		Body:          io.NopCloser(strings.NewReader(content)),
	}
	//料理が完成したら注文票の上に料理を置く
	resultReceiver <- &response
}

func processSession691(conn net.Conn) {
	//注文票の管理システムの作成
	sessionResponses := make(chan chan *http.Response, 50)
	defer close(sessionResponses)
	//料理配膳ロボット
	go writeToConn(sessionResponses, conn)
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		//客から注文
		request, _ := http.ReadRequest(reader)
		//注文票を作成
		sessionResponse := make(chan *http.Response)
		//注文票管理システムに注文票を登録
		sessionResponses <- sessionResponse
		//1人の料理人に依頼（料理人はシステムを知らない）
		go handleRequest691(request, sessionResponse)
	}
}

func f691() {
	listener, _ := net.Listen("tcp", "localhost:8888")
	fmt.Println("Server is running")
	for {
		conn, _ := listener.Accept()
		//テーブルへ案内
		go processSession691(conn)
	}
}

func main() {
	// tcpを使ったHTTPサーバー
	f651()
	// tcpを使ったHTTPクライアント
	f652()
	// keep-alive対応 HTTPサーバー
	f661()
	// keep-alive対応 HTTPクライアント
	f662()
	// チャンク形式のサーバ
	f681()
	// パイプライニングのサーバ実装
	f691()
}
