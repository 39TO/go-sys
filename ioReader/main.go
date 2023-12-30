package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func f341() {
	for {
		buffer := make([]byte, 5)
		size, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		fmt.Printf("size=%d, input=%s \n", size, string(buffer))
	}
}

func f342() {
	file, err := os.Open("main.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(os.Stdout, file)
}

func f343() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	io.WriteString(conn, "GET / HTTP/1.0\r\nHost: example.com\r\n\r\n")
	res, err := http.ReadResponse(bufio.NewReader(conn), nil)

	fmt.Println(res.Header)
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}

func f351() {
	r := strings.NewReader("abcdefghijklmnopqrstuvwxyz\n")
	sr := io.NewSectionReader(r, 14, 7)
	io.Copy(os.Stdout, sr)
}

// エンディアン変換
func f352() {
	//100000(10)
	data := []byte{0x0, 0x0, 0x27, 0x10}
	var i int32
	// to リトル
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Printf("data: %d\n", i)
}

func f353() {
	file, err := os.Open("PNG_transparency_demonstration_1.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	chunks := f353ReadChunks(file)
	for _, c := range chunks {
		f353DumpChunks(c)
	}
}

func f353ReadChunks(file *os.File) []io.Reader {
	var chunks []io.Reader

	file.Seek(8, 0) //内部ポインタを移動させる
	var offset int64 = 8

	for {
		var length int32
		err := binary.Read(file, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(length)+12))
		offset, _ = file.Seek(int64(length)+8, 1)
	}
	return chunks
}

func f353DumpChunks(c io.Reader) {
	//チャンクの長さ
	var length int32
	binary.Read(c, binary.BigEndian, &length) //内部ポインタが種類に移動する
	//チャンクの種類
	b := make([]byte, 4)
	c.Read(b)
	fmt.Printf("chunk %v (%d bytes)\n", string(b), length)
}

func f354() {
	file, err := os.Open("PNG_transparency_demonstration_1.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	newFile, err := os.Create("PNG_transparency_demonstration_secret.png")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	chunks := f353ReadChunks(file)
	io.WriteString(newFile, "\x89PNF\r\n\x1a\n")
	io.Copy(newFile, chunks[0])
	io.Copy(newFile, f354TextChunk("Lambda Note++"))
	for _, c := range chunks[1:] {
		io.Copy(newFile, c)
	}
}

func f354TextChunk(text string) io.Reader {
	byteText := []byte(text)
	crc := crc32.NewIEEE()
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, int32(len(byteText))) //チャンクの長さ
	writer := io.MultiWriter(&buffer, crc)
	io.WriteString(writer, "teXt")                       //チャンクの種類
	writer.Write(byteText)                               //チャンクのデータ
	binary.Write(&buffer, binary.BigEndian, crc.Sum32()) //チャンクのCRC
	return &buffer
}

func f354_2() {
	file, err := os.Open("PNG_transparency_demonstration_secret.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	chunks := f353ReadChunks(file)
	for _, c := range chunks {
		f353DumpChunks(c)
	}
}

func q1() {
	file, err := os.Open("old.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	newFile, err := os.Create("new.txt")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	io.Copy(newFile, file)
}

func q2() {
	data := make([]byte, 1024)
	_, _ = rand.Read(data)
	fmt.Printf("data %v", data)
}

func q3() {
	file, _ := os.Create("q3.zip")
	defer file.Close()
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()
	w, _ := zipWriter.Create("newFile.txt")
	w2, _ := zipWriter.Create("newFile2.txt")
	io.WriteString(w, "Writer1")
	io.WriteString(w2, "Writer2")
}

func qa4Handler(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Create("q4.zip")
	defer file.Close()
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()
	writer, _ := zipWriter.Create("q41.txt")
	writer.Write([]byte("Hello, world!"))
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=q4.zip")
	_, err := io.Copy(w, file)
	if err != nil {
		panic(err)
	}
}

func q4() {
	http.HandleFunc("/", qa4Handler)
	http.ListenAndServe(":8080", nil)
}

func q6() {
	var (
		computer    = strings.NewReader("COMPUTER")
		system      = strings.NewReader("SYSTEM")
		programming = strings.NewReader("PROGRAMMING")
	)

	var stream io.Reader
	a := io.NewSectionReader(programming, 5, 1)
	s := io.NewSectionReader(system, 0, 1)
	c := io.NewSectionReader(computer, 0, 1)
	i := io.NewSectionReader(programming, 8, 1)
	stream = io.MultiReader(a, s, c, i)

	io.Copy(os.Stdout, stream)
}

func main() {
	//3.4.1
	// f341()

	//3.4.2
	// f342()

	//3.4.3
	// f343()

	//3.5.1
	// f351()

	//3.5.2
	// f352()

	//3.5.3
	// f353()

	//3.5.4
	// f354()
	// f354_2()

	//question
	// q1()
	// q2()
	// q3()
	// q4()
	q6()
}
