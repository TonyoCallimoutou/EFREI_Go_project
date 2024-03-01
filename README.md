# URL Shortener

Projet d'URL Shortener réalisé pour le cours de Go.
## Installation

Créer la base de données en important le fichier `url_shortener.sql`.

Créer un fichier `.env` dans la racine du projet, ce fichier va contenir les informations pour accéder à votre base de données.

`MYSQL_USER=[user]`

`MYSQL_PASSWORD=[password]`

`MYSQL_URL=[url]`

`MYSQL_PORT=[port]`

`MYSQL_DATABASE=ulrshortener`
    
## Lancer les tests

Pour lancer les tests, utiliser la commande suivante :

```bash
  go test
```


## Lancer le projet

Pour démarrer le projet, il faut lancer la commande `go run src/server.go` à la racine du projet. Ensuite vous pouvez accéder à l'application à l'adresse suivante : http://localhost:4000.


## Documentation API

#### Redirection

```http
  GET /{ShortUrl}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `ShortUrl`| `string` | **Requis**. La ShortUrl va rediriger vers l'URL associée. |

#### Créer une ShortUrl

```http
  POST /
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Url`     | `string` | **Requis**. URL que vous voulez raccourcir. |
| `ShortUrl`| `string` | **Requis**. L'URL raccourcie. |


#### Modifier une ShortUrl

```http
  PUT /
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Url`     | `string` | **Requis**. Nouvelle URL. |
| `ShortUrl`| `string` | **Requis**. L'URL raccourcie que vous voulez modifier. |

#### Supprimer une ShortUrl

```http
  DELETE /{ShortUrl}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `ShortUrl`| `string` | **Requis**. L'URL que vous voulez supprimer. |

#### Récupère toutes les URLSs dans la base de données.

```http
  GET /
```


## Description des dossiers

### Public

Le dossier `public` contient tous les fichiers utiles au frontend.

- Dossier Public
  - Dossier Script (contient les scripts)
  - Dossier Style (contient le css)
  - Fichier index.html (contient le html du front)

Le dossier `src` contient tous les fichiers utiles à l'API.

- Dossier src
  - Dossier domain (contient le type Shortener)
  - Dossier infrastructure (contient les fichiers permettant d'initialiser la BDD et le serveur)
  - Dossier interface
  - Dossier usecase
  - Fichier server.go (le fichier qu'il faut exécuter pour lancer l'API)
    
## Authors
- [@Tonyo Callimoutou](https://github.com/TonyoCallimoutou)
- [@Nathan Chevalet](https://github.com/NtchPlayer)
- [@Lucas Tourneaux](https://www.github.com/Xubeo)


