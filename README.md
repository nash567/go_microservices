# go_microservices

## Folders
.
├── authentication-service
│   ├── authApp
│   ├── authentication-service.dockerfile
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── helper.go
│   │       ├── main.go
│   │       └── routes.go
│   ├── data
│   │   └── models.go
│   ├── go.mod
│   └── go.sum
├── broker-service
│   ├── brokerApp
│   ├── broker-service.dockerfile
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── helpers.go
│   │       ├── main.go
│   │       └── route.go
│   ├── event
│   │   ├── consumer.go
│   │   ├── emitter.go
│   │   └── event.go
│   ├── go.mod
│   ├── go.sum
│   └── logs
│       ├── logs_grpc.pb.go
│       ├── logs.pb.go
│       └── logs.proto
├── front-end
│   ├── cmd
│   │   └── web
│   │       ├── main.go
│   │       └── templates
│   │           ├── base.layout.gohtml
│   │           ├── footer.partial.gohtml
│   │           ├── header.partial.gohtml
│   │           └── test.page.gohtml
│   ├── frontEndApp
│   ├── front-end.dockerfile
│   └── go.mod
├── listener-service
│   ├── event
│   │   ├── consumer.go
│   │   └── event.go
│   ├── go.mod
│   ├── go.sum
│   ├── listenerApp
│   ├── listener-service.dockerfile
│   └── main.go
├── logger-service
│   ├── cmd
│   │   └── api
│   │       ├── grpc.go
│   │       ├── handler.go
│   │       ├── helper.go
│   │       ├── main.go
│   │       ├── routes.go
│   │       └── rpc.go
│   ├── data
│   │   └── model.go
│   ├── go.mod
│   ├── go.sum
│   ├── loggerServiceApp
│   ├── logger-service.dockerfile
│   └── logs
│       ├── logs_grpc.pb.go
│       ├── logs.pb.go
│       └── logs.proto
├── mail-service
│   ├── api
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── helpers.go
│   │       ├── mailer.go
│   │       ├── main.go
│   │       └── routes.go
│   ├── go.mod
│   ├── go.sum
│   ├── mailerApp
│   ├── mailer-service.dockerfile
│   └── templates
│       ├── mail.html.gohtml
│       └── mail.plain.gohtml
├── main.go
├── project
│   ├── caddy_config
│   ├── caddy_data
│   ├── caddy.dockerfile
│   ├── Caddyfile
│   ├── db-data
│   ├── docker-compose.yaml
│   ├── Makefile
│   └── swarm.yaml
└── README.md