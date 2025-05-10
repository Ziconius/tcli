package main

import (
	"log/slog"
)

// We should pull config from {USER_HOME}/.tcli/.config ../.auth

func AuthConfig() (Auth, error) {
	auth := Auth{}
	err := auth.LoadConfig()
	if err != nil {
		slog.Error("Failed to load local config file")
		
		return auth, err
	}

	return auth, nil
}

func LocalConfig() (StoredConfig, error) {
	cache := StoredConfig{}
	if err := cache.LoadConfig();err != nil {
		return cache, err
	}
	return cache, nil
}
