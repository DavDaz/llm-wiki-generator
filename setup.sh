#!/bin/bash
# =============================================================================
# llm-wiki-template — setup.sh
# Genera un wiki nuevo a partir del template, configurado para tu dominio.
# Uso: ./setup.sh
# =============================================================================

set -e

# ─────────────────────────────────────────
# Verificaciones previas
# ─────────────────────────────────────────

if ! command -v python3 &>/dev/null; then
    echo "Error: python3 no encontrado. Es requerido para generar el schema del wiki."
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

for required_file in "schema.md.template" "commands/wiki-ingest.md" "commands/wiki-query.md" "commands/wiki-lint.md"; do
    if [[ ! -f "${SCRIPT_DIR}/${required_file}" ]]; then
        echo "Error: archivo requerido no encontrado: ${required_file}"
        exit 1
    fi
done

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║       LLM Wiki — Setup               ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"
echo ""

# ─────────────────────────────────────────
# 1. Recolectar datos del dominio
# ─────────────────────────────────────────

echo -e "${YELLOW}→ Nombre del wiki${NC}"
echo "  Nombre legible de tu dominio. Aparece en el título del wiki y en los logs."
echo "  Ej: MIDES RENAB, Banco XYZ, Sistema de Facturación"
read -r WIKI_NAME

