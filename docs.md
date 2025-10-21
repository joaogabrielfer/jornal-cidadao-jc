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
| `GET`  | `/admin/materias/:id/edit`| Renderiza o formulário para editar uma matéria existente.                                 | `atualizar_materia.tmpl`   |

---

## 3. API - Endpoints Públicos

Estes endpoints fornecem dados em formato JSON e são consumidos principalmente por um frontend (JavaScript) ou outros serviços.

### Endpoints de Usuários

#### Criar um Novo Usuário
- **Endpoint**: `POST /api/users`
- **Descrição**: Cria um novo usuário. Espera dados de um formulário (`x-www-form-urlencoded`).
- **Respostas**:
    - **`303 See Other`**: Sucesso. Redireciona para `/login?signup=success`.
    - **`400 Bad Request`**: Erro de validação.

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

### Endpoints de Matérias (Artigos)

#### Obter Todas as Matérias
- **Endpoint**: `GET /api/materias`
- **Descrição**: Retorna uma lista de todas as matérias disponíveis, incluindo suas enquetes.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um array de objetos `Article`.

#### Obter uma Matéria por ID
- **Endpoint**: `GET /api/materia/:id`
- **Descrição**: Retorna uma matéria específica pelo seu ID, incluindo sua enquete.
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna um objeto `Article`.

### Endpoints de Enquetes (Polls)

#### Registrar um Voto
- **Endpoint**: `POST /api/enquete/votar/:id`
- **Descrição**: Incrementa a contagem de votos para uma opção de enquete específica.
- **Parâmetros da URL**: `:id` (int) - O ID da **opção da enquete** (`poll_options.id`).
- **Respostas**:
    - **`200 OK`**: Sucesso.
    - **`404 Not Found`**: ID da opção não existe.
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
- **Descrição**: Retorna uma lista de todos os usuários (para a página de admin).
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
- **Descrição**: Cria uma nova matéria e, opcionalmente, uma enquete associada via formulário.
- **Respostas**:
    - **`200 OK`**: Sucesso. Renderiza `escrever_materia.tmpl` com mensagem de sucesso.

#### Atualizar uma Matéria com Enquete
- **Endpoint**: `PUT /admin/api/materia/:id`
- **Descrição**: Atualiza uma matéria e sua enquete via formulário.
- **Respostas**:
    - **`200 OK`**: Sucesso. Renderiza `atualizar_materia.tmpl` com mensagem de sucesso.

#### Deletar uma Matéria
- **Endpoint**: `DELETE /admin/api/materia/:id`
- **Descrição**: Deleta uma matéria do banco de dados (e sua enquete, via `ON DELETE CASCADE`).
- **Respostas**:
    - **`200 OK`**: Sucesso. Retorna `{ "message": "..." }`.
