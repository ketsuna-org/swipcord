# click-ship-api
Projet GPE Etna, API V1 (Documentation incluse)

## Installation

### Prérequis

- [Devenv](https://devenv.sh)

### Installation

```bash
devenv update
```

## Utilisation

```bash
devenv up
```

## Documentation

Base URL: [http://localhost:4000](http://localhost:4000)

### Authentification

#### Enregistrer un utilisateur

```bash
curl -X POST http://localhost:4000/auth/register -d '{"email": "xxxx@test.fr", "password": "password", "username": "xxxx"}'
```

#### Se connecter

```bash
curl -X POST http://localhost:4000/auth/login -d '{"email": "xxxx@test.fr", "password": "password"}'
```
