#!/bin/bash
# SonarQube Analysis Runner
# This script runs SonarQube analysis with the token from .env file

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

docker run --rm --link sonarqube \
  -v "$(pwd):/usr/src" \
  sonarsource/sonar-scanner-cli \
  -Dsonar.token="${SONAR_TOKEN}"
