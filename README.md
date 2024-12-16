
# Debit Authorizer Challenge

## Funcionalidades
- Enviar transacao

## Tecnologias Utilizadas
- Golang
- Testfy para testes de unidade

## Como rodar o projeto via API
Clone o projeto

```bash
git clone https://github.com/andreparelho/debit-authorizer.git
```

Navegue até o diretorio do projeto
```bash
cd api
```

Execute a aplicação
```bash
go run main.go
```

A aplicação estará disponível em http://localhost:8080

## Endpoints
- `POST /authorizer-debit:` enviar transacao; 

### Teste via postman, curl ou insomnia
```bash
curl -X POST http://localhost:8080/authorizer-debit \
     -H "Content-Type: application/json" \
     -d '{"clientId": "12345", "amount": 100.50}'
```

## Como rodar o projeto via CLI
Clone o projeto

Navegue até o diretorio do projeto
```bash
cd cli
```

Execute a aplicação
```bash
go run main.go
```

### Teste via cli
```bash
Client ID: "1"
Amount: 500
```