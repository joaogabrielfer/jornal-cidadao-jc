# Guia de Contribuição - Projeto Jornal Cidadão JC

Olá! Agradecemos seu interesse em contribuir com o backend do Jornal Cidadão. Este guia irá ajudá-lo a configurar seu ambiente de desenvolvimento e a seguir nosso fluxo de trabalho.

## 🚀 Configurando o Ambiente

### Pré-requisitos

Antes de começar, certifique-se de que você tem as seguintes ferramentas instaladas:

-   [**Go**](https://golang.org/doc/install) (versão 1.21 ou superior)
-   [**Git**](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
-   [**Air**](https://github.com/cosmtrek/air) (Recomendado para live-reloading durante o desenvolvimento)

### Passos da Instalação

1.  **Clone o Repositório:**
    ```bash
    git clone https://github.com/seu-usuario/seu-repositorio.git
    cd seu-repositorio
    ```

2.  **Instale as Dependências:**
    O Go Modules cuidará de baixar e instalar todas as dependências do projeto.
    ```bash
    go mod tidy
    ```

3.  **Execute o Servidor:**
    Com o `air` instalado, você pode iniciar o servidor com live-reloading. O `air` irá automaticamente recompilar e reiniciar o servidor sempre que você salvar um arquivo `.go` ou `.tmpl`.
    ```bash
    air
    ```
    Se você não estiver usando o `air`, pode executar o projeto diretamente com:
    ```bash
    go run ./cmd/api/main.go
    ```
    O servidor estará disponível em `http://localhost:8000` (ou na porta definida em seu `.env`).

## 💻 Como Contribuir

Utilizamos o fluxo de trabalho com *Feature Branches* e *Pull Requests*. **Nunca faça commits diretamente na branch `main`**.

### Fluxo de Trabalho

1.  **Sincronize sua branch `main`:**
    ```bash
    git checkout main
    git pull origin main
    ```

2.  **Crie uma Nova Branch:**
    Crie uma branch a partir da `main` com um nome descritivo, prefixado com `feature/`, `fix/`, ou `refactor/`.
    ```bash
    # Exemplo para uma nova funcionalidade
    git checkout -b feature/implementar-autenticacao-jwt

    # Exemplo para uma correção de bug
    git checkout -b fix/corrigir-ordenacao-de-artigos
    ```

3.  **Desenvolva e Faça Commits:**
    Trabalhe na sua funcionalidade. Faça commits pequenos e atômicos com mensagens claras.
    ```bash
    git add .
    git commit -m "feat: Adiciona rota e handler para login de usuário"
    ```

4.  **Envie sua Branch para o GitHub:**
    ```bash
    git push -u origin feature/implementar-autenticacao-jwt
    ```

5.  **Abra um Pull Request (PR):**
    -   Vá para a página do seu repositório no GitHub.
    -   Clique no botão "Compare & pull request".
    -   Escreva um título claro e uma descrição do que foi feito. Se sua alteração resolve uma *Issue*, mencione-a (ex: `Resolves #25`).
    -   Aguarde a revisão do código.

6.  **Após o Merge:**
    Uma vez que seu PR for aprovado e mesclado, você pode deletar sua branch local e sincronizar sua `main` novamente.

## 🎨 Padrões de Código

-   **Formatação:** Sempre formate seu código Go antes de commitar. Use `gofmt -w .` ou configure seu editor para fazer isso automaticamente.
-   **Nomenclatura:** Siga as convenções do Go (`camelCase` para funções e variáveis, `PascalCase` para tipos e funções exportadas).
-   **Clareza:** Escreva um código claro e, se necessário, adicione comentários para explicar lógicas complexas.
