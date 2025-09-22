package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcphost/sdk"
)

func main() {
	ctx := context.Background()

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

	fmt.Println("/quit to exit")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		prompt := scanner.Text()
		if prompt == "/quit" {
			break
		}

		subCtx, cancel := context.WithCancel(ctx)
		_, err := host.PromptWithCallbacks(
			subCtx,
			prompt,
			func(name, args string) {
				fmt.Printf("Run tool: %s with args: %s\n", name, args)
				fmt.Print("Are you sure? [y/n] ")
				scanner.Scan()
				if scanner.Text() != "y" {
					cancel()
				}
			},
			func(name, args, result string, isError bool) {
				if isError {
					fmt.Printf("Tool %s failed\n", name)
					return
				}
				if errors.Is(subCtx.Err(), context.Canceled) {
					fmt.Printf("Tool %s was canceled\n", name)
					return
				}
				fmt.Printf("Tool %s ran successfully\n", name)
			},
			func(chunk string) {
				fmt.Print(chunk)
			})
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Fatal(err)
		}

		fmt.Println()
	}
}
