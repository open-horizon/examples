package main

import (
	"fmt"
	"github.com/Shopify/sarama"		// doc: https://godoc.org/github.com/Shopify/sarama
)

func main() {
	encodedStr := sarama.StringEncoder("my message")
	fmt.Println(encodedStr)
}