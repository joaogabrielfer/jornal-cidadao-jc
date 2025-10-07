# Documentação da API - Jornal Cidadão JC (Extensiva)

## 1\. Introdução

Bem-vindo à documentação detalhada da API do Jornal Cidadão JC. Esta API fornece um conjunto de endpoints para gerenciar os recursos centrais da aplicação: **Usuários** e **Charges**.

A documentação está dividida entre a API RESTful pública, projetada para ser consumida por clientes (como um frontend JavaScript ou aplicativo móvel), e os endpoints de Administração, que são ações protegidas para gerenciamento de conteúdo.

  - **URL Base**: `http://localhost:8080`
  - **Formato de Dados**: Todas as respostas de dados da API são em formato `JSON`, seguindo a codificação `UTF-8`.
  - **Autenticação**: Atualmente, os endpoints públicos não requerem autenticação. Os endpoints de administração (prefixo `/admin`) são projetados para serem protegidos por um middleware de autenticação (não implementado nesta documentação) que deve verificar se o usuário é um administrador.

## 2\. Modelos de Dados

#### Objeto `User`

Representa um usuário no sistema.

| Campo      | Tipo   | Descrição                                                              |
| :--------- | :----- | :----------------------------------------------------------------------- |
| `id`       | int    | O identificador numérico único para o usuário.                           |
| `username` | string | O nome de usuário, que deve ser único no sistema.                        |
| `email`    | string | O endereço de e-mail do usuário, que também deve ser único.              |
| `status`   | string | O status atual da conta do usuário. Valores possíveis: `active`, `suspended`, `banned`. |

#### Objeto `Charge`

Representa uma charge e seus metadados associados.

| Campo      | Tipo   | Descrição                                                                     |
| :--------- | :----- | :------------------------------------------------------------------------------ |
| `id`       | int    | O identificador numérico único para a charge.                                 |
| `title`      | string | O título descritivo da charge, fornecido durante o upload.                      |
| `filename` | string | O nome do arquivo único como está salvo no servidor (ex: `1665097200-charge.png`). |
| `url`      | string | A URL pública completa para acessar a imagem (ex: `/static/images/charges/...`). |
| `date`     | string | A data e hora de criação da charge, formatada como `DD-MM-YYYY HH:MM:SS`.        |

-----

## 3\. Páginas Web (HTML)

Endpoints que servem interfaces de usuário renderizadas no servidor.

| Endpoint                  | Método | Descrição                                    |
| :------------------------ | :----- | :------------------------------------------- |
| `/`                       | `GET`  | Renderiza a página principal da aplicação.                |
| `/cadastro`               | `GET`  | Renderiza a página de cadastro de usuário.   |
| `/charge/:id`             | `GET`  | Renderiza a página de visualização de uma charge. |
| `/admin`                  | `GET`  | Renderiza a página principal do painel de admin. |
| `/admin/adicionar-charge` | `GET`  | Renderiza o formulário de upload de charge.      |
| `/admin/charges`          | `GET`  | Renderiza a página para listar e deletar charges. |
| `/admin/users`            | `GET`  | Renderiza a página para gerenciar usuários.      |

-----

## 4\. API Endpoints

### Grupo: Usuários

#### 4.1. Criar um Novo Usuário

Cria um novo registro de usuário. Este endpoint é tipicamente chamado por um formulário HTML tradicional.

  - **Endpoint**: `POST /api/users`
  - **Descrição**: Registra um novo usuário no banco de dados. Realiza validações para campos obrigatórios, senhas coincidentes e unicidade de `username` e `email`.
  - **Corpo da Requisição**: `application/x-www-form-urlencoded`

| Parâmetro        | Tipo   | Descrição                                           | Obrigatório |
| :--------------- | :----- | :---------------------------------------------------- | :---------- |
| `username`         | string | Nome de usuário desejado.                             | Sim         |
| `email`            | string | E-mail para a conta.                                  | Sim         |
| `password`         | string | Senha da conta.                                       | Sim         |
| `password-confirm` | string | Confirmação da senha. Deve ser idêntica a `password`. | Sim         |

##### Exemplo de Requisição (`cURL`):

```sh
curl -X POST http://localhost:8080/api/users \
-d "username=novo_usuario" \
-d "email=novo@exemplo.com" \
-d "password=senha123" \
-d "password-confirm=senha123"
```

##### Respostas:

  - **`303 See Other`**: Sucesso. O navegador é redirecionado para a página `/cadastro`.
  - **`400 Bad Request`**: Erro nos dados enviados. O corpo da resposta contém um objeto JSON com a descrição do erro.
    ```json
    // Exemplo 1: Senhas não conferem
    { "error": "A senha na confirmaçao esta diferente." }

    // Exemplo 2: Usuário ou e-mail já em uso
    { "error": "Nome ou email ja podem estar em uso" }
    ```
  - **`500 Internal Server Error`**: Falha interna no servidor ao processar o cadastro.
    ```json
    { "error": "Falhou em criar conta" }
    ```

#### 4.2. Listar Todos os Usuários

Retorna uma lista de todos os usuários cadastrados.

  - **Endpoint**: `GET /api/users`

##### Exemplo de Requisição (`cURL`):

```sh
curl -X GET http://localhost:8080/api/users
```

##### Respostas:

  - **`200 OK`**: Sucesso. Retorna um array de objetos `User`.
    ```json
    [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@exemplo.com",
        "status": "active"
      },
      {
        "id": 2,
        "username": "joao_silva",
        "email": "joao@exemplo.com",
        "status": "suspended"
      }
    ]
    ```

#### 4.3. Obter um Usuário Específico

