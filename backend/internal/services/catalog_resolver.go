package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/dcorreal/coordinador/internal/repositories"
)

// CatalogResolver resolves human-readable names to UUIDs, auto-creating missing entries.
type CatalogResolver struct {
	repo repositories.CatalogRepository

	countryCache  map[string]uuid.UUID
	cityCache     map[string]uuid.UUID
	profCache     map[string]uuid.UUID
	jobCache      map[string]uuid.UUID
	uniCache      map[string]uuid.UUID
}

// NewCatalogResolver creates a new CatalogResolver with empty caches.
func NewCatalogResolver(repo repositories.CatalogRepository) *CatalogResolver {
	return &CatalogResolver{
		repo:         repo,
		countryCache: make(map[string]uuid.UUID),
		cityCache:    make(map[string]uuid.UUID),
		profCache:    make(map[string]uuid.UUID),
		jobCache:     make(map[string]uuid.UUID),
		uniCache:     make(map[string]uuid.UUID),
	}
}

func cacheKey(parts ...string) string {
	normalized := make([]string, len(parts))
	for i, p := range parts {
		normalized[i] = strings.ToLower(strings.TrimSpace(p))
	}
	return strings.Join(normalized, "|")
}

func (r *CatalogResolver) ResolveCountry(ctx context.Context, name string) (uuid.UUID, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, nil
	}
	key := cacheKey(name)
	if id, ok := r.countryCache[key]; ok {
		return id, nil
	}

	id, err := r.repo.FindCountryByName(ctx, name)
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		id, err = r.repo.CreateCountry(ctx, name)
		if err != nil {
			return uuid.Nil, err
		}
	}
	r.countryCache[key] = id
	return id, nil
}

func (r *CatalogResolver) ResolveCity(ctx context.Context, name string, countryID uuid.UUID) (uuid.UUID, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, nil
	}
	key := cacheKey(name, countryID.String())
	if id, ok := r.cityCache[key]; ok {
		return id, nil
	}

	id, err := r.repo.FindCityByName(ctx, name, countryID)
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		id, err = r.repo.CreateCity(ctx, name, countryID)
		if err != nil {
			return uuid.Nil, err
		}
	}
	r.cityCache[key] = id
	return id, nil
}

func (r *CatalogResolver) ResolveProfession(ctx context.Context, name string) (uuid.UUID, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, nil
	}
	key := cacheKey(name)
	if id, ok := r.profCache[key]; ok {
		return id, nil
	}

	id, err := r.repo.FindProfessionByName(ctx, name)
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		id, err = r.repo.CreateProfession(ctx, name)
		if err != nil {
			return uuid.Nil, err
		}
	}
	r.profCache[key] = id
	return id, nil
}

func (r *CatalogResolver) ResolveJobTitleCategory(ctx context.Context, name string) (uuid.UUID, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, nil
	}
	key := cacheKey(name)
	if id, ok := r.jobCache[key]; ok {
		return id, nil
	}

	id, err := r.repo.FindJobTitleCategoryByName(ctx, name)
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		id, err = r.repo.CreateJobTitleCategory(ctx, name)
		if err != nil {
			return uuid.Nil, err
		}
	}
	r.jobCache[key] = id
	return id, nil
}

func (r *CatalogResolver) ResolveUniversity(ctx context.Context, name string, cityID *uuid.UUID, countryID uuid.UUID) (uuid.UUID, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, nil
	}
	key := cacheKey(name, countryID.String())
	if id, ok := r.uniCache[key]; ok {
		return id, nil
	}

	id, err := r.repo.FindUniversityByName(ctx, name, countryID)
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		id, err = r.repo.CreateUniversity(ctx, name, cityID, countryID)
		if err != nil {
			return uuid.Nil, err
		}
	}
	r.uniCache[key] = id
	return id, nil
}
