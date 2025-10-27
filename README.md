# D&D Character Sheet Manager

A command-line tool for managing D&D 5e characters

## Project Structure

```
DDsheetfinal/
├── cmd/                                    # Application entry points
│   └── ddsheet/
│       └── main.go                        # CLI entry point, routes commands
│
├── internal/                              
│   ├── domain/                            # Domain Layer
│   │   ├── entities/
│   │   │   └── character.go              
│   │   ├── valueobjects/
│   │   │   ├── abilities.go              # Ability modifiers, proficiency bonus
│   │   │   ├── armor.go                  # Armor types and AC info
│   │   │   └── srd_data.go               # Skill/class/race mappings
│   │   ├── repositories/                  # Repository interfaces
│   │   │   ├── character_repository.go
│   │   │   ├── enrichment_repository.go
│   │   │   └── srd_repository.go
│   │   └── services/
│   │       └── character_service.go      # Domain logic (AC, saves, skills)
│   │
│   ├── application/                       # Application Layer
│   │   └── usecases/
│   │       ├── create_character.go
│   │       ├── view_character.go
│   │       ├── list_characters.go
│   │       ├── delete_character.go
│   │       ├── equip_item.go
│   │       ├── unequip_item.go
│   │       ├── learn_spell.go
│   │       ├── prepare_spell.go
│   │       ├── get_learnable_spells.go
│   │       ├── enrich_character.go        # API enrichment
│   │       └── serve_http.go
│   │
│   ├── infrastructure/                    # Infrastructure Layer
│   │   ├── persistence/
│   │   │   ├── json_character_repository.go  
│   │   │   └── csv_srd_repository.go         
│   │   ├── external/
│   │   │   └── dnd5eapi_enrichment_repository.go  # D&D 5e API client
│   │   └── di/
│   │       └── container.go              # Dependency injection container
│   │
│   └── presentation/                      # Presentation Layer (UI)
│       └── cli/
│           ├── create_handler.go
│           ├── view_handler.go
│           ├── list_handler.go
│           ├── delete_handler.go
│           ├── equip_handler.go
│           ├── unequip_handler.go
│           ├── learn_spell_handler.go
│           ├── prepare_spell_handler.go
│           ├── learnable_spells_handler.go
│           ├── enrich_handler.go
│           ├── serve_handler.go           # HTTP server
│           └── help_handler.go
│
├── assets/                                 # Static resources
│   └── srd/
│       ├── 5e-SRD-Equipment.csv
│       └── 5e-SRD-Spells.csv
│
├── data/                                   # Runtime data (gitignored, preserves structure)
│   ├── characters/                        # Character JSON files
│   │   └── .gitkeep
│   ├── enrichments/                       # Enriched character data
│   └── cache/                             # API response cache
│       ├── spells.json
│       └── equipment.json
│
├── web/                                    # Web frontend
│   ├── static/
│   │   ├── index.html                     
│   │   ├── charactersheet.html            
│   │   ├── css/
│   │   │   ├── normalize.css
│   │   │   └── style.css
│   │   └── js/
│   │       ├── list.js
│   │       └── charactersheet.js
│   └── manifest.json
│
├── go.mod
└── README.md
```

## Architecture

This project follows **Onion Architecture** (aka Clean Architecture):

- **Domain Layer** 
- **Application Layer**
- **Infrastructure Layer**
- **Presentation Layer**

## Building

```bash
go build -o ddsheet ./cmd/ddsheet
```

## Usage

```bash
# Create a character
./ddsheet create -name Gandalf -race Human -class Wizard -level 20 -str 10 -dex 10 -con 10 -int 20 -wis 18 -cha 14

# View character details
./ddsheet view -name Gandalf

# List all characters
./ddsheet list

# Delete a character
./ddsheet delete -name Gandalf

# Equip items
./ddsheet equip -name Gandalf -weapon Quarterstaff
./ddsheet equip -name Gandalf -armor "Leather Armor"
./ddsheet equip -name Gandalf -shield Shield

# Unequip items
./ddsheet unequip -name Gandalf -weapon
./ddsheet unequip -name Gandalf -armor
./ddsheet unequip -name Gandalf -shield

# Learn and prepare spells
./ddsheet learn-spell -name Gandalf -spell "Fireball"
./ddsheet prepare-spell -name Gandalf -spell "Fireball"
./ddsheet learnable-spells -name Gandalf

# Enrich character with D&D 5e API data
./ddsheet enrich Gandalf
./ddsheet enrich Gandalf --inplace           # Save enrichment to character file
./ddsheet enrich Gandalf --force             # Re-fetch even if cached
./ddsheet enrich Gandalf --fetch-all         # Pre-cache all spells and equipment

# Start web server
./ddsheet serve                              # Runs on http://localhost:8080
./ddsheet serve -port 3000                   # Custom port

# Show all commands
./ddsheet help
```

## Features

- Character creation with racial bonuses and skill proficiencies
- Equipment management with AC calculation
- Spell learning and preparation (class-specific)
- Derived stat calculation (initiative, passive perception, etc.)
- JSON-based character storage
- CSV-based SRD data access
- External API integration for enrichment (dnd5eapi)

## Development
All components are wired together via dependency injection (I hope)
