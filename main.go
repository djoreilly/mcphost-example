package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcphost/sdk"
)

func main() {
	ctx := context.Background()

	options := sdk.Options{
		Model:      "ollama:qwen2.5:3b",
		ConfigFile: "config.json",
	}
	host, err := sdk.New(ctx, &options)
	if err != nil {
		log.Fatal(err)
	}
	defer host.Close()

	response, err := host.Prompt(ctx, "what files are in the current directory?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
