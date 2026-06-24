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

## Monitoramento e Observabilidade

O serviço expõe métricas no padrão do Prometheus em `GET /metrics`:

- **Volume de requisições** — `http_requests_total`, um contador rotulado por
  código HTTP (`code`) e método (`method`).
- **Disponibilidade do serviço** — derivada da métrica `up` que o Prometheus
  registra a cada coleta do alvo (`1` = disponível, `0` = indisponível).

O `compose.yaml` inclui o container `otel-lgtm` (imagem `grafana/otel-lgtm`,
que empacota Grafana e Prometheus). O Prometheus coleta o `/metrics` do serviço
e o Grafana já vem com um dashboard provisionado para analisar o comportamento
do serviço.

```bash
# sobe o serviço, o proxy reverso e a stack de observabilidade
docker compose up -d --build
```

| Interface  | URL                     | Descrição                                  |
| ---------- | ----------------------- | ------------------------------------------ |
| Serviço    | http://localhost/projeto-korp | Endpoint exposto via Nginx.          |
| Grafana    | http://localhost:3000   | Dashboard _Projeto Korp_ (acesso anônimo). |
| Prometheus | http://localhost:9090   | Alvos coletados e consultas ad-hoc.        |

## Configuração

| Variável | Padrão | Descrição                          |
| -------- | ------ | ---------------------------------- |
| `PORT`   | `8080` | Porta em que o servidor escuta.    |
