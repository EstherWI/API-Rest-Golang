# ConexaoSolar user-crud
> API Rest desenvolvida em Golang que armazena as informações do usuário (Nome, Email, Senha) em um banco MongoDB (Atlas)


# Clone este repositório
```bash
$ git clone https://github.com/EstherWI/UserCRUD.git
```

# Acesse a pasta do projeto no terminal/cmd
```bash
$ cd ~/go-workspace/src (O projeto deve ficar em go-workspace/src/UserCRUD)
```

# Download dos pacotes:
```bash
go get go.mongodb.org/mongo-driver
go get github.com/gorilla/mux
go get github.com/gorilla/handlers
```
# No arquivo helper/helper.go
Especificar a URI do banco, o nome do banco e a collection (instruções nos comentários)

# Execute a aplicação em modo de desenvolvimento
```bash
go build main.go
go run main.go
```
# O servidor inciará na porta:8000 - acesse <http://localhost:8000/users> 
