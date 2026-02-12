# 📦 Guía para Crear Repositorio en GitHub

## Paso 1: Configurar Git Local (PRIMERO)

```bash
cd /home/dcorreal/Develop/coordinador

# Configurar tu identidad (usa tu email real de GitHub)
git config user.name "Dario Correal"
git config user.email "TU_EMAIL@example.com"  # ⚠️ Cambia esto por tu email real

# Verificar configuración
git config --list | grep user
```

## Paso 2: Crear el Commit Inicial

```bash
# Los archivos ya están en staging, solo falta hacer commit
git commit -m "feat: setup inicial del proyecto Coordinador

- Arquitectura monolito modular con CQRS light
- Backend Go/Fiber con estructura de módulos
- Dos frontends Svelte (admin-web y portal-web)
- Docker Compose para PostgreSQL
- Documentación completa (README, ADRs, guías)
- Agente Arquitecto configurado
- Makefile con comandos de desarrollo
- ADR-001: Decisiones arquitectónicas fundamentales"

# Verificar que el commit se creó
git log --oneline
```

## Paso 3: Crear Repositorio en GitHub

### Opción A: Via Web (Recomendado)

1. **Ve a GitHub**: https://github.com/new

2. **Configuración del repositorio**:
   - **Repository name**: `coordinador`
   - **Description**: `Sistema de Gestión Académica para Maestrías - Go/Fiber + Svelte + PostgreSQL`
   - **Visibility**: 
     - ✅ **Public** (si quieres compartir)
     - ✅ **Private** (si es solo para ti)
   - **NO marques** "Initialize with README" (ya lo tenemos)
   - **NO agregues** .gitignore (ya lo tenemos)
   - **NO agregues** license (lo haremos después si quieres)

3. **Click en "Create repository"**

### Opción B: Via GitHub CLI (si lo tienes instalado)

```bash
# Crear repositorio privado
gh repo create coordinador --private --source=. --remote=origin

# O crear repositorio público
gh repo create coordinador --public --source=. --remote=origin
```

## Paso 4: Conectar Local con GitHub

GitHub te mostrará instrucciones, pero básicamente:

```bash
# Desde /home/dcorreal/Develop/coordinador

# Agregar remote (GitHub te dará la URL exacta)
git remote add origin https://github.com/TU_USUARIO/coordinador.git

# O si usas SSH (recomendado):
git remote add origin git@github.com:TU_USUARIO/coordinador.git

# Verificar remote
git remote -v

# Push del commit inicial
git push -u origin main
```

## Paso 5: Verificar en GitHub

1. Ve a `https://github.com/TU_USUARIO/coordinador`
2. Deberías ver:
   - ✅ README.md renderizado
   - ✅ Toda la estructura de carpetas
   - ✅ Tu commit inicial
   - ✅ 21 archivos

## 🔐 Configurar Autenticación SSH (Recomendado)

Si aún no tienes SSH configurado con GitHub:

```bash
# 1. Generar clave SSH (si no tienes)
ssh-keygen -t ed25519 -C "tu_email@example.com"
# Presiona Enter para aceptar ubicación por defecto
# Opcionalmente agrega passphrase

# 2. Copiar clave pública
cat ~/.ssh/id_ed25519.pub

# 3. Agregar en GitHub:
# - Ve a https://github.com/settings/keys
# - Click "New SSH key"
# - Pega la clave pública
# - Click "Add SSH key"

# 4. Verificar conexión
ssh -T git@github.com
# Deberías ver: "Hi TU_USUARIO! You've successfully authenticated..."
```

## 📥 Clonar en Otra Máquina

Cuando estés en tu oficina:

```bash
# Con SSH (recomendado)
git clone git@github.com:TU_USUARIO/coordinador.git
cd coordinador

# Con HTTPS
git clone https://github.com/TU_USUARIO/coordinador.git
cd coordinador

# Instalar dependencias
make install

# Levantar PostgreSQL
make db-up

# Iniciar desarrollo
make dev-backend
```

## 🔄 Workflow Diario

### En el portátil (trabajar):
```bash
# Obtener últimos cambios (si trabajaste en oficina)
git pull origin main

# Trabajar...
# Hacer commits...
git add .
git commit -m "feat: implementar módulo de estudiantes"

# Subir cambios
git push origin main
```

### En la oficina (continuar):
```bash
# Obtener lo que hiciste en el portátil
git pull origin main

# Continuar trabajando...
```

## 🌿 Branches Recomendadas

Para desarrollo más organizado:

```bash
# Crear rama develop
git checkout -b develop
git push -u origin develop

# Para nuevas features
git checkout -b feature/nombre-feature
# Trabajar...
git push -u origin feature/nombre-feature
# Hacer PR a develop en GitHub
```

## 🏷️ Tags para Versiones

Cuando completes fases importantes:

```bash
# Crear tag
git tag -a v0.1.0 -m "Setup inicial del proyecto"
git push origin v0.1.0

# Listar tags
git tag -l
```

## ⚠️ Archivos que NO se subirán (por .gitignore)

Estos archivos/carpetas están ignorados:
- `node_modules/`
- `bin/`
- `*.log`
- `.env` (variables de entorno)
- `coverage.out`
- `.DS_Store`

**¡Perfecto!** Así proteges información sensible.

## 🔍 Verificación Post-Setup

```bash
# Ver estado del repo
git status

# Ver historial
git log --oneline --graph

# Ver remotes
git remote -v

# Ver todas las ramas
git branch -a
```

## 💡 Tips

1. **Commits frecuentes**: Haz commits pequeños y frecuentes
2. **Mensajes descriptivos**: Usa Conventional Commits (feat, fix, docs, etc.)
3. **Pull antes de push**: Siempre haz `git pull` antes de `git push`
4. **Branches para features**: Usa branches para experimentar
5. **README actualizado**: Mantén el README actualizado

## 📞 Ayuda

Si algo falla:

```bash
# Ver qué está pasando
git status
git log --oneline

# Deshacer cambios no commiteados
git reset --hard HEAD

# Deshacer último commit (mantiene cambios)
git reset --soft HEAD~1
```

---

## ✅ Checklist Final

Antes de considerar que está todo listo:

- [ ] Git configurado con tu nombre y email
- [ ] Commit inicial creado localmente
- [ ] Repositorio creado en GitHub
- [ ] Remote `origin` configurado
- [ ] Push exitoso a GitHub
- [ ] README se ve bien en GitHub
- [ ] SSH configurado (opcional pero recomendado)
- [ ] Clonado exitosamente en otra máquina (cuando estés en la oficina)

**¡Listo! Ya puedes trabajar desde cualquier lugar.** 🚀
