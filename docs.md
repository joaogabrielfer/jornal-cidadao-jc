# Documentação da API - Projeto Charges

Esta documentação descreve os endpoints disponíveis na API do projeto.

**URL Base**: `http://localhost:8080`

---

## Autenticação

Atualmente, os endpoints são públicos.

---

## Endpoints de Usuários

### 1. Criar um Novo Usuário

Cria um novo registro de usuário no sistema.

-   **Endpoint**: `POST /api/cadastro`
-   **Método**: `POST`
-   **Tipo de Conteúdo**: `application/x-www-form-urlencoded`

#### Parâmetros do Formulário:

| Parâmetro          | Tipo   | Descrição                                         | Obrigatório |
| ------------------ | ------ | --------------------------------------------------- | ----------- |
| `username`         | string | Nome de usuário único.                              | Sim         |
| `email`            | string | Endereço de e-mail único.                           | Sim         |
| `password`         | string | Senha do usuário.                                   | Sim         |
| `password-confirm` | string | Confirmação da senha. Deve ser igual a `password`.  | Sim         |

#### Respostas:

-   **`303 See Other`**: Sucesso. O usuário é redirecionado para a página `/cadastro`.
-   **`400 Bad Request`**:
    -   Se as senhas não conferem.
    -   Se algum campo obrigatório não foi preenchido.
    -   Se o `username` ou `email` já existem.
    ```json
    {
      "error": "Mensagem descritiva do erro"
    }
    ```
-   **`500 Internal Server Error`**: Erro inesperado no servidor.
    ```json
    {
      "error": "Falhou em criar a conta"
    }
    ```

### 2. Listar Todos os Usuários

Retorna uma lista com o nome e e-mail de todos os usuários cadastrados.

-   **Endpoint**: `GET /api/users`
-   **Método**: `GET`

#### Respostas:

-   **`200 OK`**: Sucesso. Retorna um array de objetos de usuário.
    ```json
    [
      {
        "username": "usuario1",
        "email": "usuario1@exemplo.com"
      },
      {
        "username": "usuario2",
        "email": "usuario2@exemplo.com"
      }
    ]
    ```
-   **`500 Internal Server Error`**: Erro ao buscar os dados no banco.

---

## Endpoints de Charges

### 1. Listar Todas as Charges

Retorna uma lista de todas as charges disponíveis, ordenadas da mais recente para a mais antiga.

-   **Endpoint**: `GET /api/charges`
-   **Método**: `GET`

#### Respostas:

-   **`200 OK`**: Sucesso. Retorna um array de objetos de charge.
    ```json
    [
      {
        "filename": "charge5.png",
        "date": "01-10-2025 18:30:00"
      },
      {
        "filename": "charge4.png",
        "date": "30-09-2025 12:00:00"
      }
    ]
    ```
-   **`404 Not Found`**: Se nenhuma imagem for encontrada no diretório.
    ```json
    {
      "message": "Nenhuma charge encontrada"
    }
    ```
-   **`500 Internal Server Error`**: Se ocorrer um erro ao ler o diretório de imagens.

#### Como usar as imagens no frontend:

O campo `filename` pode ser usado para construir a URL completa da imagem. O diretório de imagens é servido estaticamente.

**Exemplo**: Se o `filename` é `charge5.png`, a URL da imagem será:
`/static/images/charges/charge5.png`
