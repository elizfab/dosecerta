# Frontend | Nome do Projeto

<p align="left">
  <img alt="Angular" src="https://img.shields.io/badge/Angular-FE9491?style=for-the-badge&logo=angular&logoColor=white">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-DABDAB?style=for-the-badge&logo=typescript&logoColor=333333">
  <img alt="PrimeNG" src="https://img.shields.io/badge/PrimeNG-F9B572?style=for-the-badge">
  <img alt="Jest" src="https://img.shields.io/badge/Jest-FFF1F2?style=for-the-badge&logo=jest&logoColor=333333">
</p>

## 📘 Sobre o frontend

Descreva aqui o objetivo do frontend.

Exemplo:

> Este frontend é responsável pela interface do usuário, navegação entre páginas, exibição de dados, modais de cadastro, tabelas, formulários, validações, integração com APIs e experiência visual do projeto.

## 🎯 Objetivo

- Criar uma interface clara, responsiva e acessível.
- Consumir APIs do backend.
- Organizar páginas por módulos.
- Reutilizar componentes globais.
- Padronizar botões, cards, modais, tabelas e formulários.
- Garantir testes e cobertura mínima.
- Manter documentação de instalação, execução e deploy.
- Padrões de contribuição e governança de branch.
- Templates de issue e pull request.
- Política de segurança e changelog.
- GitHub Actions reutilizáveis para validação do template.
- Configuração de Dependabot para dependências do GitHub Actions.

## 🧱 Stack técnica

| Tecnologia           | Uso                              |
| -------------------- | -------------------------------- |
| Angular              | Framework principal              |
| TypeScript           | Linguagem principal              |
| SCSS                 | Estilização                      |
| PrimeNG              | Componentes de UI                |
| PrimeIcons           | Ícones                           |
| NgRx                 | Estado global, quando necessário |
| Jest                 | Testes unitários                 |
| Coverage obrigatório | Validação                        |
| ESLint               | Padronização de código           |
| Prettier             | Formatação                       |
| GitHub Actions       | Pipeline de validação            |
| Dependabot           | x                                |

## 📁 Estrutura sugerida

```txt
frontend/
├── src/
│   ├── app/
│   │   ├── core/
│   │   ├── shared/
│   │   ├── features/
│   │   ├── layouts/
│   │   └── app.routes.ts
│   ├── assets/
│   ├── shared/
│   |   ├── shared/
|   │   │   ├── _tokens.scss
|   │   │   ├── _buttons.scss
|   │   │   ├── _cards.scss
|   │   │   └── styles.scss
│   └── environments/
├── angular.json
├── package.json
└── README.md
```

## ▶️ Como executar

```bash
npm install
npm start:dev
```

## 🔌 Integração com o backend

Este template já vem com `HttpClient`, `environments/` (dev/prod) e um `ApiService` genérico
configurados para consumir o [template de backend Go + MongoDB + Docker](../../template-backend-go-mongo-docker).
Veja **[INTEGRACAO-BACKEND.md](INTEGRACAO-BACKEND.md)** para o passo a passo completo: como as
portas se conectam, CORS, e como criar o serviço de um novo recurso.

## 🧪 Testes

```bash
npm run test
npm run test:coverage
```

## 🧭 Padrões de UI

### Botões

| Variante    | Uso                               |
| ----------- | --------------------------------- |
| `primary`   | Cadastrar, salvar, confirmar      |
| `secondary` | Filtrar, exportar, abrir detalhes |
| `tertiary`  | Voltar, cancelar, limpar          |
| `danger`    | Excluir, remover, apagar          |

### Componentes globais

- `app-page-header`
- `app-breadcrumb`
- `app-button`
- `app-card`
- `app-card-grid`
- `app-modal-form`
- `app-data-table`
- `app-empty-state`
- `app-confirm-dialog`

## Acessibilidade

- Todo botão com ícone deve ter texto ou `aria-label`.
- Todo input deve ter `label`.
- Não remover o foco visual do teclado.
- Usar contraste adequado entre texto e fundo.
- Usar `button` para ações e `a` para navegação.
- Evitar texto muito pequeno.
- Garantir navegação por teclado em modais.

## Governança esperada

- Branch principal: `main`.
- Alterações apenas via Pull Request.
- Revisão obrigatória antes do merge.
- Status checks obrigatórios antes do merge.
- Force push e delete de branch bloqueados.

## 📌 Checklist do frontend

- [ ] Criar estrutura inicial.
- [ ] Configurar rotas.
- [ ] Configurar layout principal.
- [ ] Criar header da página.
- [x] Integrar com backend.
- [ ] Criar testes unitários.
- [ ] Configurar coverage.
- [ ] Configurar pipeline.
- [ ] Publicar deploy.
