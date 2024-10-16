package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func extractMessage(r *http.Response) (string, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func menu1() {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get("http://localhost:5678/get")
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("req timeout")
		} else {
			panic(err)
		}
		return
	}
	defer resp.Body.Close()
	text, err := extractMessage(resp)
	fmt.Println(text)
	return
}

func menu2() {
	var fileName string
	var content string
	scan := bufio.NewScanner(os.Stdin)
	fmt.Print("fileName: ")
	scan.Scan()
	fileName = scan.Text()
	fmt.Print("content: ")
	scan.Scan()
	content = scan.Text()

	fileName = fileName + ".txt"
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reqBody := &bytes.Buffer{}
	w := multipart.NewWriter(reqBody)
	filefield, err := w.CreateFormFile("file", file.Name())
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(filefield, file)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Post("http://localhost:5678/upload", w.FormDataContentType(), reqBody)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("Req timeout")
		} else {
			panic(err)
		}
		return
	}
	defer resp.Body.Close()
	data, err := extractMessage(resp)
	fmt.Println(data)
	return
}

func menu3() {
	var name string
	var age int
	scan := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Input name: ")
		name, _ = scan.ReadString('\n')
		name = strings.TrimSpace(name)

		if len(name) > 5 {
			break
		} else {
			fmt.Println("Name > 5 characters")
		}
	}
	for {
		fmt.Print("Age: ")
		fmt.Scanf("%d\n", &age)
		if age > 0 {
			break
		} else {
			fmt.Println("Age must > 0 years old")
		}
	}

	person := data.Person{Name: name, Age: age}
	jsonData, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	
	reqBody := &bytes.Buffer{}
	w := multipart.NewWriter(reqBody)
	personField, err := w.CreateFormField("Person")
	if err != nil {
		panic(err)
	}
	_, err = personField.Write(jsonData)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}

	//req timeout
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Post("http://localhost:5678/json", w.FormDataContentType(), reqBody)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("req timeout")
		} else {
			panic(err)
		}
		return
	}
	defer resp.Body.Close()

	data, err := extractMessage(resp)
	fmt.Println(data)
	return
}
func main() {
	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1. Retrieve message")
		fmt.Println("2. Post File")
		fmt.Println("3. Post JSON")
		fmt.Println("4. Exit")
		fmt.Print(">> ")
		scan.Scan()
		choice := scan.Text()
		if choice == "1" {
			menu1()
		} else if choice == "2" {
			menu2()
		} else if choice == "3" {
			menu3()
		} else if choice == "4" {
			fmt.Println("Thank you for using our services")
			break
		} else {
			fmt.Println("Please choose one of the above menu")
		}
	}
}
