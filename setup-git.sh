#!/bin/bash

# Script de Configuración Rápida de Git
# =====================================

echo "🔧 Configuración de Git para Coordinador"
echo ""

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar si git está configurado
if git config user.name > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Git ya está configurado:${NC}"
    echo "   Nombre: $(git config user.name)"
    echo "   Email: $(git config user.email)"
    echo ""
    read -p "¿Quieres cambiar la configuración? (y/N): " change
    if [[ ! $change =~ ^[Yy]$ ]]; then
        echo "Manteniendo configuración actual."
        exit 0
    fi
fi

# Solicitar nombre
echo ""
echo -e "${BLUE}Configura tu identidad de Git:${NC}"
read -p "Tu nombre completo: " git_name
read -p "Tu email (usa el mismo de GitHub): " git_email

# Configurar git
git config user.name "$git_name"
git config user.email "$git_email"

echo ""
echo -e "${GREEN}✅ Git configurado exitosamente:${NC}"
echo "   Nombre: $(git config user.name)"
echo "   Email: $(git config user.email)"
echo ""

# Preguntar si quiere hacer el commit inicial
read -p "¿Quieres hacer el commit inicial ahora? (Y/n): " do_commit
if [[ ! $do_commit =~ ^[Nn]$ ]]; then
    echo ""
    echo "📦 Creando commit inicial..."
    
    git add .
    git commit -m "feat: setup inicial del proyecto Coordinador

- Arquitectura monolito modular con CQRS light
- Backend Go/Fiber con estructura de módulos
- Dos frontends Svelte (admin-web y portal-web)
- Docker Compose para PostgreSQL
- Documentación completa (README, ADRs, guías)
- Agente Arquitecto configurado
- Makefile con comandos de desarrollo
- ADR-001: Decisiones arquitectónicas fundamentales"
    
    echo ""
    echo -e "${GREEN}✅ Commit inicial creado${NC}"
    echo ""
    echo -e "${BLUE}📋 Próximos pasos:${NC}"
    echo "1. Crea el repositorio en GitHub: https://github.com/new"
    echo "2. Ejecuta: git remote add origin git@github.com:TU_USUARIO/coordinador.git"
    echo "3. Ejecuta: git push -u origin main"
    echo ""
    echo "Ver GITHUB_SETUP.md para instrucciones detalladas"
else
    echo ""
    echo "Puedes hacer el commit cuando quieras con:"
    echo "  git commit -m 'feat: setup inicial del proyecto'"
fi

echo ""
echo -e "${GREEN}✅ Configuración completa${NC}"
