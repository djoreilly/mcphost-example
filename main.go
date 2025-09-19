package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcphost/sdk"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	options := sdk.Options{
		Model:      "ollama:qwen2.5:3b",
		ConfigFile: "config.json",
		Streaming:  true,
	}
	host, err := sdk.New(ctx, &options)
	if err != nil {
		log.Fatal(err)
	}
	defer host.Close()

	// response, err := host.Prompt(ctx, "what files are in the current directory?")
	response, err := host.PromptWithCallbacks(
		ctx,
		"what files are in the current directory?",
		func(name, args string) {
			fmt.Printf("Run tool: %s with args: %s\n", name, args)
			fmt.Print("Are you sure? [y/n]")
			var answer string
			fmt.Scan(&answer)
			if answer != "y" {
				cancel()
			}
		},
		func(name, args, result string, isError bool) {
			if isError {
				fmt.Printf("Tool %s failed\n", name)
			} else {
				fmt.Printf("Tool %s ran successfully\n", name)
			}
		},
		func(chunk string) {
			fmt.Print(chunk)
		})
	if err != nil {
		log.Fatal(err)
	}

	_ = response
	// fmt.Printf("Final response: %s\n", response)
}
