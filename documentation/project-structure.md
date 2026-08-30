llm-playground-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── health_handler.go
│   │   ├── generation_handler.go
│   │   └── model_handler.go
│   ├── service/
│   │   └── generation_service.go
│   ├── provider/
│   │   ├── provider.go
│   │   └── llm_provider.go
│   ├── model/
│   │   ├── request.go
│   │   ├── response.go
│   │   └── model_config.go
│   ├── middleware/
│   │   ├── request_id.go
│   │   ├── logging.go
│   │   ├── recovery.go
│   │   └── rate_limit.go
│   ├── pricing/
│   │   └── calculator.go
│   └── apperror/
│       └── error.go
├── tests/
│   └── integration/
├── .env.example
├── Dockerfile
├── go.mod
├── go.sum
└── README.md