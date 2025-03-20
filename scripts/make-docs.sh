#!/usr/bin/env bash

# PolyClient Documentation Structure Generator
# This script creates the entire documentation directory structure with placeholder files

set -e # Exit on error

echo "Creating PolyClient documentation structure..."

# Function to create a markdown file with basic front matter
create_file() {
    local file_path="$1"
    local title="$2"

    # Create directory if it doesn't exist
    mkdir -p "$(dirname "$file_path")"

    # Create the file with front matter
    cat >"$file_path" <<EOF
---
title: "$title"
description: "${title} for PolyClient Database Management Platform"
---

# $title

TODO: Add content for $title
EOF

    echo "Created: $file_path"
}

# Create main index file
create_file "docs/index.md" "PolyClient Documentation"

# Getting Started section
create_file "docs/getting-started/index.md" "Introduction"
create_file "docs/getting-started/what-is-polyclient.md" "What is PolyClient"

# Installation subsection
create_file "docs/getting-started/installation/binary-installation.md" "Binary Installation"
create_file "docs/getting-started/installation/build-from-source.md" "Building from Source"
create_file "docs/getting-started/installation/docker-installation.md" "Docker/Podman Installation"
create_file "docs/getting-started/installation/package-managers.md" "Package Manager Installation"
create_file "docs/getting-started/installation/system-requirements.md" "System Requirements"

# Quick Start subsection
create_file "docs/getting-started/quick-start/first-connection.md" "Your First Connection"
create_file "docs/getting-started/quick-start/basic-operations.md" "Basic Operations"

# GUI section
create_file "docs/gui/index.md" "GUI Overview"
create_file "docs/gui/layout-navigation.md" "Layout and Navigation"
create_file "docs/gui/keyboard-shortcuts.md" "Keyboard Shortcuts"
create_file "docs/gui/customization.md" "Customization Options"
create_file "docs/gui/workspace-management.md" "Workspace Management"
create_file "docs/gui/troubleshooting.md" "GUI Troubleshooting"

# GUI Query Interface subsection
create_file "docs/gui/query-interface/query-editor.md" "Query Editor Features"
create_file "docs/gui/query-interface/query-optimization.md" "Query Optimization"
create_file "docs/gui/query-interface/query-history.md" "Query History and Favorites"
create_file "docs/gui/query-interface/visual-query-builder.md" "Visual Query Builder"

# GUI Schema Management subsection
create_file "docs/gui/schema-management/browsing-objects.md" "Browsing Database Objects"
create_file "docs/gui/schema-management/creating-modifying-objects.md" "Creating and Modifying Objects"
create_file "docs/gui/schema-management/entity-relationship-diagrams.md" "Entity Relationship Diagrams"
create_file "docs/gui/schema-management/schema-comparison.md" "Schema Comparison and Synchronization"

# CLI section
create_file "docs/cli/index.md" "CLI Overview"
create_file "docs/cli/command-reference.md" "Command Reference"
create_file "docs/cli/configuration-options.md" "Configuration Options"
create_file "docs/cli/database-management.md" "Database Management Commands"
create_file "docs/cli/plugin-management.md" "Plugin Management Commands"
create_file "docs/cli/troubleshooting.md" "CLI Troubleshooting"

# CLI Configuration subsection
create_file "docs/cli/configuration/cli-flags.md" "CLI Flags"
create_file "docs/cli/configuration/config-file-structure.md" "Configuration File Structure"
create_file "docs/cli/configuration/environment-variables.md" "Environment Variables"

# API section
create_file "docs/api/index.md" "API Overview"
create_file "docs/api/api-reference.md" "API Reference"
create_file "docs/api/authentication-security.md" "Authentication and Security"
create_file "docs/api/integration-examples.md" "Integration Examples"
create_file "docs/api/troubleshooting.md" "API Troubleshooting"

# API Endpoints subsection
create_file "docs/api/endpoints/connection.md" "Connection Endpoints"
create_file "docs/api/endpoints/query.md" "Query Endpoints"
create_file "docs/api/endpoints/schema.md" "Schema Endpoints"
create_file "docs/api/endpoints/user.md" "User Endpoints"

# Plugins section
create_file "docs/plugins/index.md" "Plugins Overview"

# Plugin Architecture subsection
create_file "docs/plugins/architecture/plugin-interfaces.md" "Plugin Interfaces"
create_file "docs/plugins/architecture/security-model.md" "Security Model"
create_file "docs/plugins/architecture/wasm-runtime.md" "WebAssembly Runtime"

# Developing Plugins subsection
create_file "docs/plugins/developing-plugins/sql-plugin-development.md" "SQL Database Plugin Development"
create_file "docs/plugins/developing-plugins/nosql-plugin-development.md" "NoSQL Database Plugin Development"
create_file "docs/plugins/developing-plugins/ai-plugin-development.md" "AI Integration Plugin Development"
create_file "docs/plugins/developing-plugins/theme-plugin-development.md" "Theme Plugin Development"
create_file "docs/plugins/developing-plugins/automation-plugin-development.md" "Automation Plugin Development"
create_file "docs/plugins/developing-plugins/publishing-plugins.md" "Publishing Plugins"

# Using Plugins subsection
create_file "docs/plugins/using-plugins/discovering-plugins.md" "Discovering Plugins"
create_file "docs/plugins/using-plugins/installing-plugins.md" "Installing Plugins"
create_file "docs/plugins/using-plugins/configuring-plugins.md" "Configuring Plugins"
create_file "docs/plugins/using-plugins/updates-versioning.md" "Updates and Versioning"

# Contribute section
create_file "docs/contribute/index.md" "Contributing Overview"

# Development Environment subsection
create_file "docs/contribute/development-environment/development-tools.md" "Development Tools"
create_file "docs/contribute/development-environment/local-setup.md" "Local Development Setup"
create_file "docs/contribute/development-environment/remote-setup.md" "Remote Development Setup"

# Guidelines subsection
create_file "docs/contribute/guidelines/code-of-conduct.md" "Code of Conduct"
create_file "docs/contribute/guidelines/development-workflow.md" "Development Workflow"
create_file "docs/contribute/guidelines/pull-request-process.md" "Pull Request Process"

# Testing subsection
create_file "docs/contribute/testing/unit-testing.md" "Unit Testing"
create_file "docs/contribute/testing/integration-testing.md" "Integration Testing"
create_file "docs/contribute/testing/end-to-end-testing.md" "End-to-End Testing"
create_file "docs/contribute/testing/performance-testing.md" "Performance Testing"

echo "Documentation structure created successfully!"
echo "Total files created: $(find . -type f -not -path "./docs/*" | wc -l)"
echo "Total directories created: $(find . -type d -not -path "./docs" | wc -l)"