echo ""
echo -e "${YELLOW}→ Slug del wiki${NC}"
echo "  Identificador técnico en kebab-case. Se usa como:"
echo "    - Nombre del directorio creado (ej: banco-xyz-wiki/)"
echo "    - Campo 'dominio:' en el frontmatter de cada página del wiki"
echo "  Solo letras minúsculas, números y guiones. Sin espacios ni acentos."
echo "  Ej: mides-renab, banco-xyz, sistema-facturacion"
read -r WIKI_SLUG
if [[ ! "$WIKI_SLUG" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
    echo "Error: slug inválido. Usá solo letras minúsculas, números y guiones (ej: mides-renab)."
    exit 1
fi

echo ""
echo -e "${YELLOW}→ Idioma principal${NC} (ej: es, en) [default: es]"
echo "  La IA generará todas las páginas del wiki en este idioma."
read -r LANGUAGE
LANGUAGE=${LANGUAGE:-es}

echo ""
echo -e "${YELLOW}→ Directorio destino${NC} [default: ./${WIKI_SLUG}-wiki]"
echo "  Ruta donde se creará el wiki. Puede ser relativa o absoluta."
echo "  Ej: ../mis-wikis/banco-xyz  o  /home/usuario/wikis/banco-xyz"
read -r WIKI_DIR
WIKI_DIR=${WIKI_DIR:-"./${WIKI_SLUG}-wiki"}

echo ""
echo -e "${YELLOW}→ CLI a utilizar${NC}"
echo "  1) Claude Code  (.claude/commands/ + CLAUDE.md)"
echo "  2) OpenCode     (.opencode/commands/ + AGENTS.md)"
echo "  3) Pi           (.pi/prompts/ + AGENTS.md)"
echo "  4) Todos"
read -r -p "  Opción [1/2/3/4] (default: 1): " CLI_CHOICE
CLI_CHOICE=${CLI_CHOICE:-1}
if [[ ! "$CLI_CHOICE" =~ ^[1234]$ ]]; then
    echo "Error: opción inválida. Ingresá 1, 2, 3 o 4."
    exit 1
fi

echo ""
echo -e "${YELLOW}→ Entidades primarias del dominio${NC}"
echo "  Son los conceptos centrales ('sustantivos') de tu dominio completo."
echo "  La IA los usa como anclas al procesar documentos: cuando encuentre"
echo "  mención a una de estas entidades, creará o actualizará su página en el wiki."
echo ""
echo "  Pensá en todos los objetos, actores y sistemas que existen en tu dominio,"
echo "  independientemente de los documentos que vayas a cargar."
echo "  Podés agregar más entidades en el archivo de instrucciones si el dominio crece."
echo ""
echo "  Ej: usuario, rol, permiso, beneficiario, expediente, sistema-renab"
echo "  Escribí una por línea. Línea vacía para terminar."
echo ""

ENTITIES=()
while true; do
    read -r -p "  Entidad: " entity
    [[ -z "$entity" ]] && break
    ENTITIES+=("$entity")
done

echo ""
echo -e "${YELLOW}→ Tipos de página del dominio${NC}"
echo "  Define qué categorías de conocimiento existen en tu wiki."
echo "  Cada página generada por la IA tendrá exactamente uno de estos tipos,"
echo "  que determina su estructura y cómo se nombra el archivo."
echo ""
echo "  Tipos disponibles y cuándo usarlos:"
echo "    proceso    → pasos para hacer algo (ej: crear-usuario.md, anular-entrada.md)"
echo "    referencia → listas, tablas, definiciones (ej: roles.md, codigos-error.md)"
echo "    entidad    → descripción de un sistema o actor (ej: sistema-renab.md)"
echo "    politica   → reglas o restricciones que deben cumplirse"
echo "    regulacion → normativa legal con cita de fuente"
echo "    reporte    → generado automáticamente por /wiki-lint"
echo ""
echo "  Dejá vacío para usar los defaults: proceso, referencia, entidad, politica"
echo "  Escribí uno por línea. Línea vacía para terminar."
echo ""

PAGE_TYPES=()
while true; do
    read -r -p "  Tipo: " ptype
    [[ -z "$ptype" ]] && break
    PAGE_TYPES+=("$ptype")
done

# Defaults si no ingresó tipos
if [ ${#PAGE_TYPES[@]} -eq 0 ]; then
    PAGE_TYPES=("proceso" "referencia" "entidad" "politica")
    echo "  Usando defaults: ${PAGE_TYPES[*]}"
fi

echo ""
echo -e "${YELLOW}→ Convenciones específicas del dominio${NC}"
echo "  Reglas de negocio que la IA debe respetar en TODAS las operaciones"
echo "  (ingest, query y lint). Son restricciones que no están escritas en"
echo "  ningún documento pero que vos sabés que siempre deben cumplirse."
echo ""
echo "  Ej: 'Todo proceso debe indicar el rol responsable de ejecutarlo'"
echo "      'Los roles siempre listan sus permisos asociados'"
echo "      'Todo expediente tiene un número único de 8 dígitos'"
echo ""
echo "  Podés dejarlo vacío ahora y agregar convenciones en el archivo de instrucciones después."
echo "  Escribí una por línea. Línea vacía para terminar."
echo ""

CONVENTIONS=()
while true; do
    read -r -p "  Convención: " conv
    [[ -z "$conv" ]] && break
    CONVENTIONS+=("$conv")
done

# ─────────────────────────────────────────
# 2. Crear estructura de directorios
# ─────────────────────────────────────────

echo ""
echo -e "${BLUE}Creando wiki en: ${WIKI_DIR}${NC}"

if [[ -d "${WIKI_DIR}" ]]; then
    echo -e "${YELLOW}Advertencia: el directorio ${WIKI_DIR} ya existe.${NC}"
    read -r -p "  ¿Continuar de todas formas? Puede sobreescribir archivos. [s/N]: " confirm
    [[ "$confirm" =~ ^[sS]$ ]] || { echo "Cancelado."; exit 0; }
fi

mkdir -p "${WIKI_DIR}/raw" "${WIKI_DIR}/wiki"
[[ "$CLI_CHOICE" == "1" || "$CLI_CHOICE" == "4" ]] && mkdir -p "${WIKI_DIR}/.claude/commands"
[[ "$CLI_CHOICE" == "2" || "$CLI_CHOICE" == "4" ]] && mkdir -p "${WIKI_DIR}/.opencode/commands"
[[ "$CLI_CHOICE" == "3" || "$CLI_CHOICE" == "4" ]] && mkdir -p "${WIKI_DIR}/.pi/prompts"

touch "${WIKI_DIR}/raw/.gitkeep"

# ─────────────────────────────────────────
# 3. Construir bloques de contenido
# ─────────────────────────────────────────

CREATED_DATE=$(date +%Y-%m-%d)
NL=$'\n'

# Entities list para YAML
ENTITIES_LIST=""
for e in "${ENTITIES[@]}"; do
    ENTITIES_LIST="${ENTITIES_LIST}  - ${e}${NL}"
done

# Page types list para YAML
PAGE_TYPES_LIST=""
for t in "${PAGE_TYPES[@]}"; do
    PAGE_TYPES_LIST="${PAGE_TYPES_LIST}  - ${t}${NL}"
done

# Page types detail — descripción de cada tipo
PAGE_TYPES_DETAIL=""
for t in "${PAGE_TYPES[@]}"; do
    case "$t" in
        proceso)
            PAGE_TYPES_DETAIL="${PAGE_TYPES_DETAIL}### \`proceso\`${NL}Describe cómo hacer algo: pasos secuenciales, precondiciones, actor responsable, resultado esperado.${NL}Slug: \`verbo-objeto.md\` (ej: \`crear-usuario.md\`)${NL}${NL}"
            ;;
        referencia)
            PAGE_TYPES_DETAIL="${PAGE_TYPES_DETAIL}### \`referencia\`${NL}Define qué es algo: términos, listas, tablas, configuraciones.${NL}Slug: \`sustantivo.md\` (ej: \`roles.md\`, \`permisos.md\`)${NL}${NL}"
            ;;
        entidad)
            PAGE_TYPES_DETAIL="${PAGE_TYPES_DETAIL}### \`entidad\`${NL}Describe un sistema, componente, actor o grupo específico del dominio.${NL}Slug: \`nombre-entidad.md\` (ej: \`sistema-renab.md\`)${NL}${NL}"
            ;;
        politica)
            PAGE_TYPES_DETAIL="${PAGE_TYPES_DETAIL}### \`politica\`${NL}Establece reglas, restricciones o lineamientos que deben cumplirse.${NL}Slug: \`politica-tema.md\` (ej: \`politica-acceso.md\`)${NL}${NL}"
            ;;
        regulacion)
            PAGE_TYPES_DETAIL="${PAGE_TYPES_DETAIL}### \`regulacion\`${NL}Documenta normativa legal o regulatoria aplicable al dominio. Debe citar artículo o fuente legal.${NL}Slug: \`regulacion-tema.md\` (ej: \`regulacion-datos-personales.md\`)${NL}${NL}"
            ;;
        reporte)
            PAGE_TYPES_DETAIL="${PAGE_TYPES_DETAIL}### \`reporte\`${NL}Resultado generado automáticamente por operaciones del wiki (lint, síntesis).${NL}Slug: \`lint-YYYY-MM-DD.md\` o \`reporte-tema.md\`${NL}${NL}"
            ;;
        *)
            PAGE_TYPES_DETAIL="${PAGE_TYPES_DETAIL}### \`${t}\`${NL}Tipo personalizado para este dominio.${NL}Slug: \`${t}-tema.md\`${NL}${NL}"
            ;;
    esac
