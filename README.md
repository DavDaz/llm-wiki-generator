# llm-wiki

CLI + TUI para crear y gestionar wikis de conocimiento mantenidos por IA, basado en el patrón LLM Wiki de Karpathy (abril 2026).

**Un comando crea tu wiki. La IA lo mantiene. El conocimiento se acumula.**

`llm-wiki` no es un motor RAG ni un parser determinístico de documentos. La CLI crea y mantiene el espacio de trabajo; el agente de IA ejecuta las operaciones de conocimiento: ingestar fuentes, responder preguntas y auditar la salud del wiki.

---

## Instalación

### Homebrew (recomendado)

```bash
brew tap DavDaz/llm-wiki
brew install llm-wiki
```

### Go install

```bash
go install github.com/DavDaz/llm-wiki-generator/cmd/llm-wiki@latest
```

---

## Camino principal

El flujo base es chico a propósito:

1. Creás el workspace con `llm-wiki init`.
2. Agregás fuentes curadas e inmutables en `raw/`.
3. Le pedís al agente elegido que ejecute `/wiki-ingest`.
4. Consultás el conocimiento compilado con `/wiki-query`.
5. Auditás periódicamente el wiki con `/wiki-lint`.

El resultado es un wiki markdown que mejora con el tiempo: `raw/` conserva las fuentes, `wiki/` acumula síntesis, `wiki/index.md` organiza la navegación y `wiki/log.md` deja trazabilidad.

---

## Crear un wiki nuevo

```bash
# TUI interactivo (recomendado)
llm-wiki init

# Modo headless con flags
llm-wiki init --name "MIDES RENAB" --slug mides-renab --tools claude-code,opencode
llm-wiki init --name "Legal Wiki" --slug legal-wiki --tools all --entities "usuario,rol,permiso"
```

El wizard TUI te guía por estos pasos:

1. **Nombre y slug** — identificador técnico del wiki (kebab-case)
2. **Idioma** — `es` o `en`
3. **Tools** — Claude Code, OpenCode, Pi (o todos)
4. **Entidades primarias** — los sustantivos centrales de tu dominio
5. **Tipos de página** — taxonomía del wiki (pre-cargado con defaults editables)
6. **Convenciones** — reglas de negocio que la IA debe aplicar siempre

Al terminar tenés un directorio listo con `CLAUDE.md` / `AGENTS.md` configurado para tu dominio, estructura `raw/` y `wiki/`, y los skills `/wiki-ingest`, `/wiki-query`, `/wiki-lint`.

---

## Qué es core y qué es opcional

| Nivel | Incluye | Rol en el patrón LLM Wiki |
| --- | --- | --- |
| Core | `raw/`, `wiki/`, `wiki.toml`, `index.md`, `log.md`, `/wiki-ingest`, `/wiki-query`, `/wiki-lint`, `CLAUDE.md` / `AGENTS.md` | Hace explícito el ciclo de fuentes curadas → síntesis markdown → consulta → auditoría. |
| Conveniencia | Dashboard TUI, gestión de tools, estados de páginas, creación rápida de notas en `raw/`, migraciones | Ayuda a operar el wiki sin cambiar el modelo mental principal. |
| Avanzado/opcional | Ollama Modelfile, múltiples backends, flujos locales especializados | Útil cuando el equipo lo necesita, pero no requerido para empezar ni para entender el patrón. |

---

## Gestionar un wiki existente

```bash
# Abre el dashboard TUI de gestión (tools + estados de páginas)
cd tu-wiki/
llm-wiki manage
# Si ejecutás sin argumentos dentro del wiki:
llm-wiki
```

`llm-wiki` dentro de un wiki y `llm-wiki manage` muestran un menú raíz con:
- New raw note
- Tools backends
- Drafts (`status: borrador`)
- Published (`status: vigente`)
- Deprecated (`status: deprecado`)

