package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	gogithub "github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github" // Configuración en específico para Github
)

/*
* - w http.ResponseWriter:
Es una interfaz que permite al servidor escribir una respuesta HTTP que será enviada de vuelta al cliente.
Escribir el cuerpo de la respuesta y definir el código de estado.

* - r *http.Request:
Es un puntero a la estructura http.Request, que representa la solicitud HTTP del cliente.
Esta estructura contiene información importante como:

1. La URL que solicitó el cliente (r.URL).
2. Los encabezados (r.Header).
3. Parámetros y Querys.

Se utiliza un puntero para evitar la copia de toda la estructura http.Request, lo cual podría ser ineficiente
en términos de memoria, ya que esta estructura puede contener mucha información
(como el cuerpo de la solicitud, encabezados, cookies, etc.).

* r.URL.Query().Get()
Obtener el valor de un parámetro de consulta (query parameter) en una URL.
Los parámetros de consulta son aquellos que vienen después del símbolo ? en una URL.
Por ejemplo: https://www.ejemplo.com/search?query=golang

r.URL.Query().Get() obtiene el valor de query ósea golang.
Si no existe, devolverá una cadena vacía.
*/

/*
* Diferentes autorizaciones y las más importantes que maneja el paquete OAuth 2.0
- amazon
- github
- google
- facebook
- gitlab
- heroku
- instagram
- linkedin
- microsoft
- paypal
- spotify
- twitch
- slack
- uber
- yahoo
*/

/*
1. Redirigir a los usuarios a Github para la autenticación.
2. Después implementar la devolución de llamada, osea, http://localhost:8080/callback
*/

func main() {
	// Configurando el registrador
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	mux := http.NewServeMux()

	mux.HandleFunc("/redirect", redirect)
	mux.HandleFunc("/callback", callback)

	logger.Info("start http", "address", ":8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("failed serving http", "error", err)
		os.Exit(1)
	}
}

// Redirigir al usuario a una URL específica
func redirect(w http.ResponseWriter, r *http.Request) {
	// Redirigir una solicitud HTTP desde una URL a otra.
	// Si un cliente solicita una página específica, y el servidor quiere redirigirlo a otra página o URL.
	// Estado de redirección temporal 307 (usar en vez del estado 302).
	// getRedirectURL(): Es la nueva URL a la que el cliente será redirigido.
	http.Redirect(w, r, getRedirectURL(), http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// TODO: Validar state
	// state := r.URL.Query().Get("state")

	config := getOAuthConfig()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Si todo está bien, esto nos dará el acceso que necesitamos
	// Convierte un código de autorización en un token.
	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("unable to get token: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repos, err := getCurrentUserRepos(token.AccessToken)
	if err != nil {
		log.Fatalf("unable to get repos: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Asegúrate de enviar el código de estado antes de escribir el cuerpo
	w.WriteHeader(http.StatusOK)

	// Escribimos los nombres de los repositorios en el cuerpo de la respuesta
	for _, r := range repos {
		w.Write([]byte(r.GetFullName() + ", "))
	}
}

func getOAuthConfig() *oauth2.Config {
	// configuración básica de OAuth
	return &oauth2.Config{
		// Estos dos campos deben estar en variables de entorno
		ClientID:     "Ov23liPm2wqrMvHShyCH",
		ClientSecret: "ec97c84ad62981e6ed3468444cd76c12461eea00",
		// Los alcances te permiten especificar exactamente el tipo de acceso que necesitas.
		// read:user -> Otorga acceso para leer los datos del perfil de un usuario.
		// user:follow	-> Otorga acceso para seguir o no seguir a otros usuarios.
		// delete_repo	-> Otorga acceso para eliminar repositorios que se pueden administrar.
		Scopes:   []string{"repo", "user", "read:user"},
		Endpoint: github.Endpoint,
	}
}

func getRedirectURL() string {
	config := getOAuthConfig()
	return config.AuthCodeURL("state")
}

// Servicio de repositorios, servicio de confirmaciones, servicio de solicitud de grupo, etc.
func getCurrentUserRepos(accessToken string) ([]*gogithub.Repository, error) {
	client := gogithub.NewClient(nil).WithAuthToken(accessToken)

	opt := &gogithub.RepositoryListByAuthenticatedUserOptions{
		Affiliation: "owner",
	}

	// Opciones de paginación
	// Puedes hacer un bucle aumentando la página
	opt.PerPage = 50

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return repos, nil
}

/*
- Struct Config = AuthCodeURL(), Exchange(), Client(), TokenSource()
- Struct Token = Valid()
- Constantes = AuthStyle
- PKCE = GenerateVerifier(), VerifierOption()
*/