Busca e retorna os dados de um único usuário pelo seu ID.

  - **Endpoint**: `GET /api/user/:id`

##### Parâmetros de URL:

| Parâmetro | Tipo | Descrição |
| :--- | :--- | :--- |
| `id` | int | O ID do usuário a ser buscado. |

##### Exemplo de Requisição (`cURL`):

```sh
curl -X GET http://localhost:8080/api/user/2
```

##### Respostas:

  - **`200 OK`**: Sucesso. Retorna um único objeto `User`.
    ```json
    {
      "id": 2,
      "username": "joao_silva",
      "email": "joao@exemplo.com",
      "status": "suspended"
    }
    ```
  - **`404 Not Found`**: O usuário com o ID especificado não foi encontrado.
    ```json
    { "error": "Usuário não encontrado" }
    ```

#### 4.4. Deletar um Usuário

Remove permanentemente um usuário do sistema.

  - **Endpoint**: `DELETE /api/user/:id`

##### Exemplo de Requisição (`cURL`):

```sh
curl -X DELETE http://localhost:8080/api/user/2
```

##### Respostas:

  - **`200 OK`**: Sucesso.
    ```json
    { "message": "Usuário 'joao_silva' deletado com sucesso!" }
    ```
  - **`404 Not Found`**: O usuário a ser deletado não foi encontrado.

-----

### Grupo: Charges

#### 4.5. Listar Todas as Charges

Retorna uma lista de todas as charges, ordenadas da mais recente para a mais antiga.

  - **Endpoint**: `GET /api/charges`

##### Exemplo de Requisição (`cURL`):

```sh
curl -X GET http://localhost:8080/api/charges
```

##### Respostas:

  - **`200 OK`**: Sucesso. Retorna um array de objetos `Charge`.
    ```json
    [
      {
        "id": 15,
        "url": "/static/images/charges/1665180000-charge-recente.png",
        "filename": "1665180000-charge-recente.png",
        "title": "A Inflação e o Mercado",
        "date": "07-10-2025 21:00:00"
      },
      {
        "id": 14,
        "url": "/static/images/charges/1665093600-charge-anterior.png",
        "filename": "1665093600-charge-anterior.png",
        "title": "Eleições e Tecnologia",
        "date": "06-10-2025 21:00:00"
      }
    ]
    ```
  - **`500 Internal Server Error`**: Se ocorrer uma falha ao consultar o banco de dados.

#### 4.6. Obter uma Charge Aleatória

Retorna os dados de uma única charge selecionada aleatoriamente, ideal para uma seção "Charge do Dia".

  - **Endpoint**: `GET /api/charges/random`

##### Exemplo de Requisição (`cURL`):

```sh
curl -X GET http://localhost:8080/api/charges/random
```

##### Respostas:

  - **`200 OK`**: Sucesso. Retorna um objeto JSON contendo um objeto `Charge`.
    ```json
    {
      "charge": {
        "id": 14,
        "url": "/static/images/charges/1665093600-charge-anterior.png",
        "filename": "1665093600-charge-anterior.png",
        "title": "Eleições e Tecnologia",
        "date": "06-10-2025 21:00:00"
      }
    }
    ```
  - **`404 Not Found`**: Nenhuma charge encontrada no banco de dados.

-----

### Grupo: Administração

Endpoints destinados a serem consumidos pelas páginas do painel de administração.

#### 4.7. Fazer Upload de Charge

  - **Endpoint**: `POST /admin/charge`
  - **Descrição**: Processa o upload de uma nova charge. Requer um corpo `multipart/form-data`.
  - **Corpo da Requisição**: `multipart/form-data`
    | Parâmetro | Tipo | Descrição |
    | :--- | :--- | :--- |
    | `title` | string | O título para a nova charge. |
    | `charge_file`| file | O arquivo de imagem (`.png`, `.jpg`, `.jpeg`, `.gif`). |

##### Exemplo de Requisição (`cURL`):

```sh
curl -X POST http://localhost:8080/admin/charge \
-F "title=Nova Charge de Teste" \
-F "charge_file=@/caminho/para/sua/imagem.png"
```

##### Respostas:

  - Renderiza a página `/admin/adicionar-charge` novamente com uma mensagem de sucesso ou erro.

#### 4.8. Deletar uma Charge

  - **Endpoint**: `DELETE /admin/charge/:id`
  - **Descrição**: Deleta uma charge do banco de dados e remove o arquivo correspondente do servidor.

##### Exemplo de Requisição (`cURL`):

```sh
curl -X DELETE http://localhost:8080/admin/charge/15
```

##### Respostas:

  - **`200 OK`**: Sucesso.
    ```json
    { "message": "Charge '1665180000-charge-recente.png' deletada com sucesso!" }
    ```
  - **`500 Internal Server Error`**: Se ocorrer um erro ao deletar do banco ou do sistema de arquivos.

#### 4.9. Atualizar Status de Usuário

  - **Endpoint**: `PATCH /admin/user/:id/status`
  - **Descrição**: Atualiza o status de um usuário (para `active`, `suspended` ou `banned`).
  - **Corpo da Requisição**: `application/json`
    ```json
    {
      "status": "suspended"
    }
    ```

##### Exemplo de Requisição (`cURL`):

```sh
curl -X PATCH http://localhost:8080/admin/user/2/status \
-H "Content-Type: application/json" \
-d '{"status": "suspended"}'
```

##### Respostas:

  - **`200 OK`**: Sucesso.
    ```json
    { "message": "Status do usuário atualizado para 'suspended' com sucesso." }
    ```
  - **`400 Bad Request`**: Se o `status` fornecido no corpo JSON for inválido.
