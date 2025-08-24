# Technology Stack

> Note: This document outlines the technologies we intend to use or grow into. Not all technologies listed are in use immediately.
> The stack is designed to be cloud-agnostic and heavily oriented towards open-source solutions.

## Development Stack

### Backend Stack

| Category | Technology | Description |
|----------|------------|-------------|
| **Language** | [![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/) | Backend programming language |
| **HTTP Router** | [![chi](https://img.shields.io/badge/chi-%2300ADD8?style=flat-square)](https://github.com/go-chi/chi) | Lightweight, idiomatic HTTP router |
| **Linting** | [![staticcheck](https://img.shields.io/badge/staticcheck-%2300ADD8?style=flat-square)](https://staticcheck.io/) | Go static analysis |
| ^^ | [![golangci-lint](https://img.shields.io/badge/golangci--lint-%2300ADD8.svg?style=flat-square&logo=go&logoColor=white)](https://golangci-lint.run/) | Go linters aggregator |
| **Testing** | [![testify](https://img.shields.io/badge/testify-%23C5D9E0?style=flat-square)](https://github.com/stretchr/testify) | Testing toolkit |
| ^^ | [![mockery](https://img.shields.io/badge/mockery-%23775AA5?style=flat-square)](https://github.com/vektra/mockery) | Mock code generator |
| ^^ | [![go-cmp](https://img.shields.io/badge/go--cmp-%2300ADD8?style=flat-square)](https://github.com/google/go-cmp) | Deep comparison library |
| **Database** | [![sqlc](https://img.shields.io/badge/sqlc-%23F7B93E?style=flat-square)](https://sqlc.dev/) | Type-safe SQL in Go |
| ^^ | [![golang-migrate](https://img.shields.io/badge/golang--migrate-%234169E1?style=flat-square)](https://github.com/golang-migrate/migrate) | Database migrations |
| ^^ | [![pgx](https://img.shields.io/badge/pgx-%23336791?style=flat-square)](https://github.com/jackc/pgx) | PostgreSQL driver |
| ^^ | [![go-redis](https://img.shields.io/badge/go--redis-%23DC382D?style=flat-square)](https://github.com/redis/go-redis) | Redis client |
| ^^ | [![mongo-driver](https://img.shields.io/badge/mongo--driver-%2347A248?style=flat-square)](https://github.com/mongodb/mongo-go-driver) | MongoDB driver |
| ^^ | [![gocql](https://img.shields.io/badge/gocql-%231287B1?style=flat-square)](https://github.com/gocql/gocql) | Cassandra driver |
| ^^ | [![elastic](https://img.shields.io/badge/elastic-%23005571?style=flat-square)](https://github.com/elastic/go-elasticsearch) | Elasticsearch client |
| **Messaging** | [![kafka-go](https://img.shields.io/badge/kafka--go-black?style=flat-square)](https://github.com/segmentio/kafka-go) | Kafka client |
| **Auth** | [![golang-jwt](https://img.shields.io/badge/golang--jwt-%23000000?style=flat-square)](https://github.com/golang-jwt/jwt) | JWT implementation |
| ^^ | [![oauth2](https://img.shields.io/badge/oauth2-%234285F4?style=flat-square)](https://github.com/golang/oauth2) | OAuth 2.0 client |
| **HTTP Client** | [![resty](https://img.shields.io/badge/resty-%23000000?style=flat-square)](https://github.com/go-resty/resty) | HTTP client |
| **Observability** | [![opentelemetry-go](https://img.shields.io/badge/opentelemetry-%23425CC7?style=flat-square)](https://github.com/open-telemetry/opentelemetry-go) | OpenTelemetry SDK |
| ^^ | [![zerolog](https://img.shields.io/badge/zerolog-%23000000?style=flat-square)](https://github.com/rs/zerolog) | Zero-allocation logger |
| **Utils** | [![decimal](https://img.shields.io/badge/decimal-%23000000?style=flat-square)](https://github.com/shopspring/decimal) | Arbitrary-precision fixed-point decimal |
| ^^ | [![sonyflake](https://img.shields.io/badge/sonyflake-%23000000?style=flat-square)](https://github.com/sony/sonyflake) | Distributed unique ID generator |

### Frontend Stack

> Coming soon...

## Protocol Stack

| Category | Technology | Description |
|----------|------------|-------------|
| **Gateway** | [![Kong](https://img.shields.io/badge/kong-%23003459.svg?style=for-the-badge&logo=kong&logoColor=white)](https://konghq.com/) | API gateway |
| **Protocol Format** | [![OpenAPI](https://img.shields.io/badge/openapi-%236BA539.svg?style=for-the-badge&logo=openapi-initiative&logoColor=white)](https://www.openapis.org/) | REST API specification |
| ^^ | [![GraphQL](https://img.shields.io/badge/-GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white)](https://graphql.org/) | Query language for APIs |
| ^^ | [![Protocol Buffers](https://img.shields.io/badge/protobuf-%23244C5A.svg?style=for-the-badge&logo=google&logoColor=white)](https://protobuf.dev/) | Interface definition language |
| **Transport** | [![gRPC](https://img.shields.io/badge/grpc-%23244C5A.svg?style=for-the-badge&logo=grpc&logoColor=white)](https://grpc.io/) | High-performance RPC framework |
| ^^ | [![WebSocket](https://img.shields.io/badge/websocket-%23010101.svg?style=for-the-badge&logo=socket.io&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) | Real-time communication |
| ^^ | [![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-000?style=for-the-badge&logo=apachekafka)](https://kafka.apache.org/) | Event streaming platform |

## Infrastructure Stack

| Category | Technology | Description |
|----------|------------|-------------|
| **Containerization** | [![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/) | Container runtime and build system |
| **Orchestration** | [![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io/) | Container orchestration |
| ^^ | [![K3s](https://img.shields.io/badge/k3s-%23FFC61C.svg?style=for-the-badge&logo=k3s&logoColor=black)](https://k3s.io/) | Lightweight Kubernetes for local development |
| **GraphQL Platform** | [![Hasura](https://img.shields.io/badge/hasura-%231EB4D4.svg?style=for-the-badge&logo=hasura&logoColor=white)](https://hasura.io/) | GraphQL engine (internal tooling) |

## Storage Stack

| Category | Technology | Description |
|----------|------------|-------------|
| **Primary DB** | [![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/) | Primary relational database |
| ^^ | [![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)](https://www.mongodb.com/) | Document database |
| ^^ | [![Cassandra](https://img.shields.io/badge/cassandra-%231287B1.svg?style=for-the-badge&logo=apache-cassandra&logoColor=white)](https://cassandra.apache.org/) | Column-oriented database |
| **Cache & Search** | [![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/) | In-memory cache & store |
| ^^ | [![Elasticsearch](https://img.shields.io/badge/elasticsearch-%23005571.svg?style=for-the-badge&logo=elasticsearch&logoColor=white)](https://www.elastic.co/) | Search engine |

## Observability Stack

| Category | Technology | Description |
|----------|------------|-------------|
| **Visualization** | [![Grafana](https://img.shields.io/badge/grafana-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/) | Visualization & dashboards |
| **Metrics** | [![Grafana Mimir](https://img.shields.io/badge/mimir-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/oss/mimir/) | Metrics storage |
| **Logs** | [![Grafana Loki](https://img.shields.io/badge/loki-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/oss/loki/) | Log aggregation |
| **Traces** | [![Grafana Tempo](https://img.shields.io/badge/tempo-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/oss/tempo/) | Distributed tracing |
| **Profiling** | [![Pyroscope](https://img.shields.io/badge/pyroscope-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/oss/pyroscope/) | Continuous profiling |
| **RUM** | [![Grafana Faro](https://img.shields.io/badge/faro-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com/oss/faro/) | Real user monitoring |
| **Testing** | [![k6](https://img.shields.io/badge/k6-%23425CC7.svg?style=for-the-badge&logo=k6&logoColor=white)](https://k6.io/) | Load testing |
| **Instrumentation** | [![OpenTelemetry](https://img.shields.io/badge/opentelemetry-%23425CC7.svg?style=for-the-badge&logo=opentelemetry&logoColor=white)](https://opentelemetry.io/) | Observability framework |

## Security Stack

| Category | Technology | Description |
|----------|------------|-------------|
| **Identity & Access** | [![Vault](https://img.shields.io/badge/vault-%23000000.svg?style=for-the-badge&logo=vault&logoColor=white)](https://www.vaultproject.io/) | Secrets management |
| ^^ | [![OPA](https://img.shields.io/badge/opa-%23231F20.svg?style=for-the-badge&logo=opa&logoColor=white)](https://www.openpolicyagent.org/) | Policy enforcement |
| **Authentication** | [![OAuth2](https://img.shields.io/badge/OAuth2-%234285F4.svg?style=for-the-badge&logo=google&logoColor=white)](https://oauth.net/2/) | Authentication protocol |
| ^^ | [![OIDC](https://img.shields.io/badge/OIDC-%23F78C40.svg?style=for-the-badge&logo=openid&logoColor=white)](https://openid.net/connect/) | Identity layer on OAuth2 |
| ^^ | [![JWT](https://img.shields.io/badge/JWT-%23000000.svg?style=for-the-badge&logo=json-web-tokens&logoColor=white)](https://jwt.io/) | Token format |
| **Network Security** | [![Cilium](https://img.shields.io/badge/cilium-%23F2F4F9.svg?style=for-the-badge&logo=cilium&logoColor=black)](https://cilium.io/) | Service mesh & networking |
| ^^ | [![TLS](https://img.shields.io/badge/TLS%201.3-black?style=for-the-badge&logo=let%27s-encrypt&logoColor=white)](https://www.rfc-editor.org/rfc/rfc8446) | Transport security |
| ^^ | [![MTLS](https://img.shields.io/badge/mTLS-black?style=for-the-badge&logo=let%27s-encrypt&logoColor=white)](https://www.rfc-editor.org/rfc/rfc8446) | Mutual TLS authentication |

## Development Tooling & SDLC

| Category | Technology/Practice | Description |
|----------|-------------------|-------------|
| **Version Control** | [![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)](https://github.com/) | Code hosting & collaboration |
| ^^ | Trunk Based Development | Single main branch with short-lived feature branches |
| **Code Review** | Pull Request | Required review and checks before merge |
| **Testing** | TDD | Test-driven development approach |
| **CI/CD** | Continuous Integration | Automated testing on every commit |
| ^^ | Continuous Delivery | Automated deployment pipelines |
| ^^ | GitOps | Git as single source of truth for deployments |
| ^^ | [![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)](https://github.com/features/actions) | CI/CD pipeline |
| **Infrastructure** | [![Terraform](https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white)](https://www.terraform.io/) | Infrastructure provisioning |
| ^^ | [![Ansible](https://img.shields.io/badge/ansible-%231A1918.svg?style=for-the-badge&logo=ansible&logoColor=white)](https://www.ansible.com/) | Configuration management |
| ^^ | [![Packer](https://img.shields.io/badge/packer-%23E7EEF0.svg?style=for-the-badge&logo=packer&logoColor=%2302A8EF)](https://www.packer.io/) | Machine image automation |
| **Local Dev** | [![Docker Compose](https://img.shields.io/badge/docker--compose-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://docs.docker.com/compose/) | Container orchestration for development |
| **Documentation** | [![Markdown](https://img.shields.io/badge/markdown-%23000000.svg?style=for-the-badge&logo=markdown&logoColor=white)](https://www.markdownguide.org/) | Documentation format |
| ^^ | [![Mermaid](https://img.shields.io/badge/mermaid-%23FF3670.svg?style=for-the-badge&logo=mermaid&logoColor=white)](https://mermaid.js.org/) | Diagrams as code |
| **K8s Tools** | [![Helm](https://img.shields.io/badge/helm-%23277A9F.svg?style=for-the-badge&logo=helm&logoColor=white)](https://helm.sh/) | Kubernetes package manager |
| ^^ | [![Kustomize](https://img.shields.io/badge/kustomize-%23326CE5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kustomize.io/) | Kubernetes configuration |

### Future Development Workflow Tools

- [![Tilt](https://img.shields.io/badge/tilt-%23142641.svg?style=for-the-badge&logo=tilt&logoColor=white)](https://tilt.dev/) - Local Kubernetes development
- [![Skaffold](https://img.shields.io/badge/skaffold-%23326CE5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://skaffold.dev/) - Kubernetes development workflow
- [![DevSpace](https://img.shields.io/badge/devspace-%23326CE5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://www.devspace.sh/) - Kubernetes development platform

---

## Notes

- Technologies listed here represent both current and planned tooling
- Stack is designed to be cloud-agnostic and heavily oriented towards open-source solutions
- Tools will be adopted gradually based on project needs and maturity
- Version numbers indicate current planned versions and may be updated
- Some tools may be replaced or supplemented as the project evolves

## Implementation Notes

- This tech stack represents the planned technologies and may evolve based on project requirements
- Some components are under evaluation for their necessity in the architecture
- Frontend technology choice is pending decision
- The stack is designed to be cloud-agnostic with emphasis on open-source solutions
- Technologies will be adopted gradually based on project needs and maturity
- All tools and technologies are selected with consideration for enterprise-grade requirements
