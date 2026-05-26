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
echo "[6] Desplegar tablas e insertar usuario admin inicial"
read -r ACCIONES

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
