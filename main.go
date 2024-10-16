package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	env, err := LoadConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	oauthConfig := getOAuthConfig(env)
	app := GitHub{oauthConfig}

	mux := http.NewServeMux()
	mux.HandleFunc("/redirect", app.Redirect)
	mux.HandleFunc("/callback", app.Callback)

	addr := fmt.Sprintf(":%s", env.Port)
	logger.Info("server listening", "addr", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func getOAuthConfig(env *Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     env.Github.ClientID,
		ClientSecret: env.Github.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  env.RedirectURL,
		Scopes:       []string{"repo", "user"},
	}
}
