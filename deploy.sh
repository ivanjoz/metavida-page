#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")"

GO_PATH="go"
if [ -x /usr/local/go/bin/go ]; then
    GO_PATH="/usr/local/go/bin/go"
fi

echo "=== METAVIDA DEPLOYMENT ==="
echo "Seleccione acciones a realizar (ej: '12'):"
echo "[1] Generar schemas.generated.go"
echo "[2] Desplegar tablas e índices en Cloudflare D1"
echo "[3] Actualizar ui-components y libs desde Genix"
echo "[4] Backend (VPS)"
echo "[5] Build frontend estático a /docs (GitHub Pages)"
echo "[6] Desplegar tablas e insertar usuario admin inicial"
echo "[7] Configurar servicio backend en VPS"
read -r ACCIONES

if [[ "$ACCIONES" == *"3"* ]]; then
    echo "=== ACTUALIZANDO COMPONENTES COMPARTIDOS DESDE GENIX ==="
    GENIX_URL="${GENIX_URL:-https://github.com/ivanjoz/genix.git}"
    GENIX_TMP_DIR=""
    if [ -n "${GENIX_REPO:-}" ]; then
        if [ ! -d "$GENIX_REPO/.git" ]; then
            echo "No se encontró el repo Genix en: $GENIX_REPO"
            exit 1
        fi
    else
        GENIX_TMP_DIR="$(mktemp -d)"
        trap 'if [ -n "${GENIX_TMP_DIR:-}" ]; then rm -rf "$GENIX_TMP_DIR"; fi' EXIT
        git clone --depth 1 --filter=blob:none --sparse "$GENIX_URL" "$GENIX_TMP_DIR"
        git -C "$GENIX_TMP_DIR" sparse-checkout set frontend/ui-components frontend/libs
        GENIX_REPO="$GENIX_TMP_DIR"
    fi

    # git archive only overwrites files it extracts; remove old vendored paths first.
    # Keep Metavida-owned files in frontend/libs, especially frontend/libs/http.ts.
    rm -rf frontend/ui-components \
        frontend/libs/assets \
        frontend/libs/cache \
        frontend/libs/excel \
        frontend/libs/funcs \
        frontend/libs/vendor \
        frontend/libs/workers \
        frontend/libs/fecha.ts \
        frontend/libs/helpers.ts \
        frontend/libs/http.svelte.ts \
        frontend/libs/sharedHelpers.ts \
        frontend/libs/sse-client.ts \
        frontend/libs/sw-cache.ts \
        frontend/libs/typed-idb.ts

    mkdir -p frontend
    git -C "$GENIX_REPO" archive HEAD frontend/ui-components frontend/libs \
        | tar -x --strip-components=1 -C frontend
    git -C "$GENIX_REPO" rev-parse HEAD > frontend/ui-components/GENIX_COMMIT
    echo "Componentes actualizados desde Genix commit: $(cat frontend/ui-components/GENIX_COMMIT)"
fi

if [[ "$ACCIONES" == *"5"* ]]; then
    echo "=== BUILD FRONTEND ESTÁTICO -> /docs ==="
    (
        cd frontend
        bun install
        bun run build:docs
    )
    echo "✅ Build estático listo en /docs (commit y push para publicar en GitHub Pages)."
fi

if [[ "$ACCIONES" == *"4"* ]]; then
    echo "=== PUBLICANDO BACKEND (VPS) ==="
    cd ./scripts
    "$GO_PATH" run . deploy_vps
    cd ..
    echo "✅ El deploy VPS finalizado!"
fi

if [[ "$ACCIONES" == *"7"* ]]; then
    echo "=== CONFIGURANDO SERVICIO BACKEND (VPS) ==="
    cd ./scripts
    "$GO_PATH" run . configure_server
    cd ..
    echo "✅ La configuración VPS finalizó!"
fi

if [[ "$ACCIONES" == *"1"* || "$ACCIONES" == *"2"* || "$ACCIONES" == *"6"* ]]; then
    echo "=== GENERANDO SCHEMAS ==="
    (cd scripts && "$GO_PATH" run . generate_schemas)
fi

if [[ "$ACCIONES" == *"2"* || "$ACCIONES" == *"6"* ]]; then
    echo "=== DESPLEGANDO TABLAS EN CLOUDFLARE D1 ==="
    (cd backend && METAVIDA_CREDENTIALS=../credentials.json "$GO_PATH" run . deploy_tables)
fi

if [[ "$ACCIONES" == *"6"* ]]; then
    echo "=== INSERTANDO USUARIO ADMIN INICIAL ==="
    (cd backend && METAVIDA_CREDENTIALS=../credentials.json "$GO_PATH" run . insert_admin)
fi

echo "Finalizado."
