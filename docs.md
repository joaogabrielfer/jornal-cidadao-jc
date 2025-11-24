# Documentação da API e Páginas - Jornal Cidadão JC

Esta documentação descreve todos os endpoints, tanto as páginas web renderizadas quanto a API REST, disponíveis no projeto.

**URL Base**: `http://localhost:8000`

---

## 1. Páginas Web (Públicas)

Estes endpoints servem páginas HTML diretamente ao navegador do usuário.

| Método | Endpoint      | Descrição                                                                                             | Template Renderizado       |
| :----- | :------------ | :---------------------------------------------------------------------------------------------------- | :------------------------- |
| `GET`  | `/`           | Renderiza a página principal da aplicação.                                                            | `index.tmpl`               |
| `GET`  | `/login`      | Renderiza a página de login. Pode exibir uma mensagem de sucesso se o parâmetro `?signup=success` existir. | `login.tmpl`               |
| `GET`  | `/cadastro`   | Renderiza o formulário de cadastro de novos usuários.                                                 | `cadastro.tmpl`            |
| `GET`  | `/ultimas`    | Renderiza a página "Jornal Cidadão", o dashboard para o usuário visualizar seus posts enviados.       | `jc.tmpl`                  |
| `GET`  | `/charge/:id` | Exibe uma charge específica baseada no seu `id`.                                                      | `vizualizar_charge.tmpl`   |
| `GET`  | `/charge`     | Redireciona para a charge com o maior ID (a mais recente).                                            | Redireciona para `/charge/:id` |
| `GET`  | `/charges`    | Redireciona para a charge com o maior ID (a mais recente). Rota para o link do menu.                  | Redireciona para `/charge/:id` |

---

## 2. Páginas Web (Administrativas)

Acesso à área de gerenciamento do site. Estes endpoints servem páginas HTML com formulários e listas para administração.

| Método | Endpoint                  | Descrição                                                                               | Template Renderizado     |
| :----- | :------------------------ | :-------------------------------------------------------------------------------------- | :----------------------- |
| `GET`  | `/admin`                  | Renderiza a página principal do painel de admin.                                        | `admin.tmpl`             |
| `GET`  | `/admin/users`            | Renderiza a página de gerenciamento de usuários.                                        | `admin_users.tmpl`       |
| `GET`  | `/admin/adicionar-charge` | Renderiza o formulário para enviar uma nova charge.                                     | `adicionar_charge.tmpl`  |
| `GET`  | `/admin/charges`          | Renderiza a página para gerenciar a exclusão de charges.                                | `deletar_charge.tmpl`    |
| `GET`  | `/admin/materia`          | Renderiza o formulário para criar uma nova matéria (artigo).                            | `escrever_materia.tmpl`  |
| `GET`  | `/admin/materias`         | Renderiza a página para listar e gerenciar matérias existentes.                         | `materias.tmpl`          |
| `GET`  | `/admin/materias/:id/edit`| Renderiza o formulário para editar uma matéria existente.                                 | `atualizar_materia.tmpl`   |

---

## 3. API - Endpoints Públicos

Estes endpoints fornecem dados em formato JSON e são consumidos principalmente pelo frontend (JavaScript).

### Endpoints de Usuários

#### Criar um Novo Usuário
- **Endpoint**: `POST /api/users`
- **Descrição**: Cria um novo usuário. Espera dados de um formulário (`x-www-form-urlencoded`).
- **Respostas**:
    - **`303 See Other`**: Sucesso. Redireciona para `/login?signup=success`.
    - **`400 Bad Request`**: Erro de validação (ex: email já em uso).

### Endpoints de Charges

#### Listar Todas as Charges
- **Endpoint**: `GET /api/charges`
- **Descrição**: Retorna uma lista de todas as charges, ordenadas da mais recente para a mais antiga.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um array de objetos de charge.

#### Obter uma Charge Aleatória
- **Endpoint**: `GET /api/charges/random`
- **Descrição**: Retorna os dados de uma única charge selecionada aleatoriamente.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "charge": { ... } }`.
    - **`404 Not Found`**: Nenhuma charge encontrada.

#### Obter a Charge do Dia
- **Endpoint**: `GET /api/charge-do-dia`
- **Descrição**: Retorna a charge mais recente, considerada a "Charge do Dia".
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um único objeto de charge.
    - **`404 Not Found`**: Nenhuma charge encontrada.

### Endpoints de Matérias (Artigos)

#### Obter Todas as Matérias
- **Endpoint**: `GET /api/materias`
- **Descrição**: Retorna uma lista de todas as matérias (artigos editoriais), incluindo suas enquetes.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um array de objetos `Article`.

#### Obter uma Matéria por ID
- **Endpoint**: `GET /api/materia/:id`
- **Descrição**: Retorna uma matéria específica pelo seu ID, incluindo sua enquete.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um objeto `Article`.

### Endpoints do Jornal Cidadão (Posts)

#### Enviar um Novo Post
- **Endpoint**: `POST /api/post`
- **Descrição**: Permite que um usuário envie uma notícia (post). Espera dados `multipart/form-data`.
- **Campos do Formulário**:
    - `title` (string): Título do post.
    - `description` (string): Corpo/descrição do post.
    - `media_file` (file): Arquivo de imagem ou vídeo.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "status": "sucess", "message": "Notícia enviada para a moderação." }`.
    - **`400 Bad Request`**: Erro de validação.

