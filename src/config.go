package main

import (
	"log/slog"
)

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
	if err := cache.LoadConfig(); err != nil {
		return cache, err
	}

	return cache, nil
}
