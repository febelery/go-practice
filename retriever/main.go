package main

import (
	"fmt"
	"learn/retriever/mock"
	real2 "learn/retriever/real"
	"time"
)

type Retriever interface {
	Get(url string) string
}

type Poster interface {
	Post(url string, form map[string]string) string
}

type RetrieverPoster interface {
	Retriever
	Poster
}

const url = "http://www.baidu.com"

func download(r Retriever) string {
	return r.Get(url)
}

func post(poster Poster) {
	poster.Post(url, map[string]string{
		"name": "yes",
		"sex":  "male",
	})
}

func session(s RetrieverPoster) string {
	s.Post(url, map[string]string{
		"contents": "this is a faker",
	})

	return s.Get(url)
}

func inspect(r Retriever) {
	fmt.Println("Inspecting", r)
	fmt.Printf(" > Type:%T Value:%v\n", r, r)
	fmt.Print(" > Type switch: ")

	switch v := r.(type) {
	case *mock.Retriever:
		fmt.Println("Contents:", v.Contents)
	case *real2.Retriever:
		fmt.Println("UserAgent:", v.UserAgent)
	}

	fmt.Println()
}

func main() {
	var r Retriever

	mockRetriever := mock.Retriever{
		Contents: "this is a faker",
	}
	r = &mockRetriever
	inspect(r)

	// Type assertion
	if mockRetriever, ok := r.(*mock.Retriever); ok {
		fmt.Println(mockRetriever.Contents)
	} else {
		fmt.Println("r is not a mock retriever")
	}

	fmt.Println("Try a session with mockRetriever")
	fmt.Println(session(&mockRetriever))

	// Real retriever
	realRetriever := real2.Retriever{
		UserAgent: "UA",
		TimeOut:   time.Second * 2,
	}
	t := &realRetriever
	inspect(t)
	fmt.Println(download(t))

}
