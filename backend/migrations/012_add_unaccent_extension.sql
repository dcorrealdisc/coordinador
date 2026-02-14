-- Migration 012: Agregar extensión unaccent para búsquedas insensibles a acentos
CREATE EXTENSION IF NOT EXISTS unaccent;
