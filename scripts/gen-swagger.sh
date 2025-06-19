#!/bin/sh

# Generate swagger.json
swag init --generalInfo internal/interfaces/http/v1/doc.go --output ./docs/v1 --instanceName v1
