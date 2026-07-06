---
name: wiki-lint
description: Audita el wiki completo y genera un reporte de salud con errores, advertencias e info.
---

# wiki-lint

Audita el wiki completo y genera un reporte de salud con severidades.

## Cuándo usar

- Periódicamente (sugerido: cada vez que se hayan procesado 5+ fuentes nuevas)
- Cuando se sospeche inconsistencia
- Antes de compartir el wiki con nuevas personas
- Después de actualizar el schema del dominio

Activadores:
- "audita el wiki"
- "revisa la consistencia"
- "¿hay problemas en el wiki?"
- `/wiki-lint`

---

## Protocolo de ejecución

### Paso 0 — Cargar el schema
Leer el schema del dominio (`CLAUDE.md` o `AGENTS.md`, el que exista en el proyecto) completo. Las reglas del dominio definen qué está mal y qué no.

### Paso 1 — Inventario
Listar todos los archivos en `wiki/` (excepto `index.md`, `log.md`, y reportes de lint anteriores).
Leer `wiki/index.md` para comparar contra archivos reales.

### Paso 2 — Ejecutar los checks semánticos y estructurales

**Checks de Error 🔴** (bloquean confiabilidad del wiki):

1. **Frontmatter incompleto** — páginas sin algún campo obligatorio (tipo, titulo, dominio, status, confianza, fuentes, actualizado)
2. **Slugs inválidos** — nombres de archivo que no siguen las reglas de nomenclatura del schema
3. **Wikilinks rotos** — `[[referencias]]` a páginas que no existen en `wiki/`
4. **Deprecados sin sucesor** — páginas con `status: deprecado` sin campo `ver_sucesor`
5. **Tipo inválido** — páginas con un `tipo` que no está definido en el schema
6. **Contradicciones entre páginas o contra fuentes** — afirmaciones incompatibles entre páginas del wiki o frente a evidencia más reciente en `raw/`, `wiki/log.md` o `wiki/sources.json`

**Checks de Advertencia 🟡** (degradan calidad del wiki):

7. **Claims potencialmente desactualizados** — páginas que parecen stale frente a fechas, fuentes nuevas o cambios ya registrados en el log
8. **Entidades o conceptos inconsistentes** — nombres, definiciones o límites distintos para la misma entidad/concepto en páginas relacionadas
9. **Cobertura de enlaces insuficiente** — páginas huérfanas, islas temáticas o conceptos mencionados repetidamente sin `[[wikilinks]]` útiles
10. **Citas o trazabilidad de fuentes insuficientes** — afirmaciones relevantes sin respaldo claro en `fuentes`, `raw/` o el historial del wiki
11. **Research gaps o preguntas abiertas sin seguimiento** — dudas, supuestos o zonas de baja confianza que deberían derivar en revisión humana o nuevas fuentes

**Checks de Info 🔵** (oportunidades de mejora):

12. **Páginas largas** — páginas con más de 500 palabras que podrían dividirse
13. **Fuentes sin procesar** — archivos en `raw/` que no tienen entrada en `wiki/log.md`

### Paso 3 — Generar reporte
Guardar en `wiki/lint-YYYY-MM-DD.md`:

```markdown
---
tipo: reporte
titulo: Lint Report YYYY-MM-DD
dominio: [wiki-slug]
status: vigente
confianza: alta
fuentes: []
actualizado: YYYY-MM-DD
---

# Lint Report — YYYY-MM-DD

**Resumen:** X errores 🔴 · Y advertencias 🟡 · Z info 🔵

---

## 🔴 Errores (X)

### Frontmatter incompleto
- `wiki/nombre-pagina.md` — falta campo: `confianza`
- `wiki/otra-pagina.md` — falta campo: `fuentes`

### Wikilinks rotos
- `wiki/crear-usuario.md` → [[rol-supervisor]] (no existe)

### Contradicciones entre páginas o contra fuentes
- `wiki/politica-acceso.md` afirma revisión anual, pero `raw/manual-v4.pdf` indica revisión trimestral

---

## 🟡 Advertencias (Y)

### Claims potencialmente desactualizados
- `wiki/politica-acceso.md` — no refleja cambios registrados en `wiki/log.md` desde 2026-03-01

### Entidades o conceptos inconsistentes
- `wiki/roles.md` y `wiki/supervisor.md` describen responsabilidades incompatibles para "Supervisor"

---

## 🔵 Info (Z)

### Páginas largas
- `wiki/sistema-renab.md` — 823 palabras, considerar dividir

### Research gaps o preguntas abiertas sin seguimiento
- `wiki/sistema-renab.md` — marca integración externa como "pendiente validar" sin fuente ni tarea de seguimiento

---

## Acciones recomendadas

1. [inspección manual concreta para resolver el error más crítico]
2. [follow-up humano o nueva fuente para el segundo]
3. ...
```

### Paso 4 — Actualizar index.md
Agregar el reporte generado al índice.

### Paso 5 — Registrar en log.md
```markdown
## YYYY-MM-DD HH:MM — lint

**Resultado:** X errores, Y advertencias, Z info
**Reporte:** [[lint-YYYY-MM-DD]]
**Acción requerida:** sí / no

---
```

---

## Qué NO hacer

- ❌ Nunca modificar páginas automáticamente durante el lint — solo reportar
- ❌ Nunca borrar páginas huérfanas automáticamente — solo reportar
- ❌ Nunca corregir automáticamente errores sin confirmación del usuario
- ❌ Nunca convertir el lint en reparación automática, ingest determinístico o query/RAG automático
