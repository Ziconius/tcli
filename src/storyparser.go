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
	for x, err := range tinesObject.SDK.ListStories(context.Background(), lf) {
		if err != nil {
			slog.Error("Failed to decode", "error", err)
		}
		if slices.Contains(x.Tags, "tcli") {
			s, err := GetStoryConfig(tinesObject, x)
			if err != nil {
				continue
			} // TODO
			stories = append(stories, s)
		}
	}

	return stories
}

func GetStoryConfig(api connector.TinesAPI, story tines.Story) (StoryConfig, error) {
	// We error if we cannot file the "tcli note"
	// sdk does not work as notes are not yet integrated.
	_, err := api.API.GetNotes()
	// story.
	if err != nil {
		return StoryConfig{}, err
	}

	return StoryConfig{
		StoryID: story.ID,
	}, nil
}
