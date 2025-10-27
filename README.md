# D&D Character Sheet Manager

A command-line tool for managing D&D 5e characters

## Project Structure

```
DDsheetfinal/
├── cmd/                        # Application entry points
│   └── ddsheet/               # Main CLI application
│       └── main.go
├── internal/                   # Private application code
│   ├── domain/                # Core business logic (no dependencies)
│   │   ├── entities/          # Character and other domain entities
│   │   ├── valueobjects/      # Immutable values (abilities, armor, SRD data)
│   │   ├── repositories/      # Repository interfaces
│   │   └── services/          # Domain services (AC calculation, etc.)
│   ├── application/           # Use cases and application logic
│   │   └── usecases/          # Command handlers (create, view, equip, etc.)
│   ├── infrastructure/        # External implementations
│   │   ├── persistence/       # File storage, CSV access
│   │   ├── external/          # API clients (dnd5eapi)
│   │   └── di/                # Dependency injection container
│   └── presentation/          # User interfaces
│       └── cli/               # Command-line handlers
├── assets/                    # Static resources
│   └── srd/                   # D&D 5e SRD data files
│       ├── 5e-SRD-Equipment.csv
│       └── 5e-SRD-Spells.csv
├── data/                      # Runtime data (gitignored)
│   ├── characters/            # Character JSON files
│   ├── enrichments/           # Enriched character data
│   └── cache/                 # API response cache
├── web/                       # Web frontend
│   ├── static/                # HTML, CSS, JS files
│   └── manifest.json
├── go.mod
└── README.md
```

## Architecture

This project follows **Onion Architecture** (also known as Clean Architecture):

- **Domain Layer** (Core): Pure business logic with no external dependencies
- **Application Layer**: Use cases that orchestrate domain logic
- **Infrastructure Layer**: Concrete implementations (file storage, APIs, etc.)
- **Presentation Layer**: User interfaces (CLI, HTTP handlers)

**Dependency Rule**: All dependencies point inward toward the domain. The domain has no knowledge of outer layers.

## Building

```bash
go build -o ddsheet ./cmd/ddsheet
```

## Usage

```bash
# Create a character
./ddsheet create -name Gandalf -race human -class wizard -level 20 -int 20 -wis 18 -cha 14

# View character details
./ddsheet view -name Gandalf

# List all characters
./ddsheet list

# Equip items
./ddsheet equip -name Gandalf -weapon Quarterstaff -armor "robe"

# Learn and prepare spells
./ddsheet learn-spell -name Gandalf -spell "fireball"
./ddsheet prepare-spell -name Gandalf -spell "fireball"

# Delete a character
./ddsheet delete -name Gandalf

# Show all commands
./ddsheet help
```

## Features

- Character creation with racial bonuses and skill proficiencies
- Equipment management with automatic AC calculation
- Spell learning and preparation (class-specific)
- Derived stat calculation (initiative, passive perception, etc.)
- JSON-based character storage
- CSV-based SRD data access
- External API integration for enrichment (dnd5eapi)

## Development

The codebase is organized into layers with clear separation of concerns:

1. **Domain** - Define entities, value objects, and interfaces
2. **Application** - Create use cases that implement business workflows
3. **Infrastructure** - Implement interfaces with concrete technologies
4. **Presentation** - Build user-facing handlers

All components are wired together via dependency injection in `internal/infrastructure/di/container.go`.
