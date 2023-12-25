package main

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "http.ResponseWriter")
}
func question1() {
	file, err := os.Create("question1.txt")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(file, "question%d", 1)
}
func question2() {
	file, _ := os.Create("output.csv")
	defer file.Close()
	records := [][]string{
		{"name", "age", "city"},
		{"Alice", "23", "New York"},
		{"Bob", "34", "San Francisco"},
	}
	// csvWriter := csv.NewWriter(os.Stdout)
	csvWriter := csv.NewWriter(file)
	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			panic(err)
		}
	}
	csvWriter.Flush()
}
func handler2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")
	source := map[string]string{
		"name": "Taiyo",
	}
	gz := gzip.NewWriter(w)
	defer gz.Close()

	mw := io.MultiWriter(gz, os.Stdout)
	encoder := json.NewEncoder(mw)
	encoder.Encode(source)
}
func question3() {
	http.HandleFunc("/", handler2)
	http.ListenAndServe(":8080", nil)
}
func main() {
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	file.Write([]byte("os.File example\n"))
	file.Close()

	os.Stdout.Write([]byte("os.Stdout.Write example\n"))

	var buffer bytes.Buffer
	buffer.Write([]byte("buffer.Write example"))
	fmt.Println(buffer.String())

	var buffer2 bytes.Buffer
	buffer2.WriteString("buffer2.WriteString\n")
	fmt.Println(buffer2.String())

	// var buffer3 bytes.Buffer
	io.WriteString(&buffer2, "io.WriteString")
	fmt.Println(buffer2.String())

	//2.4.4書き出し専用
	var builder strings.Builder
	builder.Write([]byte("strings.Builder\n"))
	fmt.Println(builder.String())

	//2.4.5
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	io.WriteString(conn, "GET / HTTP/1.0\r\nHost: example.com\r\n\r\n")
	io.Copy(os.Stdout, conn)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	req.Write(conn)
	// io.Copy(os.Stdout, conn)

	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)

	//2.4.6
	file2, err := os.Create("multiwriter.txt")
	if err != nil {
		panic(err)
	}
	writer := io.MultiWriter(file2, os.Stdout)
	io.WriteString(writer, "io.MultiWriter\n")

	file3, err := os.Create("test.txt.gz")
	if err != nil {
		panic(err)
	}
	writer3 := gzip.NewWriter(file3)
	writer3.Header.Name = "gzip.txt"
	io.WriteString(writer3, "gzip.Writer")
	writer3.Close()

	//2.4.7
	fmt.Fprintf(os.Stdout, "Write os.stdout sy %v \n\n\n", time.Now())

	fmt.Println("--Question--")
	question1()
	question2()
	question3()
}
