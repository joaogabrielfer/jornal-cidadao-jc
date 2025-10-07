# Documentação da API - Jornal Cidadão JC (Revisado)

Esta documentação descreve os endpoints da API e as páginas web servidas pelo backend do projeto.

**URL Base**: `http://localhost:8080`

-----

## Páginas Web (HTML)

Estes endpoints servem páginas HTML diretamente para o navegador.

| Endpoint | Método | Descrição |
| :--- | :--- | :--- |
| `/` | `GET` | Renderiza a página principal da aplicação. |
| `/cadastro` | `GET` | Renderiza o formulário de cadastro de novos usuários. |
| `/charge/:id` | `GET` | Renderiza a página para visualizar uma charge específica pelo seu ID do banco de dados. |
| `/admin` | `GET` | Renderiza a página principal do painel de administração. |
| `/admin/adicionar-charge` | `GET` | Renderiza o formulário para upload de uma nova charge. |
| `/admin/charges` | `GET` | Renderiza a página para listar e deletar charges existentes. |
| `/admin/users` | `GET` | Renderiza a página para listar e gerenciar usuários. |

-----

## API

Endpoints que manipulam dados e retornam JSON.

### Usuários

#### 1\. Criar um Novo Usuário

  - **Endpoint**: `POST /api/users`
  - **Descrição**: Cria um novo usuário a partir de dados de um formulário.
  - **Respostas**: Redireciona para `/cadastro` em caso de sucesso ou retorna JSON de erro.

#### 2\. Listar Todos os Usuários

  - **Endpoint**: `GET /api/users`
  - **Descrição**: Retorna uma lista com todos os usuários e seus dados.
  - **Resposta `200 OK`**:
    ```json
    [
      {
        "id": 1,
        "username": "usuario1",
        "email": "usuario1@exemplo.com",
        "status": "active"
      }
    ]
    ```

#### 3\. Obter um Usuário Específico

  - **Endpoint**: `GET /api/user/:id`
  - **Descrição**: Retorna os dados de um único usuário pelo seu ID.
  - **Resposta `200 OK`**:
    ```json
    {
      "id": 1,
      "username": "usuario1",
      "email": "usuario1@exemplo.com",
      "status": "active"
    }
    ```
  - **Respostas de Erro**: `400 Bad Request` (ID inválido), `404 Not Found` (usuário não encontrado).

#### 4\. Deletar um Usuário

  - **Endpoint**: `DELETE /api/user/:id`
  - **Descrição**: Deleta permanentemente um usuário pelo seu ID.
  - **Resposta `200 OK`**:
    ```json
    {
      "message": "Usuário 'nome_do_usuario' deletado com sucesso!"
    }
    ```

### Charges

#### 1\. Listar Todas as Charges

  - **Endpoint**: `GET /api/charges`
  - **Descrição**: Retorna uma lista de todas as charges, ordenadas da mais recente para a mais antiga, com metadados vindos do banco de dados.
  - **Resposta `200 OK`**:
    ```json
    [
      {
        "id": 2,
        "url": "/static/images/charges/1665097320-charge_recente.png",
        "filename": "1665097320-charge_recente.png",
        "title": "Título da Charge Recente",
        "date": "07-10-2025 19:22:00"
      },
      {
        "id": 1,
        "url": "/static/images/charges/1665097200-charge_antiga.png",
        "filename": "1665097200-charge_antiga.png",
        "title": "Título da Charge Antiga",
        "date": "07-10-2025 19:20:00"
      }
    ]
    ```

#### 2\. Obter uma Charge Aleatória

  - **Endpoint**: `GET /api/charges/random`
  - **Descrição**: Retorna os dados de uma única charge selecionada aleatoriamente.
  - **Resposta `200 OK`**:
    ```json
    {
      "charge": {
        "id": 1,
        "url": "/static/images/charges/1665097200-charge_antiga.png",
        "filename": "1665097200-charge_antiga.png",
        "title": "Título da Charge Antiga",
        "date": "07-10-2025 19:20:00"
      }
    }
    ```

### Ações de Administração (via Web Pages)

As seguintes ações são executadas a partir das páginas de admin e retornam respostas para serem processadas pelo JavaScript no frontend.

#### 1\. Fazer Upload de Charge

  - **Endpoint**: `POST /admin/charge`
  - **Descrição**: Recebe um arquivo de imagem e um título de um formulário `multipart/form-data`. Salva o arquivo e insere seus metadados no banco.
  - **Respostas**: Renderiza a página de upload novamente com uma mensagem de sucesso ou erro.

#### 2\. Deletar uma Charge

  - **Endpoint**: `DELETE /admin/charge/:id`
  - **Descrição**: Deleta o registro de uma charge do banco de dados e remove o arquivo de imagem correspondente do servidor.
  - **Resposta `200 OK`**:
    ```json
    {
      "message": "Charge 'nome_do_arquivo.png' deletada com sucesso!"
    }
    ```