En cada bucket de páginas podés cambiar status con opciones contextuales: nunca aparece el status actual, siempre aparece `cancel`.

La opción `New raw note` sirve para capturar material manual en `raw/` sin salir del dashboard. Es una conveniencia: la nota queda como fuente curada y después el agente la incorpora al wiki con `/wiki-ingest`.

### Comandos headless

```bash
llm-wiki status                   # estado del wiki y tools instalados
llm-wiki add-tool opencode        # habilitar un tool backend
llm-wiki remove-tool pi           # deshabilitar un tool backend
llm-wiki migrate                  # aplicar cambios del manifest al filesystem
llm-wiki doctor                   # chequeo estructural read-only; devuelve exit code != 0 si hay errores
```

`llm-wiki doctor` hace un chequeo estructural read-only del wiki actual. Reporta problemas en el manifest, archivos core, salidas de tools y enlaces, y devuelve exit code no cero cuando encuentra errores estructurales, así que sirve para validaciones locales o CI sin mutar el workspace.

---

## Cómo funciona

### Primera vez — Crear wiki

```mermaid
flowchart TD
    A([Tú]) --> B["llm-wiki init"]

    B --> D{"Wizard TUI"}
    D --> D1["📝 Nombre del wiki"]
    D --> D2["🏷️ Entidades primarias"]
    D --> D3["📄 Tipos de página"]
    D --> D4["⚙️ Convenciones del dominio"]

    D1 & D2 & D3 & D4 --> E["Elige tools\n(Claude Code / OpenCode / Pi / Todos)"]

    E --> F["Wiki listo\ncon estructura completa"]

    F --> G["📂 raw/\nVacío — tus fuentes van aquí"]
    F --> H["📂 wiki/\nindex.md · log.md · sources.json"]
    F --> I["⚙️ commands/\nwiki-ingest · wiki-query · wiki-lint"]
    F --> J["🧠 CLAUDE.md / AGENTS.md\nCerebro del sistema"]

    style A fill:#4f46e5,color:#fff
    style F fill:#059669,color:#fff
    style J fill:#d97706,color:#fff
```

### Flujo continuo — Agregar conocimiento y consultar

```mermaid
flowchart TD
    A([Tú]) --> B["Copiás un documento\na raw/"]

    B --> C["/wiki-ingest"]

    subgraph IA ["La IA hace esto automáticamente"]
        C --> D["Lee el schema\n(reglas del dominio)"]
        D --> E["Lee wiki/log.md\n(y sources.json en modo dirigido)"]
        E --> F["Analiza el documento\ne identifica conceptos"]
        F --> G{"¿Existe una página\npara este concepto?"}
        G -- "Sí" --> H["Actualiza página\nexistente"]
        G -- "No" --> I["Crea página nueva\ncon frontmatter"]
        H & I --> J["Agrega wikilinks\nentre páginas relacionadas"]
        J --> K["Actualiza wiki/index.md"]
        K --> L["Registra en wiki/log.md"]
    end

    L --> M([Conocimiento compilado])

    M --> N["/wiki-query\n¿Qué permisos tiene\nel rol Supervisor?"]

    subgraph QUERY ["La IA responde así"]
        N --> O["Lee wiki/index.md"]
        O --> P["Abre solo las páginas\nrelevantes (máx. 5)"]
        P --> Q["Responde con citas\npagina-fuente"]
    end

    Q --> R([Respuesta con referencias])

    M --> S["/wiki-lint\n(periódico)"]
    S --> T["Reporte wiki/lint-YYYY-MM-DD.md\nErrores · Advertencias · Info"]

    style A fill:#4f46e5,color:#fff
    style M fill:#059669,color:#fff
    style R fill:#059669,color:#fff
```

---

## Agente local con Ollama (opcional)

Si usás **OpenCode o Pi con Ollama**, podés crear un modelo especializado que ya sabe exactamente cómo operar el wiki — sin tener que explicarle nada cada vez que lo abrís.

