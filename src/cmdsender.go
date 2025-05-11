package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// TODO: Remove process args into a func.
func ExecuteCommand(story StoryConfig, cmd *cobra.Command) error {
	fmt.Printf("... Sending: %v: %v\n", story.CommandName, story.URL)

	args := make(map[string]string)

	if cmd.HasFlags() {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			slog.Debug("Flags provided with command", "name", f.Name, "value", f.Value.String())

			args[f.Name] = f.Value.String()
		})
	}

	resp, err := sendCommand(story.URL, args)
	if err != nil {
		fmt.Printf("Failed to run %s\n", cmd.DisplayName())
	}
	fmt.Printf("Response: %s\n", resp)

	return nil
}


func sendCommand(url string, args map[string]string) (string, error) {
	// Args to Body
	// buildRequestBody(args)

	// Http Request
	resp, _ := httpRequest(url, args)

	// Response parse
	processResponse(resp)

	return "", nil
}

// func buildRequestBody(args []cobra.Completion) ([]byte, error) {
// 	fmt.Println("Printing args & making body.")
// 	for i, x := range args {
// 		fmt.Printf("%d: %v\n", i, x)
// 	}

// 	return []byte("Testing"), nil
// }

func httpRequest(url string, body map[string]string) (*http.Response, error) {
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err) // TODO: Replace
	}

	defer res.Body.Close()

	// TODO: Output response
	fmt.Printf("Body: %v\n", res.Body)


	return res, nil

}

// processResponse will output the content of the body. If JSON this should be prettyPrinted.
// Once the response definition is build into the tCLI schema this will take that object and process the response in accordance.
// Additionally, we may want to create user created customisation which will layer over the default. i.e. colours etc.
func processResponse(resp *http.Response) error {


	fmt.Print()
	return nil
}