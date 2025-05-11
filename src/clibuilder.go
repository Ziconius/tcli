package main

import (
	"tcli/src/connector"

	"github.com/spf13/cobra"
)

func InitCLI() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "tcli",
		Short: "tCLI  - A tenant managed CLI tool for Tines.",
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	// TODO: Add other global arg, i.e. config
	// TODO: Add global flag for verbosity.

	return rootCmd
}

func CmdConfig(tApi connector.TinesAPI, cache StoredConfig)*cobra.Command{
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "manage tCLI config",
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) {
			// TODO: Convert to proper args.
			if args[0] == "pull"{
				UpdateCache(tApi, &cache)
			}
		},
	}

	return configCmd
}

func storySubCommand(x StoryConfig) *cobra.Command {
	var tmp = &cobra.Command{
		Use: x.CommandName,
		// Aliases: []string{"command"},
		Short: x.Description,
		// Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteCommand(x, cmd)
		},
	}
	for _, flag := range x.Request {
		tmp.Flags().String(flag, "", "testing")
	}

	return tmp
}

func BuildCliParser(sc StoredConfig) *cobra.Command {
	// Based command
	var tinesCommand = &cobra.Command{
		Use:     "cmd",
		Aliases: []string{"command"},
		Short:   "executes a tines story from the tenant.",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// fmt.Printf("PArty: %#v\n", sc)
	for _, x := range sc.Commands {
		// fmt.Printf("Command: %v\n", x.CommandName)
		g := storySubCommand(x)
		tinesCommand.AddCommand(g)
	}

	return tinesCommand
}

/*
TODO:
Reserved commands

Info - returns tenant URL + other info
List - lists the name of all commands we have
Update - pulls the most recent config and displays results.

global flags
 --no-cache, does not use or modify the local cache
 --no-validate: Does not run the validation command.

*/
