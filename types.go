package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type GitHub struct {
	Config *oauth2.Config
}

func (g *GitHub) pkceChallenge() oauth2.AuthCodeOption {
	verifier := oauth2.GenerateVerifier()
	return oauth2.S256ChallengeOption(verifier)
}

// oauthURL devuelve la URL de autorización de GitHub. Esta URL se utiliza para redirigir
// al usuario a la página de autorización de GitHub, donde el usuario concede permisos
// a la aplicación. El parámetro "state" se utiliza para evitar ataques de CSRF.
func (g *GitHub) oauthURL() string {
	return g.Config.AuthCodeURL("state", oauth2.AccessTypeOffline, g.pkceChallenge())
}

// Redirect redirige al usuario a la página de autorización de GitHub. El usuario es
// redirigido a la URL de autorización de GitHub, donde se le pedirá que conceda
// permisos a la aplicación. Después de autorizar, GitHub redirigirá al usuario de
// vuelta a la ruta "/callback" de la aplicación.
func (g *GitHub) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, g.oauthURL(), http.StatusTemporaryRedirect)
}

// Callback maneja la respuesta de GitHub después de que el usuario otorga permisos.
// Intercambia el código de autorización por un token de acceso y recupera la información del usuario.
func (g *GitHub) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
		return
	}

	token, err := g.Config.Exchange(context.Background(), code, g.pkceChallenge())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !token.Valid() {
		http.Error(w, "Invalid token", http.StatusInternalServerError)
		return
	}

	body, err := g.GetUserInfo(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

// GetUserInfo obtiene la información del usuario autenticado en GitHub.
//
// Devuelve el cuerpo de la respuesta como un slice de bytes. Si ocurre un error,
// devuelve un error.
func (g *GitHub) GetUserInfo(token *oauth2.Token) ([]byte, error) {
	client := g.Config.Client(context.Background(), token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("error creating request to GitHub API: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
