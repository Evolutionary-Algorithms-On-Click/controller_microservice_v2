#!/bin/sh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- VALIDATION ---
if [ -z "$1" ]; then
  echo "Error: No resource name provided."
  echo "Usage: ./scripts/make_endpoint.sh <resource_name>"
  echo "Example: ./scripts/make_endpoint.sh user"
  exit 1
fi

# --- VARIABLE SETUP ---
RESOURCE_NAME=$1
# Example: user -> users
RESOURCE_NAME_PLURAL="${RESOURCE_NAME}s"
# Example: user -> User
RESOURCE_NAME_PASCAL=$(echo "$RESOURCE_NAME" | sed -e "s/\b\(.\)/\u\1/g")
# Example: user -> UserModule
MODULE_PASCAL="${RESOURCE_NAME_PASCAL}Module"
# Example: user -> UserController
CONTROLLER_PASCAL="${RESOURCE_NAME_PASCAL}Controller"
# Example: user -> UserRepository
REPOSITORY_PASCAL="${RESOURCE_NAME_PASCAL}Repository"

# File paths
CONTROLLER_FILE="controllers/${RESOURCE_NAME}_controller.go"
MODULE_FILE="modules/${RESOURCE_NAME}_module.go"
REPOSITORY_DIR="db/repository"
REPOSITORY_FILE="${REPOSITORY_DIR}/${RESOURCE_NAME}_repository.go"
ROUTES_FILE="routes/api.go"

echo "--- Generating files for resource: $RESOURCE_NAME ---"

# --- PRE-FLIGHT CHECKS ---
if [ ! -f "$ROUTES_FILE" ]; then
    echo "Error: Main routes file not found at '$ROUTES_FILE'. Cannot proceed."
    exit 1
fi

if [ -f "$CONTROLLER_FILE" ] || [ -f "$MODULE_FILE" ]; then
    echo "Error: Files for resource '$RESOURCE_NAME' already exist. Aborting to prevent overwrite."
    exit 1
fi

# --- CREATE DIRECTORIES ---
mkdir -p "$REPOSITORY_DIR"
echo "Directory '$REPOSITORY_DIR' ensured."

# --- GENERATE CONTROLLER FILE ---
cat <<EOF > "$CONTROLLER_FILE"
package controllers

import (
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/modules"
	"github.com/rs/zerolog"
)

// ${CONTROLLER_PASCAL} holds the dependencies for the ${RESOURCE_NAME} handlers.
type ${CONTROLLER_PASCAL} struct {
	Module *modules.${MODULE_PASCAL}
	Logger zerolog.Logger
}

// New${CONTROLLER_PASCAL} creates and returns a new ${CONTROLLER_PASCAL}.
func New${CONTROLLER_PASCAL}(module *modules.${MODULE_PASCAL}, logger zerolog.Logger) *${CONTROLLER_PASCAL} {
	return &${CONTROLLER_PASCAL}{
		Module: module,
		Logger: logger,
	}
}

// Create${RESOURCE_NAME_PASCAL}Handler handles POST /api/v1/${RESOURCE_NAME_PLURAL}
func (c *${CONTROLLER_PASCAL}) Create${RESOURCE_NAME_PASCAL}Handler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to create a new ${RESOURCE_NAME}.
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// List${RESOURCE_NAME_PASCAL}sHandler handles GET /api/v1/${RESOURCE_NAME_PLURAL}
func (c *${CONTROLLER_PASCAL}) List${RESOURCE_NAME_PASCAL}sHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to list all ${RESOURCE_NAME_PLURAL}.
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// Get${RESOURCE_NAME_PASCAL}ByIDHandler handles GET /api/v1/${RESOURCE_NAME_PLURAL}/{id}
func (c *${CONTROLLER_PASCAL}) Get${RESOURCE_NAME_PASCAL}ByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// TODO: Implement logic to get a ${RESOURCE_NAME} by its ID.
	http.Error(w, "Not Implemented: get by id "+id, http.StatusNotImplemented)
}

// Update${RESOURCE_NAME_PASCAL}ByIDHandler handles PUT /api/v1/${RESOURCE_NAME_PLURAL}/{id}
func (c *${CONTROLLER_PASCAL}) Update${RESOURCE_NAME_PASCAL}ByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// TODO: Implement logic to update a ${RESOURCE_NAME}.
	http.Error(w, "Not Implemented: update by id "+id, http.StatusNotImplemented)
}

// Delete${RESOURCE_NAME_PASCAL}ByIDHandler handles DELETE /api/v1/${RESOURCE_NAME_PLURAL}/{id}
func (c *${CONTROLLER_PASCAL}) Delete${RESOURCE_NAME_PASCAL}ByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// TODO: Implement logic to delete a ${RESOURCE_NAME}.
	http.Error(w, "Not Implemented: delete by id "+id, http.StatusNotImplemented)
}
EOF
echo "Created $CONTROLLER_FILE"

# --- GENERATE MODULE FILE ---
cat <<EOF > "$MODULE_FILE"
package modules

import (
	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
)

// ${MODULE_PASCAL} encapsulates the business logic for ${RESOURCE_NAME_PLURAL}.
type ${MODULE_PASCAL} struct {
	Repo repository.${REPOSITORY_PASCAL}
}

