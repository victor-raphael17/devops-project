# http-server-projeto-korp

Serviço HTTP em Go que expõe um único endpoint retornando o nome do projeto e o
horário atual em UTC, resolvido dinamicamente a cada requisição.

## Endpoint

`GET /projeto-korp` → `200 OK`

```json
{
  "nome": "Projeto Korp",
  "horario": "2026-06-22T14:30:00Z"
}
```

O campo `horario` está em UTC, no formato RFC 3339, e é recalculado a cada
requisição.

## Executando com Docker

O `Dockerfile` usa multi-stage build: um estágio compila o binário e o outro
apenas o executa em uma imagem mínima (`alpine`), como usuário não-root.

```bash
# build da imagem
docker build -t http-server-projeto-korp .

# execução do container
docker run --rm -p 8080:8080 http-server-projeto-korp
```

Em outro terminal:

```bash
curl http://localhost:8080/projeto-korp
```

## Executando localmente (sem Docker)

Requer Go 1.23+.

```bash
go run .
# ou
go build -o http-server-projeto-korp . && ./http-server-projeto-korp
```

## Configuração

| Variável | Padrão | Descrição                          |
| -------- | ------ | ---------------------------------- |
| `PORT`   | `8080` | Porta em que o servidor escuta.    |
