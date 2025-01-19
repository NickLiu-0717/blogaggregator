package main

import (
	"fmt"

	config "github.com/NickLiu-0717/blogaggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("NickLiu")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Read config: %+v\n", cfg)
}
