#!/bin/bash

# Verifica se o git está limpo
if [[ -n $(git status -s) ]]; then
  echo "Error: Working directory is not clean. Please commit or stash changes."
  exit 1
fi

# Pega a última tag
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null)

if [ -z "$LAST_TAG" ]; then
  # Se não houver tags, começa com v0.1.0
  NEW_TAG="v0.1.0"
else
  # Incrementa a versão
  VERSION=$(echo $LAST_TAG | sed 's/v//')
  MAJOR=$(echo $VERSION | cut -d. -f1)
  MINOR=$(echo $VERSION | cut -d. -f2)
  PATCH=$(echo $VERSION | cut -d. -f3)
  
  # Incrementa o patch
  NEW_PATCH=$((PATCH + 1))
  NEW_TAG="v${MAJOR}.${MINOR}.${NEW_PATCH}"
fi

# Cria a tag
git tag -a $NEW_TAG -m "Release $NEW_TAG"
git push origin $NEW_TAG

echo "Created and pushed tag $NEW_TAG" 