done

# Domain conventions
DOMAIN_CONVENTIONS=""
if [ ${#CONVENTIONS[@]} -eq 0 ]; then
    DOMAIN_CONVENTIONS="> Sin convenciones específicas definidas al momento del setup.${NL}> Agregar aquí las reglas particulares de este dominio a medida que emerjan."
else
    for c in "${CONVENTIONS[@]}"; do
        DOMAIN_CONVENTIONS="${DOMAIN_CONVENTIONS}- ${c}${NL}"
    done
fi

# ─────────────────────────────────────────
# 4. Generar archivo(s) de instrucciones
# ─────────────────────────────────────────

generate_instructions() {
    local dest_file="$1"
    local commands_dir="$2"
    local instructions_file
    instructions_file="$(basename "${dest_file}")"

    cp "${SCRIPT_DIR}/schema.md.template" "${dest_file}"

    python3 - \
        "${dest_file}" \
        "${WIKI_NAME}" \
        "${WIKI_SLUG}" \
        "${LANGUAGE}" \
        "${CREATED_DATE}" \
        "${ENTITIES_LIST}" \
        "${PAGE_TYPES_LIST}" \
        "${PAGE_TYPES_DETAIL}" \
        "${DOMAIN_CONVENTIONS}" \
        "${commands_dir}" \
        "${instructions_file}" \
<<'PYEOF'
import sys

path, wiki_name, wiki_slug, language, created_date, \
    entities_list, page_types_list, page_types_detail, domain_conventions, \
    commands_dir, instructions_file = sys.argv[1:]

with open(path, "r") as f:
    content = f.read()

content = content.replace("{{WIKI_NAME}}", wiki_name)
content = content.replace("{{WIKI_SLUG}}", wiki_slug)
content = content.replace("{{WIKI_ROOT}}", wiki_slug + "-wiki")
content = content.replace("{{LANGUAGE}}", language)
content = content.replace("{{CREATED_DATE}}", created_date)
content = content.replace("{{ENTITIES_LIST}}", entities_list)
content = content.replace("{{PAGE_TYPES_LIST}}", page_types_list)
content = content.replace("{{PAGE_TYPES_DETAIL}}", page_types_detail)
content = content.replace("{{DOMAIN_CONVENTIONS}}", domain_conventions)
content = content.replace("{{COMMANDS_DIR}}", commands_dir)
content = content.replace("{{INSTRUCTIONS_FILE}}", instructions_file)

with open(path, "w") as f:
    f.write(content)
PYEOF
}

[[ "$CLI_CHOICE" == "1" || "$CLI_CHOICE" == "4" ]] && generate_instructions "${WIKI_DIR}/CLAUDE.md"  ".claude/commands"
[[ "$CLI_CHOICE" == "2" || "$CLI_CHOICE" == "4" ]] && generate_instructions "${WIKI_DIR}/AGENTS.md" ".opencode/commands"
[[ "$CLI_CHOICE" == "3" || "$CLI_CHOICE" == "4" ]] && generate_instructions "${WIKI_DIR}/AGENTS.md" ".pi/prompts"

# ─────────────────────────────────────────
# 5. Generar index.md y log.md
# ─────────────────────────────────────────

cat > "${WIKI_DIR}/wiki/index.md" << EOF
# Índice — ${WIKI_NAME}

> Catálogo central del wiki. La IA lo lee primero en cada operación.
> No editar manualmente — se actualiza automáticamente con cada \`/wiki-ingest\`.

| Página | Descripción | Tipo | Status | Actualizado |
|--------|-------------|------|--------|-------------|

<!-- Las páginas se agregan aquí automáticamente durante el ingest -->
EOF

cat > "${WIKI_DIR}/wiki/log.md" << EOF
# Log de Operaciones — ${WIKI_NAME}

> Historial append-only. Nunca modificar entradas anteriores.
> Se actualiza automáticamente con cada \`/wiki-ingest\` y \`/wiki-lint\`.

---

## ${CREATED_DATE} — setup

**Evento:** Wiki inicializado con setup.sh
**Dominio:** ${WIKI_NAME} (${WIKI_SLUG})
**Entidades primarias:** $(IFS=', '; echo "${ENTITIES[*]}")
**Tipos de página:** $(IFS=', '; echo "${PAGE_TYPES[*]}")

---
EOF

# ─────────────────────────────────────────
# 6. Copiar skills
# ─────────────────────────────────────────

copy_commands() {
    local dest_dir="$1"
    cp "${SCRIPT_DIR}/commands/wiki-ingest.md" "${dest_dir}/"
    cp "${SCRIPT_DIR}/commands/wiki-query.md"  "${dest_dir}/"
    cp "${SCRIPT_DIR}/commands/wiki-lint.md"   "${dest_dir}/"
}

[[ "$CLI_CHOICE" == "1" || "$CLI_CHOICE" == "4" ]] && copy_commands "${WIKI_DIR}/.claude/commands"
[[ "$CLI_CHOICE" == "2" || "$CLI_CHOICE" == "4" ]] && copy_commands "${WIKI_DIR}/.opencode/commands"
[[ "$CLI_CHOICE" == "3" || "$CLI_CHOICE" == "4" ]] && copy_commands "${WIKI_DIR}/.pi/prompts"

# ─────────────────────────────────────────
# 7. Inicializar Git
# ─────────────────────────────────────────

cd "${WIKI_DIR}"

cat > .gitignore << EOF
.DS_Store
*.swp
*.tmp
EOF

git init -q
git add .
git commit -q -m "chore: init wiki ${WIKI_SLUG}"

cd - > /dev/null

# ─────────────────────────────────────────
# 8. Resumen final
# ─────────────────────────────────────────

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   ✓ Wiki creado exitosamente                 ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  📁 Directorio: ${BLUE}${WIKI_DIR}${NC}"
echo -e "  📂 Fuentes:    ${BLUE}${WIKI_DIR}/raw/${NC}"
echo -e "  📂 Wiki:       ${BLUE}${WIKI_DIR}/wiki/${NC}"
echo ""

if [[ "$CLI_CHOICE" == "1" || "$CLI_CHOICE" == "4" ]]; then
    echo -e "  ${YELLOW}Claude Code${NC}"
    echo -e "    📄 Schema:   ${BLUE}${WIKI_DIR}/CLAUDE.md${NC}"
    echo -e "    ⚙️  Commands: ${BLUE}${WIKI_DIR}/.claude/commands/${NC}"
fi
if [[ "$CLI_CHOICE" == "2" || "$CLI_CHOICE" == "4" ]]; then
    echo -e "  ${YELLOW}OpenCode${NC}"
    echo -e "    📄 Schema:   ${BLUE}${WIKI_DIR}/AGENTS.md${NC}"
    echo -e "    ⚙️  Commands: ${BLUE}${WIKI_DIR}/.opencode/commands/${NC}"
fi
if [[ "$CLI_CHOICE" == "3" || "$CLI_CHOICE" == "4" ]]; then
    echo -e "  ${YELLOW}Pi${NC}"
    echo -e "    📄 Schema:   ${BLUE}${WIKI_DIR}/AGENTS.md${NC}"
    echo -e "    ⚙️  Commands: ${BLUE}${WIKI_DIR}/.pi/prompts/${NC}"
fi

echo ""
echo -e "${YELLOW}Próximos pasos:${NC}"
echo "  1. Copia tus documentos existentes en raw/"

if [[ "$CLI_CHOICE" == "1" ]]; then
    echo "  2. Abre Claude Code en el directorio del wiki"
    echo "  3. Ejecuta: /wiki-ingest"
    echo "  4. Pregunta lo que necesites: /wiki-query ¿qué roles existen?"
elif [[ "$CLI_CHOICE" == "2" ]]; then
    echo "  2. Abre una terminal en el directorio del wiki y ejecuta: opencode"
    echo "  3. Ejecuta: /wiki-ingest"
    echo "  4. Pregunta lo que necesites: /wiki-query ¿qué roles existen?"
elif [[ "$CLI_CHOICE" == "3" ]]; then
    echo "  2. Abre una terminal en el directorio del wiki y ejecuta: pi"
    echo "  3. Ejecuta: /wiki-ingest"
    echo "  4. Pregunta lo que necesites: /wiki-query ¿qué roles existen?"
else
    echo "  2. Abre el CLI de tu preferencia en el directorio del wiki"
    echo "  3. Ejecuta: /wiki-ingest"
    echo "  4. Pregunta lo que necesites: /wiki-query ¿qué roles existen?"
fi
echo ""
