package usecases

import (
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/services"
)

// ServeHTTPUseCase provides dependencies for the HTTP server
type ServeHTTPUseCase struct {
	characterRepo    repositories.CharacterRepository
	srdRepo          repositories.SRDRepository
	enrichmentRepo   repositories.EnrichmentRepository
	characterService *services.CharacterService
}

// NewServeHTTPUseCase creates a new serve HTTP use case
func NewServeHTTPUseCase(
	characterRepo repositories.CharacterRepository,
	srdRepo repositories.SRDRepository,
	enrichmentRepo repositories.EnrichmentRepository,
	characterService *services.CharacterService,
) *ServeHTTPUseCase {
	return &ServeHTTPUseCase{
		characterRepo:    characterRepo,
		srdRepo:          srdRepo,
		enrichmentRepo:   enrichmentRepo,
		characterService: characterService,
	}
}

// GetCharacterRepository returns the character repository
func (uc *ServeHTTPUseCase) GetCharacterRepository() repositories.CharacterRepository {
	return uc.characterRepo
}

// GetSRDRepository returns the SRD repository
func (uc *ServeHTTPUseCase) GetSRDRepository() repositories.SRDRepository {
	return uc.srdRepo
}

// GetEnrichmentRepository returns the enrichment repository
func (uc *ServeHTTPUseCase) GetEnrichmentRepository() repositories.EnrichmentRepository {
	return uc.enrichmentRepo
}

// GetCharacterService returns the character service
func (uc *ServeHTTPUseCase) GetCharacterService() *services.CharacterService {
	return uc.characterService
}