// New${MODULE_PASCAL} creates and returns a new ${MODULE_PASCAL}.
func New${MODULE_PASCAL}(repo repository.${REPOSITORY_PASCAL}) *${MODULE_PASCAL} {
	return &${MODULE_PASCAL}{
		Repo: repo,
	}
}

// TODO: Implement business logic methods for ${MODULE_PASCAL}
// Example:
// func (m *${MODULE_PASCAL}) Get${RESOURCE_NAME_PASCAL}ByID(ctx context.Context, id string) (*models.${RESOURCE_NAME_PASCAL}, error) {
// 	 return m.Repo.GetByID(ctx, id)
// }
EOF
echo "Created $MODULE_FILE"

# --- GENERATE REPOSITORY INTERFACE FILE ---
cat <<EOF > "$REPOSITORY_FILE"
package repository

import (
	"context"
	// "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
)

// ${REPOSITORY_PASCAL} defines the data access methods for a ${RESOURCE_NAME}.
type ${REPOSITORY_PASCAL} interface {
	// TODO: Define repository methods.
	// Example:
	// GetByID(ctx context.Context, id string) (*models.${RESOURCE_NAME_PASCAL}, error)
	// Create(ctx context.Context, ${RESOURCE_NAME} *models.${RESOURCE_NAME_PASCAL}) error
}

// TODO: Create a new file in this directory for the concrete implementation, e.g., postgres_${RESOURCE_NAME}.go
EOF
echo "Created $REPOSITORY_FILE"

# --- GENERATE TEST FILES ---
cat <<EOF > "${CONTROLLER_FILE%.go}_test.go"
package controllers_test

import "testing"

func Test${CONTROLLER_PASCAL}(t *testing.T) {
	// TODO: Write unit tests for the ${CONTROLLER_PASCAL}
}
EOF
echo "Created ${CONTROLLER_FILE%.go}_test.go"

cat <<EOF > "${MODULE_FILE%.go}_test.go"
package modules_test

import "testing"

func Test${MODULE_PASCAL}(t *testing.T) {
	// TODO: Write unit tests for the ${MODULE_PASCAL}
}
EOF
echo "Created ${MODULE_FILE%.go}_test.go"

# --- APPEND ROUTES TO routes/api.go ---
ROUTES_BLOCK="
	// ${RESOURCE_NAME_PASCAL} Routes
	mux.HandleFunc(\"POST /api/v1/${RESOURCE_NAME_PLURAL}\", ${RESOURCE_NAME}Controller.Create${RESOURCE_NAME_PASCAL}Handler)
	mux.HandleFunc(\"GET /api/v1/${RESOURCE_NAME_PLURAL}\", ${RESOURCE_NAME}Controller.List${RESOURCE_NAME_PASCAL}sHandler)
	mux.HandleFunc(\"GET /api/v1/${RESOURCE_NAME_PLURAL}/{id}\", ${RESOURCE_NAME}Controller.Get${RESOURCE_NAME_PASCAL}ByIDHandler)
	mux.HandleFunc(\"PUT /api/v1/${RESOURCE_NAME_PLURAL}/{id}\", ${RESOURCE_NAME}Controller.Update${RESOURCE_NAME_PASCAL}ByIDHandler)
	mux.HandleFunc(\"DELETE /api/v1/${RESOURCE_NAME_PLURAL}/{id}\", ${RESOURCE_NAME}Controller.Delete${RESOURCE_NAME_PASCAL}ByIDHandler)
"

CONTROLLER_INSTANTIATION_BLOCK="
	${RESOURCE_NAME}Module := modules.New${MODULE_PASCAL}(nil) // TODO: Provide real repository implementation
	${RESOURCE_NAME}Controller := controllers.New${CONTROLLER_PASCAL}(${RESOURCE_NAME}Module, *pkg.Logger)
"

# Use a temporary file to make the awk edits safer
TMP_ROUTES_FILE=$(mktemp)

# Insert the controller instantiation before the problemController instantiation
awk -v block="$CONTROLLER_INSTANTIATION_BLOCK" '/problemModule :=/ {print block} 1' "$ROUTES_FILE" > "$TMP_ROUTES_FILE"
mv "$TMP_ROUTES_FILE" "$ROUTES_FILE"

# Re-create a new temp file for the next operation
TMP_ROUTES_FILE=$(mktemp)

# Insert the routes block before the Kernel Routes
awk -v block="$ROUTES_BLOCK" '/\/\/ Kernel Routes/ {print block} 1' "$ROUTES_FILE" > "$TMP_ROUTES_FILE"
mv "$TMP_ROUTES_FILE" "$ROUTES_FILE"

echo "Updated $ROUTES_FILE with new routes and controller for '$RESOURCE_NAME'."

echo "--- Scaffolding complete! ---"
echo "Next steps:"
echo "1. Create a model for '${RESOURCE_NAME_PASCAL}' in the 'pkg/models' directory."
echo "2. Create the concrete repository implementation (e.g., 'db/repository/postgres_${RESOURCE_NAME}.go')."
echo "3. Update the 'nil' in 'routes/api.go' with the real repository instance."
echo "4. Implement the 'TODO' sections in the newly generated files."