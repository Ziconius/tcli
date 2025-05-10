package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func ExecuteCommand(story StoryConfig, cmd *cobra.Command) error {
	fmt.Printf("... Sending: %v: %v\n", story.CommandName, story.URL)
	
	fmt.Printf("Flags: %#v\n", cmd.Flags())
	// bb := cmd.Flags()
	// for _, x := range bb. {

	// }



	resp, err := sendCommand(story.URL, cmd.ValidArgs)
	if err != nil {
		fmt.Printf("Failed to run %s\n", cmd.DisplayName())
	}
	fmt.Printf("Response: %s\n", resp)

	return nil
}


func sendCommand(url string, args []cobra.Completion) (string, error){
	// Args to Body
	buildRequestBody(args)
	// Http Request
	httpRequest(url)
	// Response parse

	return "", nil
}


func buildRequestBody(args []cobra.Completion) ([]byte, error){
	fmt.Println("Printing args & making body.")
	for i, x := range args {
		fmt.Printf("%d: %v\n", i, x)
	}

	return []byte("Testing"), nil
}

func httpRequest(url string) (string, error){
	// Need to convert to send POST requests and proper headers.
	resp, err :=http.Get(url)
	if err != nil {
		return "", err
	}
	fmt.Printf("Resp: %v", resp.Status)
	
	return "Ok", nil
}