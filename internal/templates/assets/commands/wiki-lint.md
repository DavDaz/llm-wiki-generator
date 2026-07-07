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

### Paso 2 — Ejecutar los checks obligatorios

Antes de clasificar severidades, revisar SIEMPRE estas categorías semánticas:

1. **Contradicciones entre páginas o contra las fuentes** — detectar afirmaciones incompatibles entre páginas del wiki o contra la evidencia disponible en `raw/`.
2. **Claims desactualizados (stale claims)** — detectar páginas vigentes que quedaron atrasadas respecto de fuentes más recientes, cambios de proceso o notas del log.
3. **Consistencia de conceptos y entidades** — detectar conceptos/entidades nombrados o definidos de forma inconsistente, relaciones conflictivas o conceptos recurrentes sin página canónica.
4. **Wikilinks y enlaces críticos** — detectar `[[wikilinks]]` rotos, cross-links faltantes entre páginas estrechamente relacionadas y sucesores deprecados no enlazados.
5. **Citas y trazabilidad hacia `raw/`** — detectar afirmaciones importantes sin soporte suficiente en `fuentes`, citas débiles o trazabilidad incompleta hacia material base.
6. **Research gaps o vacíos de investigación** — detectar preguntas abiertas, cobertura insuficiente, ambigüedad no resuelta o áreas donde falta evidencia antes de subir la confianza.

Luego complementar con checks estructurales y operativos:

**Checks de Error 🔴** (bloquean confiabilidad del wiki):

7. **Frontmatter incompleto** — páginas sin algún campo obligatorio (tipo, titulo, dominio, status, confianza, fuentes, actualizado)
8. **Slugs inválidos** — nombres de archivo que no siguen las reglas de nomenclatura del schema
9. **Tipo inválido** — páginas con un `tipo` que no está definido en el schema

**Checks de Advertencia 🟡** (degradan calidad del wiki):

10. **Borradores viejos** — páginas con `status: borrador` creadas hace más de 30 días sin revisión
11. **Páginas largas** — páginas con más de 500 palabras que podrían dividirse
12. **Fuentes sin procesar** — archivos en `raw/` que no tienen entrada en `wiki/log.md`

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

### Contradicciones entre páginas o contra las fuentes
- `wiki/roles.md` afirma que Supervisor aprueba accesos, pero `wiki/politica-acceso.md` indica que solo Seguridad puede aprobarlos

### Citas y trazabilidad hacia `raw/`
- `wiki/crear-usuario.md` describe una excepción operativa sin fuente verificable en `raw/`

---

## 🟡 Advertencias (Y)

### Claims desactualizados (stale claims)
- `wiki/configuracion-smtp.md` no refleja el procedimiento documentado en `raw/manual-v3.pdf`

### Consistencia de conceptos y entidades
- `wiki/roles.md` usa "Supervisor" y `wiki/crear-usuario.md` usa "Aprobador" para la misma entidad sin aclaración

---

## 🔵 Info (Z)

### Wikilinks y enlaces críticos
- `wiki/onboarding.md` debería enlazar `[[crear-usuario]]` y `[[roles]]` pero no lo hace

### Research gaps o vacíos de investigación
- Falta evidencia en `raw/` para confirmar quién puede revocar permisos temporales

---

## Acciones recomendadas

1. [inspección concreta para validar la contradicción o claim desactualizado más crítico]
2. [seguimiento para completar citas, enlaces o evidencia faltante]
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

- ❌ Nunca modificar páginas durante el lint — solo reportar
- ❌ Nunca borrar páginas huérfanas automáticamente — solo reportar
- ❌ Nunca corregir automáticamente errores sin confirmación del usuario
- ❌ Nunca convertir el lint en ingest, query determinista, vector search, RAG o reparación automática
