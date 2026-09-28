# Changelog

Todas as alterações relevantes deste template-base serão documentadas aqui.

Este projeto segue um formato próximo de Keep a Changelog.

## [Unreleased]

### Added

- Screenshots automáticos (`npm run screenshots`, Puppeteer) com dados de demonstração, temas claro/escuro e docs em `docs/SCREENSHOTS.md`.
- Scripts `lint` (ESLint + angular-eslint), `format`, `format:check` e `typecheck`, no padrão da Suplementos Store.
- Header com menu dropdown no mobile/tablet.
- Governança base para templates Angular da organização.
- Templates de issue e pull request.
- Validações de GitHub Actions e Dependabot.
- Política de contribuição e segurança.

### Changed

- Tema claro/escuro: cores de links e botões com contraste AA; cantos de cards e botões em `0.2rem`; medidas em `rem`/`%`.
- Formulários (bottom sheet) fecham com **Esc**; clique fora checado no overlay.

### Removed

- Badge de status do backend e `HealthService`.

## [0.1.0]

### Added

- Estrutura inicial do template-base.
