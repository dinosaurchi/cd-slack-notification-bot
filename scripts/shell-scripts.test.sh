#!/bin/sh

echo "Testing get-toml-value.sh"
./scripts/get-toml-value.sh env/ci.toml PostgreSQLDatabase.Host
./scripts/get-toml-value.sh env/beta.toml PostgreSQLDatabase.Host

echo "----------"

echo "Testing get-secret-value.sh"
./scripts/get-secret-value.sh ValensExchangeCredentials eu-central-1 .PostgreSQLDatabase.username

echo "----------"
