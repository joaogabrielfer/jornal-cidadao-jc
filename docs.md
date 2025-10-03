# Documentação da API - Jornal Cidadão JC (Atualizado)

Esta documentação descreve os endpoints da API e as páginas web servidas pelo backend do projeto.

**URL Base**: `http://localhost:8080`
**Versão da API**: `v1` (sugestão)

-----

## Páginas Web

Estes endpoints servem páginas HTML diretamente para o navegador.

### 1\. Página Inicial

  - **Endpoint**: `/`
  - **Método**: `GET`
  - **Descrição**: Renderiza a página principal da aplicação (`index.tmpl`).

### 2\. Página de Cadastro

  - **Endpoint**: `/cadastro`
  - **Método**: `GET`
  - **Descrição**: Renderiza o formulário de cadastro de novos usuários (`cadastro.tmpl`).

### 3\. Visualizar Charge Específica

  - **Endpoint**: `/charge/:id`
  - **Método**: `GET`
  - **Descrição**: Renderiza uma página HTML para visualizar uma charge específica, identificada pelo seu `ID`. O `ID` corresponde à ordem cronológica (ID 1 é a mais antiga).
  - **Parâmetro de URL**:
    | Parâmetro | Tipo    | Descrição                                 |
    | --------- | ------- | ----------------------------------------- |
    | `id`      | integer | O ID da charge que deve ser exibida. |

-----

## API v1

Endpoints que retornam dados no formato JSON.

### Usuários

#### 1\. Criar um Novo Usuário

Cria um novo registro de usuário no sistema.

  - **Endpoint**: `/api/users`
  - **Método**: `POST`
  - **Tipo de Conteúdo**: `application/x-www-form-urlencoded`

##### Parâmetros do Formulário:

| Parâmetro | Tipo | Descrição | Obrigatório |
| :--- | :--- | :--- | :--- |
| `username` | string | Nome de usuário único. | Sim |
| `email` | string | Endereço de e-mail único. | Sim |
| `password` | string | Senha do usuário. | Sim |
| `password-confirm` | string | Confirmação da senha. Deve ser igual a `password`. | Sim |

##### Respostas:

  - **`303 See Other`**: Sucesso. Redireciona para `/cadastro`.
  - **`400 Bad Request`**: Erro de validação ou usuário/email já existente.
    ```json
    {"error": "A senha na confirmaçao esta diferente."}
    // ou
    {"error": "Nome ou email ja podem estar em uso"}
    ```
  - **`500 Internal Server Error`**: Erro inesperado no servidor.

#### 2\. Listar Todos os Usuários

Retorna uma lista com os dados de todos os usuários cadastrados.

  - **Endpoint**: `/api/users`
  - **Método**: `GET`

##### Respostas:

  - **`200 OK`**: Sucesso. Retorna um array de objetos de usuário.
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
  - **`500 Internal Server Error`**: Erro ao buscar os dados no banco.

-----

### Charges

#### 1\. Listar Todas as Charges

Retorna uma lista de todas as charges disponíveis, **ordenadas da mais recente para a mais antiga**. O `ID` é atribuído de forma inversa (a charge mais antiga tem o `ID: 1`).

  - **Endpoint**: `/api/charges`
  - **Método**: `GET`

##### Respostas:

  - **`200 OK`**: Sucesso. Retorna um array de objetos de charge.
    ```json
    [
      {
        "id": 2,
        "url": "/static/images/charges/charge_recente.png",
        "filename": "charge_recente.png",
        "title": "",
        "date": "03-10-2025 00:02:00"
      },
      {
        "id": 1,
        "url": "/static/images/charges/charge_antiga.png",
        "filename": "charge_antiga.png",
        "title": "",
        "date": "02-10-2025 15:30:00"
      }
    ]
    ```
  - **`404 Not Found`**: Nenhuma charge encontrada.
  - **`500 Internal Server Error`**: Erro ao ler o diretório.

#### 2\. Obter uma Charge Aleatória

Retorna os dados de uma única charge selecionada aleatoriamente.

  - **Endpoint**: `/api/charges/random`
  - **Método**: `GET`

##### Respostas:

  - **`200 OK`**: Sucesso. Retorna um objeto JSON contendo os detalhes da charge sorteada.
    ```json
    {
      "charge": {
        "id": 1,
        "url": "/static/images/charges/charge_sorteada.png",
        "filename": "charge_sorteada.png",
        "title": "",
        "date": "02-10-2025 15:30:00"
      }
    }
    ```
  - **`404 Not Found`**: Nenhuma charge encontrada.
  - **`500 Internal Server Error`**: Erro ao ler o diretório.
