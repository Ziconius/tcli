package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"time"

	"tcli/src/connector"
)

var tinesAPI connector.TinesAPI

func init() {
	auth, err := AuthConfig()
	if err != nil {
		slog.Warn("Failed to load auth config file", "error", err)
	}

	tinesAPI, err = connector.NewTinesAPI(auth.TenantURL, auth.APIKey)
	if err != nil {
		slog.Error("failed to create API")
	}
}

func main() {

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	// Build the root of the application.
	cli := InitCLI()

	cache, err := LocalConfig()
	if err != nil {
		// Failed to load a file
	}
	ok, err := cache.ValidateState()
	if err != nil {
		// Checking cache state has failed.
	}
	if !ok {
		// Pull new cache state
		if err = UpdateCache(tinesAPI, &cache); err != nil {
			slog.Error("Failed to update config cache", "error", err)
		}
	}

	// Now that we have a valid API config & created cache we can create config args.
	cli.AddCommand(CmdConfig(tinesAPI, cache))

	// Remote command builder & parser
	bc := BuildCliParser(cache)
	cli.AddCommand(bc)
	err = cli.Execute()
	if err != nil {
		fmt.Printf("Failed to generate CLI arguements, error: %s\n", err)
	}

	// TODO: Response parser
}

func UpdateCache(tinesAPI connector.TinesAPI, cache *StoredConfig) error {
	if err := GetRemoteConfig(tinesAPI, cache); err != nil {
		slog.Error("Failed to update config cache", "error", err)
		panic("failed to update cache from remote, quitting.")
	}
	err := cache.WriteConfig()
	if err != nil {
		fmt.Printf("Warning: failed to save configuration file, error: %s\n", err)

		return err
	}

	return nil
}

func GetRemoteConfig(api connector.TinesAPI, sc *StoredConfig) error {
	// TODO
	slog.Info("Downloading config...")
	stories := GetAllStoryConfigs(api)
	rl, err := api.API.ListResources()
	if err != nil {
		return err
	}

	// Building regex for resource name
	r, err := regexp.Compile(`^tcli_\d{4}$`)
	if err != nil {
		fmt.Printf("Failed to compile Regex...")
		return err
	}

	// tcli resources - Maybe this should be tcli_resources only?
	tResources := []connector.Resource{}
	for _, x := range rl.GlobalResources {
		if r.MatchString(x.Name) {
			tResources = append(tResources, x)
		}
	}

	// Loop
	final := []StoryConfig{}
	for _, y := range stories {
		for _, z := range tResources {
			st := "tcli_" + strconv.Itoa(y.StoryID)
			if st == z.Name {
				// Build story from z.Value.
				err := ValueToStory(z.Value, &y)
				if err != nil {
					fmt.Printf("Failed to unmarshal story(%v), error: %v\n", y.StoryID, z.Slug)
				}
				final = append(final, y)
			}
		}
	}
	sc.Commands = final
	sc.LastUpdated = time.Now()

	return nil
}

// Used to convert from value string to a usable struct.
type CommandCfg struct {
	Cmd         string   `json:"cmd"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Request     []string `json:"request"`
}

func ValueToStory(v string, story *StoryConfig) error {

	vv := []byte(v)
	ccfg := CommandCfg{}
	err := json.Unmarshal(vv, &ccfg)
	if err != nil {
		fmt.Printf("Failed to unmarshal.")
		return err
	}

	story.CommandName = ccfg.Cmd
	story.Description = ccfg.Description
	story.URL = ccfg.URL
	story.Request = ccfg.Request
	return nil
}
