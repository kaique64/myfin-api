# 📄 Novo `README.md`

# Finix API – Gin + MongoDB REST API

Este projeto é uma API REST simples construída em **Golang**, usando o framework **Gin** e o banco de dados **MongoDB**.  
A estrutura segue boas práticas de organização, separando **configuração**, **conexão com banco**, **handlers**, **repositories** e **services**.

---

## 📦 Requisitos
- Go 1.18+ instalado
- MongoDB em execução localmente ou em um container  
  (URI padrão: `mongodb://localhost:27017`)

---

## 🚀 Como rodar o projeto

1. **Clonar o repositório**
   ```bash
   git clone https://github.com/your-username/finix-api.git
   cd finix-api/myapp
   ```

2. **Instalar dependências Go**

   ```bash
   go mod tidy
   ```

3. **Configurar variáveis de ambiente**
   Crie um arquivo `.env` na raiz do projeto (`myapp/`) com:

   ```env
   MONGODB_DATABASE_URL=mongodb://localhost:27017
   MONGODB_DATABASE=finixdb
   ```

4. **Rodar a aplicação**
   Execute o binário principal em `cmd/server`:

   ```bash
   go run cmd/server/main.go
   ```

5. **Testar os endpoints**

   * Health Check: [http://localhost:8080/health](http://localhost:8080/health)
   * Exemplo de items: [http://localhost:8080/items](http://localhost:8080/items)

---

## 📖 Estrutura do Projeto

```
myapp/
├── cmd/                 # Ponto de entrada da aplicação
│   └── server/          
│       └── main.go      # main principal que sobe o servidor
│
├── internal/            # Código interno
│   ├── config/          # Configurações (env, variáveis globais, etc.)
│   │   └── config.go
│   ├── db/              # Conexão com o banco de dados MongoDB
│   │   └── mongo.go
│   ├── models/          # Estruturas (structs) que representam collections
│   │   └── item.go
│   ├── repository/      # Operações de acesso ao banco (CRUD)
│   │   └── item_repository.go
│   ├── service/         # Regras de negócio
│   │   └── item_service.go
│   └── handler/         # Handlers/Controllers do Gin
│       └── item_handler.go
│
├── pkg/                 # Pacotes utilitários opcionais
│   └── logger/
│
├── go.mod
├── go.sum
└── README.md
```

---

## 🛠 Notas

* Banco de dados padrão: `finixdb`
* Collection padrão: `items`
* Você pode alterar as configs no arquivo `.env`.
