package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var fileName = "test.txt"
var clientCount = 0
var duration = time.Second * 4

func main() {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	mu := &sync.RWMutex{}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		clientCount++
		go handleConnection(conn, fileName, mu)
	}
}

func handleConnection(conn net.Conn, f string, mutex *sync.RWMutex) {
	curCount := clientCount
	name := conn.RemoteAddr().String()
	fmt.Printf("%+v connected\n", name)
	fmt.Fprintln(conn, "Hello, "+name+"\n\r")
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
LOOP:
	for scanner.Scan() {
		text := scanner.Text()
		args := strings.Split(text, " ")
		comm := args[0]
		switch comm {
		case "exit":
			fmt.Fprintf(conn, "bye!\n\r")
			fmt.Println(name, "disconnected")
			break LOOP
		case "read":
			var whatId = ""
			if len(args) > 1 {
				whatId = args[1]
			}
			fmt.Fprintf(conn, read(f, whatId, mutex, curCount)+"\n\r")
			fmt.Println(name, "read file")
		case "delete":
			var whatId = ""
			if len(args) > 1 {
				whatId = args[1]
			}
			fmt.Fprintf(conn, deleteRecord(f, whatId, mutex)+"\n\r")
			fmt.Println(name, "delete something")
		case "change":
			whatId := args[1]
			rec := strings.Join(args[2:], " ")
			fmt.Fprintf(conn, changeRecord(f, whatId, rec, mutex)+"\n\r")
			fmt.Println(name, "change record")
		case "add":
			var whatId = ""
			var rec = ""
			if strings.Contains(strings.Join(args, " "), "-id") {
				whatId = args[2]
				rec = strings.Join(args[3:], " ")
			} else {
				rec = strings.Join(args[1:], " ")
			}
			fmt.Fprintf(conn, addRecord(f, whatId, rec, mutex)+"\n\r")
			fmt.Println(name, "add record")
		default:
			fmt.Println(name, "enters", text)
			fmt.Fprintln(conn, "You enter "+text+"\n\r")
		}
	}
}

func read(f string, id string, mutex *sync.RWMutex, count int) string {
	mutex.RLock()
	defer mutex.RUnlock()
	fileBuffer, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		indent := ""
		for i := 0; i < count-1; i++ {
			indent += "	"
		}
		i := math.Round(duration.Seconds())
		for range ticker.C {
			i--
			fmt.Println("reader #", count, indent, i)
		}
	}(ticker)
	time.Sleep(duration)
	ticker.Stop()
	inputData := string(fileBuffer)
	if inputData == "" {
		return "file is empty"
	}
	if id != "" {
		lines := strings.Split(inputData, "\n")
		newId, _ := strconv.Atoi(id)
		if newId > len(lines) || newId < 1 {
			return "invalid index"
		} else {
			return lines[newId-1]
		}
	} else {
		return inputData
	}
}

func deleteRecord(f string, id string, mutex *sync.RWMutex) string {
	mutex.Lock()
	defer mutex.Unlock()
	fileBuffer, _ := ioutil.ReadFile(f)
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		i := math.Round(duration.Seconds())
		for range ticker.C {
			i--
			fmt.Println("deleting...", i)
		}
	}(ticker)
	time.Sleep(duration)
	ticker.Stop()
	data := string(fileBuffer)
	if data == "" {
		return "file is empty"
	}
	file, err := os.OpenFile(f, os.O_RDWR, 7777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if id == "" {
		file.Truncate(0)
		return "success"
	}
	lines := strings.Split(data, "\n")
	newId, _ := strconv.Atoi(id)
	if newId > len(lines) || newId < 1 {
		return "invalid index"
	}
	lines = append(lines[:newId-1], lines[newId:]...)
	data = strings.Join(lines, "\n")
	file.Truncate(0)
	file.WriteAt([]byte(data), 0)
	return "success"
}

func changeRecord(f string, id string, record string, mutex *sync.RWMutex) string {
	mutex.Lock()
	defer mutex.Unlock()
	fileBuffer, _ := ioutil.ReadFile(f)
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		i := math.Round(duration.Seconds())
		for range ticker.C {
			i--
			fmt.Println("changing...", i)
		}
	}(ticker)
	time.Sleep(duration)
	ticker.Stop()
	data := string(fileBuffer)
	if data == "" {
		return "file is empty"
	}
	lines := strings.Split(data, "\n")
	newId, _ := strconv.Atoi(id)
	if newId > len(lines) || newId < 1 {
		return "invalid index"
	}
	lines[newId-1] = record
	data = strings.Join(lines, "\n")
	file, err := os.OpenFile(f, os.O_RDWR, 7777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Truncate(0)
	file.WriteAt([]byte(data), 0)
	return "success"
}

func addRecord(f string, id string, record string, mutex *sync.RWMutex) string {
	mutex.Lock()
	defer mutex.Unlock()
	file, err := os.OpenFile(f, os.O_RDWR|os.O_APPEND, 7777)
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		i := math.Round(duration.Seconds())
		for range ticker.C {
			i--
			fmt.Println("adding...", i)
		}
	}(ticker)
	time.Sleep(duration)
	ticker.Stop()
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if id == "" {
		info, _ := file.Stat()
		if info.Size() == 0 {
			file.Write([]byte(record))
		} else {
			file.Write([]byte("\n" + record))
		}
		return "You write: " + record
	} else {
		fileBuffer, _ := ioutil.ReadFile(f)
		data := string(fileBuffer)
		lines := strings.Split(data, "\n")
		newId, _ := strconv.Atoi(id)
		if newId > len(lines) || newId < 1 {
			return "invalid index"
		}
		var temp = make([]string, len(lines[newId-1:]))
		copy(temp, lines[newId-1:])
		lines = append(lines[:newId-1], record)
		lines = append(lines, temp...)
		data = strings.Join(lines, "\n")
		file.Truncate(0)
		file.WriteAt([]byte(data), 0)
		return "You write: " + record
	}
}