#### Listar Posts Aprovados (com Paginação)
- **Endpoint**: `GET /api/posts-aprovados`
- **Descrição**: Retorna uma lista paginada de todos os posts que foram aprovados pela moderação.
- **Parâmetros de Query**:
    - `page` (int, opcional, default: 1): O número da página.
    - `limit` (int, opcional, default: 10): O número de itens por página.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um objeto `PaginatedPosts` com `posts` e `metadata`.

#### Obter um Post por ID
- **Endpoint**: `GET /api/post/:id`
- **Descrição**: Retorna os dados de um post específico pelo seu ID.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um objeto `Post`.

#### Obter Todos os Posts de um Autor
- **Endpoint**: `GET /api/status-noticias/:id`
- **Descrição**: Retorna todos os posts (com qualquer status) de um autor específico, baseado no ID do autor.
- **Parâmetros da URL**: `:id` (int) - O ID do **autor** (`users.id`).
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um array de objetos `Post`.
    - **`404 Not Found`**: Usuário não encontrado ou não possui posts.

#### Denunciar um Post
- **Endpoint**: `POST /api/post/:id/report`
- **Descrição**: Permite que um usuário denuncie um post por conteúdo impróprio ou outras razões.
- **Parâmetros da URL**: `:id` (int) - O ID do post a ser denunciado.
- **Corpo da Requisição (JSON)**:
    ```json
    {
      "reason": "Motivo da denúncia..."
    }
    ```
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "message": "Denúncia enviada com sucesso. Obrigado pela colaboração." }`.
    - **`400 Bad Request`**: Motivo não fornecido ou ID inválido.
    - **`404 Not Found`**: Post não encontrado.
    - **`500 Internal Server Error`**: Erro ao registrar no banco de dados.

### Endpoints de Enquetes (Polls)

#### Registrar um Voto
- **Endpoint**: `POST /api/enquete/votar/:id`
- **Descrição**: Incrementa a contagem de votos para uma opção de enquete específica.
- **Parâmetros da URL**: `:id` (int) - O ID da **opção da enquete** (`poll_options.id`).
- **Respostas**:
    - **`200 OK`**: Sucesso.
    - **`500 Internal Server Error`**: Erro ao registrar o voto.

---

## 4. API - Endpoints Administrativos

Endpoints para gerenciamento do sistema, acessíveis através das páginas de administração.

### Gerenciamento de Charges

#### Upload de Nova Charge
- **Endpoint**: `POST /admin/api/charge`
- **Descrição**: Envia uma nova charge. Espera dados `multipart/form-data`.
- **Respostas**:
    - **`200 OK`**: Sucesso. Renderiza `adicionar_charge.tmpl` com mensagem de sucesso.

#### Deletar uma Charge
- **Endpoint**: `DELETE /admin/api/charge/:id`
- **Descrição**: Deleta uma charge (registro no DB e arquivo físico).
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "message": "..." }`.

### Gerenciamento de Usuários

#### Listar Todos os Usuários (Admin)
- **Endpoint**: `GET /admin/api/users`
- **Descrição**: Retorna uma lista de todos os usuários.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um array de `User`.

#### Deletar um Usuário
- **Endpoint**: `DELETE /admin/api/user/:id`
- **Descrição**: Deleta um usuário do banco de dados.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "message": "..." }`.

### Gerenciamento de Matérias (Artigos)

#### Criar uma Matéria com Enquete
- **Endpoint**: `POST /admin/api/materia`
- **Descrição**: Cria uma nova matéria e, opcionalmente, uma enquete associada.
- **Respostas**:
    - **`200 OK`**: Sucesso. Renderiza `escrever_materia.tmpl` com mensagem de sucesso.

#### Atualizar uma Matéria com Enquete
- **Endpoint**: `PUT /admin/api/materia/:id`
- **Descrição**: Atualiza uma matéria e sua enquete.
- **Respostas**:
    - **`200 OK`**: Sucesso. Renderiza `atualizar_materia.tmpl` com mensagem de sucesso.

#### Deletar uma Matéria
- **Endpoint**: `DELETE /admin/api/materia/:id`
- **Descrição**: Deleta uma matéria do banco de dados.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "message": "..." }`.

### Gerenciamento de Posts (Jornal Cidadão)

#### Atualizar Status de um Post (Moderação)
- **Endpoint**: `PUT /admin/api/post-status/:id/:status`
- **Descrição**: Altera o status de um post enviado por um usuário (aprova, rejeita, etc.).
- **Parâmetros da URL**:
    - `:id` (int) - O ID do post a ser moderado.
    - `:status` (string) - O novo status. Valores possíveis: `aprovado`, `rejeitado`, `analise`.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "message": "Status do post atualizado com sucesso!" }`.
    - **`400 Bad Request`**: ID ou status inválido.
    - **`404 Not Found`**: Post não encontrado.
