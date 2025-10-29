#!/bin/bash
# SonarQube Analysis Runner
# This script runs tests with coverage, then runs SonarQube analysis

# Load token from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Check if token is set
if [ -z "$SONAR_TOKEN" ]; then
  echo "Error: SONAR_TOKEN not set"
  echo "Please create a .env file with: SONAR_TOKEN=your_token_here"
  exit 1
fi

# Run tests with coverage first
echo "Running tests with coverage..."
mkdir -p coverage
go test ./... -coverprofile=coverage/coverage.out -covermode=count

echo ""
echo "Coverage Summary:"
go tool cover -func=coverage/coverage.out | grep total:

echo ""
echo "Running SonarQube analysis..."
docker run --rm --link sonarqube \
  -v "$(pwd):/usr/src" \
  sonarsource/sonar-scanner-cli \
  -Dsonar.token="${SONAR_TOKEN}"
