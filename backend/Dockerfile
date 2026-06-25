FROM golang:1.25-alpine AS builder

# Instala dependências de compilação
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copia os arquivos de dependências primeiro (otimiza o cache)
COPY go.mod ./
# Se o arquivo go.sum existir, copia ele também. Se não, ignora.
COPY go.sum* ./

# Baixa as dependências
RUN go mod download

# Copia o resto do código
COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Estágio final (imagem leve)
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env
# ESTA É A LINHA QUE FALTAVA:
COPY --from=builder /app/config ./config 

EXPOSE 8080
CMD ["./main"]
