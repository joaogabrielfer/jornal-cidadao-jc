# Documentação da API - Jornal Cidadão JC

Esta documentação descreve os endpoints da API e as páginas web servidas pelo backend do projeto.

**URL Base**: `http://localhost:8080`

---

## Páginas Web

Estes endpoints servem páginas HTML diretamente para o navegador e não são considerados parte da API REST.

### 1. Página Inicial

-   **Endpoint**: `/`
-   **Método**: `GET`
-   **Descrição**: Renderiza a página principal da aplicação (`index.tmpl`).

### 2. Página de Cadastro

-   **Endpoint**: `/cadastro`
-   **Método**: `GET`
-   **Descrição**: Renderiza o formulário de cadastro de novos usuários (`cadastro.tmpl`).

---

## API - Endpoints de Usuários

### 1. Criar um Novo Usuário

Cria um novo registro de usuário no sistema.

-   **Endpoint**: `/api/cadastro`
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
-   **`400 Bad Request`**: Erro de validação dos dados enviados.
    ```json
    // Exemplo 1: Senhas não conferem
    {
      "error": "A senha na confirmaçao esta diferente."
    }
    
    // Exemplo 2: Campos faltando
    {
      "error": "Todos os campos sao requeridos"
    }

    // Exemplo 3: Usuário ou email já existem
    {
      "error": "Nome ou email ja podem estar em uso"
    }
    ```
-   **`500 Internal Server Error`**: Erro inesperado no servidor (ex: erro ao gerar hash da senha).
    ```json
    {
      "error": "Falhou em criar conta"
    }
    ```

### 2. Listar Todos os Usuários

Retorna uma lista com o nome e e-mail de todos os usuários cadastrados. (Endpoint de utilidade/teste).

-   **Endpoint**: `/api/users`
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
    ```json
    {
        "error": "Erro processando lista de usuarios"
    }
    ```
---

## API - Endpoints de Charges

### 1. Listar Todas as Charges

Retorna uma lista de todas as charges disponíveis, ordenadas da **mais recente para a mais antiga**.

-   **Endpoint**: `/api/charges`
-   **Método**: `GET`

#### Respostas:

-   **`200 OK`**: Sucesso. Retorna um array de objetos de charge com nome e data formatada.
    ```json
    [
      {
        "filename": "charge_mais_nova.png",
        "date": "02-10-2025 18:30:00"
      },
      {
        "filename": "charge_antiga.png",
        "date": "01-10-2025 12:00:00"
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
    ```json
    {
      "error": "Não foi possível ler diretório de charges."
    }
    ```

### 2. Obter uma Charge Aleatória

Busca e retorna os dados de uma única charge selecionada aleatoriamente do diretório.

-   **Endpoint**: `/api/charges/random`
-   **Método**: `GET`

#### Respostas:

-   **`200 OK`**: Sucesso. Retorna um objeto JSON contendo os detalhes da charge sorteada.
    ```json
    {
      "charge": {
        "url": "/static/images/charges/charge_sorteada.png",
        "filename": "charge_sorteada.png",
        "title": "",
        "modtime": "2025-10-01T12:00:00Z"
      }
    }
    ```
    **Nota**: O campo `url` contém o caminho público que o frontend deve usar para exibir a imagem. O campo `modtime` está no formato ISO 8601 (UTC).

-   **`404 Not Found`**: Se nenhum arquivo de charge for encontrado no diretório.
    ```json
    {
      "message": "Nenhuma charge (arquivo) encontrada"
    }
    ```
-   **`500 Internal Server Error`**: Se ocorrer um erro ao ler o diretório ou as informações do arquivo.
    ```json
    {
      "error": "Não foi possível ler diretório de charges."
    }
    ```
