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

Para rodar os testes:

```bash
go test ./...
```

## Monitoramento e Observabilidade

O serviço expõe métricas no padrão do Prometheus em `GET /metrics`:

- **Volume de requisições** — `http_requests_total`, um contador rotulado por
  código HTTP (`code`) e método (`method`).
- **Disponibilidade do serviço** — derivada da métrica `up` que o Prometheus
  registra a cada coleta do alvo (`1` = disponível, `0` = indisponível).

O `HEALTHCHECK` do container consulta `GET /healthz`, um endpoint de liveness
fora da instrumentação — assim as sondas não inflam as métricas de requisições.
O Nginx expõe publicamente apenas `GET /projeto-korp`; `/metrics` e `/healthz`
ficam restritos à rede interna do Docker.

O `compose.yaml` inclui o container `otel-lgtm` (imagem `grafana/otel-lgtm`,
que empacota Grafana e Prometheus). O Prometheus coleta o `/metrics` do serviço
e o Grafana já vem com um dashboard provisionado para analisar o comportamento
do serviço.

```bash
# a rede é externa ao Compose; crie-a uma vez (idempotente)
docker network create devops-korp

# sobe o serviço, o proxy reverso e a stack de observabilidade
docker compose up -d --build
```

| Interface  | URL                           | Descrição                                  |
| ---------- | ----------------------------- | ------------------------------------------ |
| Serviço    | http://localhost/projeto-korp | Endpoint exposto via Nginx.                |
| Grafana    | http://localhost:3000         | Dashboard _Projeto Korp_ (acesso anônimo). |
| Prometheus | http://localhost:9090         | Alvos coletados e consultas ad-hoc.        |

## Provisionamento com Ansible (Parte 3)

O diretório [`ansible/`](ansible/) contém um playbook que provisiona todo o
ambiente (Partes 1 e 2) em um host **Ubuntu** com um único comando. Ele instala
o Docker, cria a rede, faz o build da imagem, sobe a stack com o Docker Compose
(que monta as configurações do Nginx e do monitoramento) e, ao final, valida o
serviço via HTTP e imprime a resposta no console.

### 1. Instale o Ansible (uma vez)

```bash
sudo apt update
sudo apt install -y ansible
```

### 2. Provisione o ambiente

```bash
cd ansible
ansible-playbook playbook.yml -K   # -K pergunta a senha do sudo (instalar Docker)
```

Por padrão o alvo é a própria máquina (conexão local, sem SSH). Para provisionar
um host Ubuntu remoto, clone este repositório nele e ajuste o `inventory.ini`.

| Etapa do playbook        | Como é atendida                                              |
| ------------------------ | ----------------------------------------------------------- |
| Instalação do Docker     | Repositório oficial via `apt` (engine + plugin do Compose). |
| Criação da rede Docker   | `devops-korp`, criada de forma idempotente.                 |
| Build + execução         | `docker compose up --detach --build`.                       |
| Nginx (proxy reverso)    | `compose.yaml` monta `nginx/` no container Nginx.           |
| Monitoramento (Grafana)  | `compose.yaml` monta `observability/` no `otel-lgtm`.       |
| Validação                | Requisição HTTP via proxy e exibição da resposta.           |

## Configuração

| Variável | Padrão | Descrição                          |
| -------- | ------ | ---------------------------------- |
| `PORT`   | `8080` | Porta em que o servidor escuta.    |
