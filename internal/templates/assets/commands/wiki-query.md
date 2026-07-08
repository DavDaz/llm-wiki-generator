---
name: wiki-query
description: Responde preguntas usando exclusivamente el conocimiento compilado en wiki/.
argument-hint: [tu pregunta]
---

# wiki-query

Responde preguntas usando exclusivamente el conocimiento compilado en `wiki/`.

## Cuándo usar

Cuando se haga una pregunta sobre el dominio documentado en el wiki.

Activadores:
- Cualquier pregunta sobre el dominio
- "¿qué dice el wiki sobre...?"
- "¿cómo se hace...?"
- "¿qué permisos tiene...?"
- "¿cuáles son los...?"
- `/wiki-query [pregunta]`

---

## Protocolo de ejecución

### Paso 0 — Cargar el schema
Leer el schema del dominio (`CLAUDE.md` o `AGENTS.md`, el que exista en el proyecto) para entender el dominio antes de buscar.

### Paso 1 — Leer el índice
Leer `wiki/index.md` completo.
**Nunca responder desde memoria.** Si el índice no está accesible → informar al usuario.

### Paso 2 — Identificar páginas relevantes
De las páginas listadas en `index.md`, identificar cuáles son relevantes para la pregunta.
Máximo 5 páginas para mantener el contexto manejable.
Si hay duda entre varias páginas → priorizar las de `confianza: alta` y `status: vigente`.

### Paso 3 — Leer páginas seleccionadas
Abrir y leer solo las páginas identificadas como relevantes.
Nunca leer todo `wiki/` de forma indiscriminada.

### Paso 4 — Formular respuesta
```
[Respuesta clara y directa basada en el wiki]

Fuentes consultadas: [[pagina-1]], [[pagina-2]]
```

**Reglas al responder:**
- Si la información está en el wiki → responder con certeza y citar fuentes
- Si la información está parcialmente → responder lo que hay, indicar qué falta
- Si la información NO está en el wiki → decirlo explícitamente:
  ```
  Esta información no está documentada en el wiki todavía.
  Para incorporarla, agrega una fuente en `raw/` y ejecuta `/wiki-ingest raw/[archivo]`.
  ```
- Si hay páginas con `confianza: baja` o `status: borrador` en las fuentes → advertirlo:
  ```
  ⚠️ Nota: parte de esta respuesta proviene de páginas marcadas como borrador/baja confianza.
  Se recomienda verificación.
  ```

### Paso 5 — Ofrecer ingest si falta conocimiento
Si la pregunta revela conocimiento que no está documentado en `wiki/`, no crear páginas desde `/wiki-query`. Ofrecer iniciar un ingest separado con fuentes en `raw/`:

```
Esta información no está documentada en el wiki todavía.
Si deseas incorporarla, agrega la fuente en `raw/` y ejecuta `/wiki-ingest raw/[archivo]`.
```

---

## Qué NO hacer

- ❌ Nunca responder de memoria sin leer el wiki primero
- ❌ Nunca inventar información que no está en el wiki
- ❌ Nunca abrir más de 5 páginas por consulta
- ❌ Nunca omitir las fuentes consultadas en la respuesta