```bash
# Requiere Ollama con gpt-oss:20b descargado
ollama create gpt-oss-wiki-agent -f ollama/wiki-agent.modelfile
```

El modelo tiene el sistema de instrucciones del wiki integrado: sabe que si le hacés una pregunta sobre el dominio tiene que ejecutar `/wiki-query`, si le decís que hay un archivo en `raw/` ejecuta `/wiki-ingest`, si pedís una auditoría ejecuta `/wiki-lint`. Rechaza cualquier consulta fuera del scope del wiki.

**Verificar que quedó bien:**

```bash
ollama show gpt-oss-wiki-agent
# Capabilities: completion · tools · thinking ✓
```

**Usar con OpenCode o Pi:**

Apuntá el tool al modelo `gpt-oss-wiki-agent` en tu configuración de Ollama. Los comandos en `.opencode/commands/` y `.pi/prompts/` funcionan igual — el sistema prompt del Modelfile amplía lo que el agente ya sabe.

> Ver `ollama/wiki-agent.modelfile` para el sistema prompt completo y los parámetros.

---

## Flujo de trabajo diario

### Agregar conocimiento nuevo

```bash
cp mi-manual.pdf tu-wiki/raw/

# En Claude Code, OpenCode o Pi:
/wiki-ingest

# Reprocesar una fuente específica (sin releer todo raw/):
/wiki-ingest raw/manual-operativo.md
```

También podés crear una nota desde el dashboard:

```bash
cd tu-wiki/
llm-wiki manage
# Elegí "New raw note" y luego ejecutá /wiki-ingest con tu agente.
```

### Hacer preguntas

```bash
# En Claude Code, OpenCode o Pi:
/wiki-query ¿qué permisos tiene el rol Supervisor?
/wiki-query ¿cómo se registra un beneficiario nuevo?
```

### Auditar el wiki

```bash
# En Claude Code, OpenCode o Pi:
/wiki-lint
```

---

## Archivos clave generados

### `CLAUDE.md` / `AGENTS.md`

El schema del dominio — el archivo más importante. Define entidades, tipos de página, reglas de nomenclatura y convenciones. La IA lo lee antes de cualquier operación.

- `CLAUDE.md` si usás Claude Code
- `AGENTS.md` si usás OpenCode o Pi
- Ambos si usás múltiples tools

### `wiki/index.md`

Catálogo central. Una línea por página. La IA lo lee primero en cada query.

### `wiki/log.md`

Historial append-only de todas las operaciones. Sirve como auditoría de qué se procesó y cuándo.

### `wiki/sources.json`

Estado operativo de ingest por fuente (`fingerprint`, `processed_at`, `pages_touched`). Permite reprocesar un archivo puntual solo cuando cambió, sin barrer todo `raw/`.

---

## Estructura de una página wiki

```markdown
---
tipo: proceso
titulo: Crear Usuario
dominio: mides-renab
status: vigente
confianza: alta
fuentes: [raw/manual-usuarios-v2.pdf]
actualizado: 2026-04-21
---

# Crear Usuario

## Precondiciones

- El solicitante debe tener rol [[rol-administrador]]

## Pasos

1. ...

## Ver también

- [[asignar-rol]]
- [[politica-acceso]]
```

---

## Cuándo usar este patrón

✅ **Ideal para:**
- Documentación de sistemas internos (hasta ~200 artículos)
- Knowledge base de equipos pequeños (2-10 personas)
- Procesos, roles, permisos, manuales operativos

⚠️ **Considerá RAG si:**
- Tenés miles de documentos que cambian constantemente
- Necesitás búsqueda semántica sobre texto libre masivo

---

## Basado en

- [Andrej Karpathy — LLM Wiki (abril 2026)](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- [LLM Wiki v2 — lessons from production](https://gist.github.com/rohitg00/2067ab416f7bbe447c1977edaaa681e2)
- [agentskills.io open standard](https://agentskills.io)
