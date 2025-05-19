package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (story *StoryConfig) ExecuteCommand(tenant string, cmd *cobra.Command) error {
	slog.Debug("Sending request to Tines.", "Command", story.CommandName, "webhook", story.Path, "tenant", tenant)

	args := story.convertFlagToArgs(cmd)

	err := story.sendCommand(tenant, args)
	if err != nil {
		slog.Error("Failed to send command.", "error", err)

		return err
	}

	return nil
}

func (story *StoryConfig) convertFlagToArgs(cmd *cobra.Command) map[string]string {
	args := make(map[string]string)

	if cmd.HasFlags() {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			slog.Debug("Flags provided with command", "name", f.Name, "value", f.Value.String())

			args[f.Name] = f.Value.String()
		})
	}

	return args
}

func (story *StoryConfig) sendCommand(host string, args map[string]string) error {
	// Http Request
	url := host + story.Path
	resp, err := story.httpRequest(url, args)
	if err != nil {
		return err
	}

	// Response parse
	err = story.processResponse(resp)
	if err != nil {
		return err
	}

	return nil
}

func (story *StoryConfig) httpRequest(url string, body map[string]string) (*http.Response, error) {
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest(story.Request.Method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		slog.Error("Failed to create a valid request", "error", err)

		return &http.Response{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "tCLI-client/0.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return res, nil
}

// processResponse will output the content of the body. If JSON this should be prettyPrinted.
// Once the response definition is build into the tCLI schema this will take that object and process the response in accordance.
// Additionally, we may want to create user created customisation which will layer over the default. i.e. colours etc.
func (story *StoryConfig) processResponse(resp *http.Response) error {
	/*
		TODO: Implement child logger of type JSON.
			We will need to expand what this func takes to ensure we have to output opts.

		BUGFIX: Check for existing newline before adding another in. Needs to be done after the above.
	*/

	/*
		if response.Format is not present, or is auto then decide, else.
	*/

	format := story.getOutputFormat(resp.Header.Get("content-type"))

	data, _ := io.ReadAll(resp.Body)

	if format == "text" {
		slog.Info(string(data))
	} else if format == "json" {
		// This needs to be passed to a JSON  handler.
		var output bytes.Buffer
		json.Indent(&output, data, "", "\t")
		slog.Info(output.String())
	}

	return nil
}

// Takes the content-type header as a string, and returns a valid tCLI output format for response printing.
// If a default is set this is always returned, if none is present the output format is based on the webhooks MIMEType.
//
// If no MIMEType can be determined, or an unsupported type is used we default to text output.
func (story *StoryConfig) getOutputFormat(contentType string) string {
	// Currently supported output
	supportedFormats := []string{"text", "json"}

	//
	if slices.Contains(supportedFormats, story.Response.Format) {
		return story.Response.Format
	}

	// If no response format if defined, fall back to using auto parsing based on content-type.
	// Defaults to text if no content type can be determined.
	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == "application/json" {
			return "json"
		} else if t == "text/plain" {
			return "text"
		}
	}

	return "text"
}
