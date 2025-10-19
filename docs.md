# Documentação da API e Páginas - Jornal Cidadão JC

Esta documentação descreve todos os endpoints, tanto as páginas web renderizadas quanto a API REST, disponíveis no projeto.

**URL Base**: `http://localhost:8080`

---

## 1. Páginas Web (Públicas)

Estes endpoints servem páginas HTML diretamente ao navegador do usuário.

| Método | Endpoint      | Descrição                                                                                             | Template Renderizado       |
| :----- | :------------ | :---------------------------------------------------------------------------------------------------- | :------------------------- |
| `GET`  | `/`           | Renderiza a página principal da aplicação.                                                            | `index.tmpl`               |
| `GET`  | `/login`      | Renderiza a página de login. Pode exibir uma mensagem de sucesso se o parâmetro `?signup=success` existir. | `login.tmpl`               |
| `GET`  | `/cadastro`   | Renderiza o formulário de cadastro de novos usuários.                                                 | `cadastro.tmpl`            |
| `GET`  | `/charge/:id` | Exibe uma charge específica baseada no seu `id`. A navegação é feita no frontend via API.             | `vizualizar_charge.tmpl`   |
| `GET`  | `/charge`     | Redireciona para a charge com o maior ID, caso a URL seja acessada sem um ID específico.               | Redireciona para `/charge/:id` |

---

## 2. Páginas Web (Administrativas)

Acesso à área de gerenciamento do site. Estes endpoints servem páginas HTML com formulários e listas para administração.

| Método | Endpoint                  | Descrição                                                                               | Template Renderizado     |
| :----- | :------------------------ | :-------------------------------------------------------------------------------------- | :----------------------- |
| `GET`  | `/admin`                  | Renderiza a página principal do painel de admin.                                        | `admin.tmpl`             |
| `GET`  | `/admin/users`            | Renderiza a página de gerenciamento de usuários.                                        | `admin_users.tmpl`       |
| `GET`  | `/admin/adicionar-charge` | Renderiza o formulário para enviar uma nova charge.                                     | `adicionar_charge.tmpl`  |
| `GET`  | `/admin/charges`          | Renderiza a página para gerenciar a exclusão de charges.                                | `deletar_charge.tmpl`    |
| `GET`  | `/admin/materia`          | Renderiza o formulário para criar uma nova matéria.                                     | `escrever_materia.tmpl`  |
| `GET`  | `/admin/materias`         | Renderiza a página para listar e gerenciar matérias existentes.                         | `materias.tmpl`          |

---

## 3. API - Endpoints Públicos

Estes endpoints fornecem dados em formato JSON e são consumidos principalmente por um frontend (JavaScript) ou outros serviços.

### Endpoints de Usuários

#### Criar um Novo Usuário
- **Endpoint**: `POST /api/users`
- **Descrição**: Cria um novo usuário. Espera dados de um formulário (`x-www-form-urlencoded`).
- **Parâmetros do Formulário**:
    - `username` (string, obrigatório)
    - `email` (string, obrigatório)
    - `password` (string, obrigatório)
    - `password-confirm` (string, obrigatório)
- **Respostas**:
    - **`303 See Other`**: Sucesso. Redireciona para `/login?signup=success`.
    - **`400 Bad Request`**: Erro de validação.
      ```json
      { "error": "Mensagem descritiva do erro" }
      ```

#### Listar Todos os Usuários
- **Endpoint**: `GET /api/users`
- **Descrição**: Retorna uma lista com todos os usuários cadastrados.
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      [ { "id": "1", "username": "usuario1", "email": "user1@jc.com" } ]
      ```

### Endpoints de Charges

#### Listar Todas as Charges
- **Endpoint**: `GET /api/charges`
- **Descrição**: Retorna uma lista de todas as charges, ordenadas da mais recente para a mais antiga. **Usado pelo frontend para a navegação de charges**.
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      [ { "id": 2, "url": "/static/...", "filename": "...", "title": "...", "date": "..." } ]
      ```

#### Obter uma Charge Aleatória
- **Endpoint**: `GET /api/charges/random`
- **Descrição**: Retorna os dados de uma única charge selecionada aleatoriamente.
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      { "charge": { "id": 1, "url": "/static/...", "filename": "...", "title": "...", "date": "..." } }
      ```
    - **`404 Not Found`**: Nenhuma charge encontrada.

### Endpoints de Matérias (Artigos)

#### Obter Todas as Matérias
- **Endpoint**: `GET /api/materias`
- **Descrição**: Retorna uma lista de todas as matérias disponíveis.
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      [ { "id": 1, "title": "...", "author": "...", "body": "..." } ]
      ```

#### Obter uma Matéria por ID
- **Endpoint**: `GET /api/materia/:id`
- **Descrição**: Retorna uma matéria específica pelo seu ID.
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      { "id": 1, "title": "...", "author": "...", "body": "..." }
      ```
    - **`400 Bad Request`**: ID inválido.
    - **`500 Internal Server Error`**: Erro ao obter a matéria.

---

## 4. API - Endpoints Administrativos

Endpoints para gerenciamento do sistema, acessíveis através das páginas de administração.

### Gerenciamento de Charges

#### Upload de Nova Charge
- **Endpoint**: `POST /admin/charge`
- **Descrição**: Envia uma nova charge. Espera dados `multipart/form-data`.
- **Parâmetros do Formulário**:
    - `title` (string, obrigatório)
    - `charge_file` (arquivo, obrigatório)
- **Respostas**:
    - **`200 OK`**: Sucesso. Renderiza `adicionar_charge.tmpl` com mensagem de sucesso.
    - **`400/500 Bad Request`**: Erro. Renderiza `error.tmpl`.

#### Deletar uma Charge
- **Endpoint**: `DELETE /admin/charge/:id`
- **Descrição**: Deleta uma charge (registro no DB e arquivo físico).
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      { "message": "Charge 'nome.png' deletada com sucesso!" }
      ```

### Gerenciamento de Usuários

#### Deletar um Usuário
- **Endpoint**: `DELETE /admin/user/:id`
- **Descrição**: Deleta um usuário do banco de dados.
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      { "message": "Usuário 'nome_usuario' deletada com sucesso!" }
      ```

### Gerenciamento de Matérias (Artigos)

#### Criar uma Matéria
- **Endpoint**: `POST /admin/materia`
- **Descrição**: Cria uma nova matéria via formulário.
- **Parâmetros do Formulário**:
    - `title`, `author`, `body` (strings, obrigatórios)
- **Respostas**:
    - **`200 OK`**: Sucesso. Renderiza `escrever_materia.tmpl` com mensagem de sucesso.
    - **`400/500 Bad Request`**: Erro. Renderiza `error.tmpl`.

#### Atualizar uma Matéria
- **Endpoint**: `PUT /admin/materia/:id`
- **Descrição**: Atualiza uma matéria existente via formulário.
- **Parâmetros da URL**: `:id` (int)
- **Parâmetros do Formulário**:
    - `title`, `author`, `body` (strings, obrigatórios)
- **Respostas**:
    - **`200 OK`**: Sucesso.
    - **`400/500 Bad Request`**: Erro. Renderiza `error.tmpl`.

#### Deletar uma Matéria
- **Endpoint**: `DELETE /admin/materia/:id`
- **Descrição**: Deleta uma matéria do banco de dados.
- **Respostas**:
    - **`200 OK`**: Sucesso.
      ```json
      { "message": "Matéria 'Título da Matéria' deletada com sucesso!" }
      ```
