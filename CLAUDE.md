# llm-wiki — Contexto del proyecto

## Qué es esto

Generador de workspaces para wikis de conocimiento mantenidos por IA. `llm-wiki` prepara la estructura, instala las instrucciones por backend y migra cambios del manifest; el agente ejecuta `/wiki-ingest`, `/wiki-query` y `/wiki-lint`.

## Cómo publicar una nueva release

```bash
make release VERSION=v0.3.0
```

Hoy `make release` valida formato semver, crea el tag y lo pushea. El empaquetado de releases vive en `.goreleaser.yaml`.

Estado actual a tener en cuenta:
- No hay un workflow de CI/preflight versionado en `.github/workflows/` dentro de este repo.
- Cualquier guardrail previo al release debe agregarse explícitamente en el repo; no asumirlo.
- El tap de Homebrew esperado sigue siendo `DavDaz/homebrew-llm-wiki`.

**El usuario actualiza con:**
```bash
brew update && brew upgrade llm-wiki
```

## Cómo instalar (usuario final)

```bash
brew tap DavDaz/llm-wiki
brew install llm-wiki
```

O con Go:

```bash
go install github.com/DavDaz/llm-wiki-generator/cmd/llm-wiki@latest
```

## Repositorios involucrados

| Repo | Propósito |
|------|-----------|
| `DavDaz/llm-wiki-generator` | Repo principal — código fuente |
| `DavDaz/homebrew-llm-wiki` | Tap de Homebrew — fórmula auto-generada por GoReleaser |

## Secrets requeridos en llm-wiki-generator

| Secret | Para qué |
|--------|----------|
| `HOMEBREW_TAP_TOKEN` | Fine-grained PAT con Contents R/W sobre `homebrew-llm-wiki` |

## Estructura del proyecto

```
cmd/llm-wiki/          → entrypoint del binario
internal/
  cmd/                 → comandos Cobra (init, manage, status, add-tool, remove-tool, migrate)
  generator/           → crea y migra wikis en el filesystem
  manifest/            → lee/escribe wiki.toml
  templates/           → render + archivos embebidos
  tools/               → registry de tool backends (claude-code, opencode, pi)
  tui/
    wizard/            → form TUI para llm-wiki init (huh + bubbletea)
    dashboard/         → panel de gestión para llm-wiki manage (bubbletea)
    styles/            → estilos Lipgloss compartidos
  version/             → versión inyectada por GoReleaser via ldflags
internal/templates/assets/
  GUIDE.md             → guía base generada para el wiki
  schema.md.template   → schema/instrucciones compartidas
  commands/            → prompts fuente para wiki-ingest/query/lint
.goreleaser.yaml       → config de build y distribución
```

## Comandos disponibles

```bash
llm-wiki                               # launcher fuera del wiki; dashboard dentro del wiki
llm-wiki init                          # wizard TUI para crear wiki nuevo
llm-wiki init --name X --slug x        # modo headless
llm-wiki manage                        # dashboard TUI para gestionar wiki
llm-wiki status                        # estado del wiki actual
llm-wiki add-tool opencode             # habilitar tool backend
llm-wiki remove-tool pi                # deshabilitar tool backend
llm-wiki migrate                       # sincronizar manifest con filesystem
llm-wiki version                       # imprimir versión
```

## Nota técnica importante — bubbletea + huh

Los valores del form TUI viven en un `*formValues` (puntero al heap), no como campos del `Model`. Esto es necesario porque Bubbletea pasa el Model por valor y los punteros que huh necesita para hacer binding quedan inválidos si están en el stack. No cambiar esto sin entender la razón.
