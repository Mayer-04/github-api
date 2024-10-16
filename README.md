# github api

## Paquete OAuth 2.0 en Go

El paquete que mencionas, golang.org/x/oauth2, sí es oficial. Es mantenido por el equipo de Go (Google) y se encuentra en el repositorio oficial del proyecto de Go en GitHub bajo golang/oauth2. Se considera una librería confiable y sólida para manejar OAuth 2.0 en aplicaciones Go.

¿Qué te ofrece este paquete?

Facilita el flujo de OAuth 2.0 (manejo de tokens, generación de URLs de autorización, etc.).
Soporte para varios servicios populares (Google, GitHub, etc.).
Manejo de refrescos de token, en caso de que el token expire.

## Hola

client_id = El id de cliente que ha recibido de github al registrarse (obligatorio).
redirect_uri = La URL en tu aplicación a donde se enviará a los usuarios después de la autorización.
login = Sugiere una cuenta específica para utilizar para registrarse y autorizar la app.
scope = Lista de ámbitos delimitados por espacio.
state = Una secuencia aleatoria indescifrable. Proteger de ataques de falsificación de solicitudes (recomendado).

- Si el usuario acepta la solicitud, GitHub le `redirecciona` de vuelta al sitio con un valor `code temporal` en un `parámetro de código` y con el estado que haya proporcionado en el paso anterior en un parámetro `state`.
- El código temporal caducará después de 10 minutos.
- Intercambie este valor code por un token de acceso:

```bash
POST https://github.com/login/oauth/access_token
```

- Este punto de conexión toma los siguientes parámetros de entrada:
- client_id (obligatorio)
- client_secret (obligatorio)
- code (obligatorio)
- redirect_uri: La URL en tu aplicación, hacia la cual se envía a los usuarios después de su autorización.

Predeterminadamente, la respuesta toma la siguiente forma:

```bash
access_token=gho_16C7e42F292c6912E7710c838347Ae178B4a&scope=repo%2Cgist&token_type=bearer
```

- El token de acceso te permite hacer solicitudes a la API a nombre de un usuario.

```bash
Authorization: Bearer OAUTH-TOKEN
GET https://api.github.com/user
```

Cada vez que recibas un token de acceso, debes usar el token para volver a validar la identidad del usuario.


## Otro

// https://github.com/login/oauth/authorize?
// client_id=Ov23liPm2wqrMvHShyCH
// &response_type=code
// &scope=repo+user+read%3Auser
// &state=state

/*
Si omites el RedirectURL en oauth2.Config, GitHub utilizará el valor predeterminado configurado cuando registraste
la aplicación en GitHub.
- El RedirectURL debe coincidir exactamente con el que has configurado en el panel de desarrollador de GitHub al registrar tu aplicación.
- GitHub no sabe dónde redirigir al usuario a menos que se lo especifiques en la configuración de OAuth.
- GitHub necesita saber de antemano cuál es la ruta de callback en tu aplicación, para que pueda redirigir al lugar
correcto después de la autenticación.
*/
