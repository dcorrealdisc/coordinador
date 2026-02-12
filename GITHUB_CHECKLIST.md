# ✅ PASOS PARA SUBIR A GITHUB - Lista de Verificación

## 📋 Resumen Rápido

Tu proyecto ya está listo para GitHub. Solo necesitas:
1. Configurar git con tu identidad
2. Hacer el commit inicial
3. Crear repo en GitHub
4. Conectar y subir

---

## 🚀 Paso a Paso (5 minutos)

### Paso 1: Configurar Git (2 minutos)

```bash
cd /home/dcorreal/Develop/coordinador

# Opción A: Script automático (RECOMENDADO)
./setup-git.sh

# Opción B: Manual
git config user.name "Dario Correal"
git config user.email "TU_EMAIL_DE_GITHUB@example.com"  # ⚠️ USA TU EMAIL REAL
```

### Paso 2: Verificar Archivos (30 segundos)

```bash
# Ver qué se va a commitear
git status

# Deberías ver 24 archivos listos
```

### Paso 3: Commit Inicial (30 segundos)

```bash
# Si usaste el script setup-git.sh, ya está hecho ✅
# Si no, ejecuta:

git commit -m "feat: setup inicial del proyecto Coordinador

- Arquitectura monolito modular con CQRS light
- Backend Go/Fiber con estructura de módulos
- Dos frontends Svelte (admin-web y portal-web)
- Docker Compose para PostgreSQL
- Documentación completa (README, ADRs, guías)
- Agente Arquitecto configurado
- Makefile con comandos de desarrollo
- ADR-001: Decisiones arquitectónicas fundamentales"
```

### Paso 4: Crear Repo en GitHub (1 minuto)

1. Ve a: **https://github.com/new**

2. Configura:
   - **Name**: `coordinador`
   - **Description**: `Sistema de Gestión Académica para Maestrías - Go/Fiber + Svelte + PostgreSQL`
   - **Visibility**: Private o Public (tú decides)
   - **NO marques** "Add README" (ya lo tenemos)
   - **NO marques** "Add .gitignore" (ya lo tenemos)

3. Click **"Create repository"**

### Paso 5: Conectar y Subir (1 minuto)

GitHub te mostrará comandos, pero básicamente:

```bash
# Agregar remote (GitHub te dará tu URL exacta)
git remote add origin git@github.com:TU_USUARIO/coordinador.git

# O si prefieres HTTPS:
git remote add origin https://github.com/TU_USUARIO/coordinador.git

# Subir a GitHub
git push -u origin main
```

### Paso 6: Verificar en GitHub (30 segundos)

Ve a: `https://github.com/TU_USUARIO/coordinador`

Deberías ver:
- ✅ README.md bien formateado
- ✅ 24 archivos
- ✅ Estructura de carpetas completa
- ✅ Tu commit inicial

---

## 🎯 ¡LISTO! Ahora puedes:

```bash
# En cualquier máquina:
git clone git@github.com:TU_USUARIO/coordinador.git
cd coordinador
make install
make dev-backend
```

---

## 🔐 Bonus: Configurar SSH (Recomendado)

Si prefieres no escribir contraseña cada vez:

```bash
# 1. Generar clave SSH
ssh-keygen -t ed25519 -C "tu_email@example.com"

# 2. Copiar la clave pública
cat ~/.ssh/id_ed25519.pub

# 3. En GitHub:
# - Ve a: https://github.com/settings/keys
# - Click "New SSH key"
# - Pega la clave
# - Click "Add SSH key"

# 4. Verificar
ssh -T git@github.com
```

---

## 📁 Archivos Creados

Total: **24 archivos**

### Documentación (9):
- ✅ README.md
- ✅ QUICKSTART.md
- ✅ CONTRIBUTING.md
- ✅ PROJECT_STATUS.md
- ✅ GITHUB_SETUP.md (guía detallada)
- ✅ GIT_REFERENCE.md (comandos útiles)
- ✅ backend/README.md
- ✅ docs/adrs/README.md + ADR-001
- ✅ docs/agents/README.md + Agente Arquitecto

### Código (3):
- ✅ backend/go.mod
- ✅ backend/cmd/api/main.go
- ✅ admin-web/package.json
- ✅ portal-web/package.json

### Configuración (5):
- ✅ .gitignore
- ✅ Makefile
- ✅ docker-compose.yml
- ✅ setup-git.sh
- ✅ varios .gitkeep

---

## ❓ Si Algo Sale Mal

### "Git dice que no sabe quién soy"
```bash
git config user.name "Tu Nombre"
git config user.email "tu@email.com"
```

### "No puedo hacer push"
```bash
# Si es primera vez, asegúrate de:
git remote -v  # Verificar que origin está configurado

# Si ves "permission denied":
# Necesitas configurar SSH o usar HTTPS con token
```

### "Olvidé configurar git antes del commit"
```bash
git config user.name "Tu Nombre"
git config user.email "tu@email.com"
git commit --amend --reset-author --no-edit
```

---

## 📞 Ayuda Completa

- **Guía detallada**: Ver [GITHUB_SETUP.md](./GITHUB_SETUP.md)
- **Comandos Git**: Ver [GIT_REFERENCE.md](./GIT_REFERENCE.md)
- **Problemas**: Abre un issue o revisa la documentación

---

## ✨ Próximo Paso

Una vez subido a GitHub:

1. **En tu oficina**: `git clone` del repo
2. **Instalar**: `make install`
3. **Desarrollar**: Seguir [PROJECT_STATUS.md](./PROJECT_STATUS.md)

**Siguiente fase**: Crear Agente DBA y diseñar base de datos 🗄️

---

**Estado actual**: Todo listo para subir a GitHub ✅
