package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func ExecuteCommand(story StoryConfig, cmd *cobra.Command) error {
	slog.Debug("Sending request to Tines tenant", "Command", story.CommandName, "webhook", story.URL)

	args := convertFlagToArgs(cmd)

	err := sendCommand(story.URL, args)
	if err != nil {
		slog.Error("Failed to send command.", "error", err)

		return err
	}

	return nil
}

func convertFlagToArgs(cmd *cobra.Command) map[string]string {
	args := make(map[string]string)

	if cmd.HasFlags() {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			slog.Debug("Flags provided with command", "name", f.Name, "value", f.Value.String())

			args[f.Name] = f.Value.String()
		})
	}

	return args
}

func sendCommand(url string, args map[string]string) error {
	// Http Request
	resp, err := httpRequest(url, args)
	if err != nil {
		return err
	}

	// Response parse
	err = processResponse(resp)
	if err != nil {
		return err
	}

	return nil
}

func httpRequest(url string, body map[string]string) ([]byte, error) {
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		slog.Error("Failed to create a valid request", "error", err)

		return []byte{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "tCLI-client/0.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	// TODO: Check response code.
	d, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error in processResponse", "error", err)
	}

	return d, nil
}

// processResponse will output the content of the body. If JSON this should be prettyPrinted.
// Once the response definition is build into the tCLI schema this will take that object and process the response in accordance.
// Additionally, we may want to create user created customisation which will layer over the default. i.e. colours etc.
func processResponse(data []byte) error {
	/*
		TODO: Implement child logger of type JSON.
			We will need to expand what this func takes to ensure we have to output opts.

		BUGFIX: Check for existing newline before adding another in. Needs to be done after the above.
	*/

	slog.Info(string(data))

	return nil
}
