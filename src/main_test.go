package main

import (
	"log"
	"log/slog"
	"testing"

	"github.com/tines/go-sdk/tines"
)

func TestValidateState(t *testing.T) {
	auth, _ := AuthConfig()
	sdk, err := tines.NewClient(
		tines.SetTenantUrl(auth.TenantURL),
		tines.SetApiKey(auth.APIKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Local API
	api, err := NewTinesAPI(auth.TenantURL, auth.APIKey)
	if err != nil {
		slog.Error("failed to create API")
	}
	tines := Tines{
		API: api,
		SDK: sdk,
	}
	sc := StoredConfig{}
	err = UpdateConfigCache(tines, &sc)
	if err != nil {
		t.Logf("Failed to update, err: %s\n", err)
	}
	err = sc.WriteConfig()
	if err != nil {
		t.Logf("Failed to write config, err: %s\n", err)
	}
}
