package main

import (
	"context"
	"log/slog"
	"slices"
	"tcli/src/connector"

	"github.com/tines/go-sdk/tines"
)

func GetAllStoryConfigs(tinesObject connector.TinesAPI) []StoryConfig {
	// TODO: Not working
	lf := tines.ListFilter{
		Tags: []string{"tcli"},
	}

	stories := []StoryConfig{}
	// BUG: SDK ListFilter isn't working as expecting and no filtering occuring.
	for x, err := range tinesObject.SDK.ListStories(context.Background(), lf) {
		if err != nil {
			slog.Error("Failed to decode", "error", err)
		}

		if slices.Contains(x.Tags, "tcli") {
			sc := StoryConfig{
				StoryID: x.ID,
			}
			stories = append(stories, sc)
		}
	}

	return stories
}
