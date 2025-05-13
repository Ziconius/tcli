package main

import (
	"context"
	"log/slog"
	"slices"
	"tcli/src/connector"

	"github.com/tines/go-sdk/tines"
)

func GetAllStoryConfigs(tinesObject connector.TinesAPI) []StoryConfig {
	// BUG: SDK ListFilter isn't working as expecting as no filtering occuring.
	lf := tines.ListFilter{
		Tags: []string{"tcli"},
	}

	stories := []StoryConfig{}
	for x, err := range tinesObject.SDK.ListStories(context.Background(), lf) {
		if err != nil {
			slog.Error("Failed to decode", "error", err)
		}
		
		// TODO: Remove slices.Contains once the tines.ListFilter issue is resolved.
		if slices.Contains(x.Tags, "tcli") {
			sc := StoryConfig{
				StoryID: x.ID,
			}
			stories = append(stories, sc)
		}
	}

	return stories
}
