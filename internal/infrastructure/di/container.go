package di

import (
	"ddsheetfinal/internal/application/usecases"
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/services"
	"ddsheetfinal/internal/infrastructure/external"
	"ddsheetfinal/internal/infrastructure/persistence"
	"ddsheetfinal/internal/presentation/cli"
)

// Container holds all dependencies for the application
type Container struct {
	// Repositories
	CharacterRepo  repositories.CharacterRepository
	SRDRepo        repositories.SRDRepository
	EnrichmentRepo repositories.EnrichmentRepository

	// Services
	CharacterService *services.CharacterService

	// Use Cases
	CreateCharacterUseCase      *usecases.CreateCharacterUseCase
	ViewCharacterUseCase        *usecases.ViewCharacterUseCase
	GetCharacterUseCase         *usecases.GetCharacterUseCase
	SaveCharacterUseCase        *usecases.SaveCharacterUseCase
	DeleteCharacterUseCase      *usecases.DeleteCharacterUseCase
	ListCharactersUseCase       *usecases.ListCharactersUseCase
	EquipItemUseCase            *usecases.EquipItemUseCase
	UnequipItemUseCase          *usecases.UnequipItemUseCase
	LearnSpellUseCase           *usecases.LearnSpellUseCase
	PrepareSpellUseCase         *usecases.PrepareSpellUseCase
	GetLearnableSpellsUseCase   *usecases.GetLearnableSpellsUseCase
	EnrichCharacterUseCase      *usecases.EnrichCharacterUseCase
	DeriveStatsUseCase          *usecases.DeriveCharacterStatsUseCase
	EnrichDataUseCase           *usecases.EnrichDataUseCase

	// CLI Handlers
	CreateHandler          *cli.CreateHandler
	ViewHandler            *cli.ViewHandler
	DeleteHandler          *cli.DeleteHandler
	ListHandler            *cli.ListHandler
	EquipHandler           *cli.EquipHandler
	UnequipHandler         *cli.UnequipHandler
	LearnSpellHandler      *cli.LearnSpellHandler
	PrepareSpellHandler    *cli.PrepareSpellHandler
	LearnableSpellsHandler *cli.LearnableSpellsHandler
	HelpHandler            *cli.HelpHandler
	ServeHandler           *cli.ServeHandler
	EnrichHandler          *cli.EnrichHandler
}

// NewContainer creates and wires up all dependencies
func NewContainer(charactersDir, equipmentCSV, spellsCSV, cacheDir string) *Container {
	c := &Container{}

	// Initialize repositories (Infrastructure layer)
	c.CharacterRepo = persistence.NewJSONCharacterRepository(charactersDir)
	c.SRDRepo = persistence.NewCSVSRDRepository(equipmentCSV, spellsCSV)
	c.EnrichmentRepo = external.NewDND5EAPIEnrichmentRepository(cacheDir)

	// Initialize domain services
	c.CharacterService = services.NewCharacterService()

	// Initialize use cases (Application layer)
	c.CreateCharacterUseCase = usecases.NewCreateCharacterUseCase(c.CharacterRepo)
	c.ViewCharacterUseCase = usecases.NewViewCharacterUseCase(c.CharacterRepo, c.CharacterService)
	c.GetCharacterUseCase = usecases.NewGetCharacterUseCase(c.CharacterRepo)
	c.SaveCharacterUseCase = usecases.NewSaveCharacterUseCase(c.CharacterRepo)
	c.DeleteCharacterUseCase = usecases.NewDeleteCharacterUseCase(c.CharacterRepo)
	c.ListCharactersUseCase = usecases.NewListCharactersUseCase(c.CharacterRepo)
	c.EquipItemUseCase = usecases.NewEquipItemUseCase(c.CharacterRepo, c.SRDRepo)
	c.UnequipItemUseCase = usecases.NewUnequipItemUseCase(c.CharacterRepo)
	c.LearnSpellUseCase = usecases.NewLearnSpellUseCase(c.CharacterRepo, c.SRDRepo)
	c.PrepareSpellUseCase = usecases.NewPrepareSpellUseCase(c.CharacterRepo, c.SRDRepo)
	c.GetLearnableSpellsUseCase = usecases.NewGetLearnableSpellsUseCase(c.CharacterRepo, c.SRDRepo)
	c.EnrichCharacterUseCase = usecases.NewEnrichCharacterUseCase(c.CharacterRepo, c.EnrichmentRepo)
	c.DeriveStatsUseCase = usecases.NewDeriveCharacterStatsUseCase(c.CharacterRepo, c.CharacterService)
	c.EnrichDataUseCase = usecases.NewEnrichDataUseCase(c.EnrichmentRepo)

	// Initialize CLI handlers (Presentation layer)
	c.CreateHandler = cli.NewCreateHandler(c.CreateCharacterUseCase)
	c.ViewHandler = cli.NewViewHandler(c.ViewCharacterUseCase)
	c.DeleteHandler = cli.NewDeleteHandler(c.DeleteCharacterUseCase)
	c.ListHandler = cli.NewListHandler(c.ListCharactersUseCase)
	c.EquipHandler = cli.NewEquipHandler(c.EquipItemUseCase)
	c.UnequipHandler = cli.NewUnequipHandler(c.UnequipItemUseCase)
	c.LearnSpellHandler = cli.NewLearnSpellHandler(c.LearnSpellUseCase)
	c.PrepareSpellHandler = cli.NewPrepareSpellHandler(c.PrepareSpellUseCase)
	c.LearnableSpellsHandler = cli.NewLearnableSpellsHandler(c.GetLearnableSpellsUseCase)
	c.HelpHandler = cli.NewHelpHandler()
	c.ServeHandler = cli.NewServeHandler(
		c.GetCharacterUseCase,
		c.SaveCharacterUseCase,
		c.DeleteCharacterUseCase,
		c.ListCharactersUseCase,
		c.DeriveStatsUseCase,
		c.EnrichDataUseCase,
	)
	c.EnrichHandler = cli.NewEnrichHandler(c.EnrichCharacterUseCase)

	return c
}
