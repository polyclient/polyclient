#!/usr/bin/env bash

# PolyClient Documentation Structure Generator
# This script creates the entire documentation directory structure with placeholder files

set -e # Exit on error

lang=$1

if [ -z "$lang" ]; then
    echo "Usage: $0 <lang>"
    exit 1
fi

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
create_file "docs/$lang/index.md" "PolyClient Documentation"

# Getting Started section
create_file "docs/$lang/getting-started/index.md" "Getting Started"
create_file "docs/$lang/getting-started/what-is-polyclient.md" "What is PolyClient"
create_file "docs/$lang/getting-started/key-features.md" "Key Features and Benefits"

# Installation subsection
create_file "docs/$lang/getting-started/installation/index.md" "Installation Guide"
create_file "docs/$lang/getting-started/installation/system-requirements.md" "System Requirements"
create_file "docs/$lang/getting-started/installation/binary-installation.md" "Binary Installation"
create_file "docs/$lang/getting-started/installation/package-managers.md" "Package Manager Installation"
create_file "docs/$lang/getting-started/installation/docker-installation.md" "Docker/Podman Installation"
create_file "docs/$lang/getting-started/installation/build-from-source.md" "Building from Source"

# Quick Start subsection
create_file "docs/$lang/getting-started/quick-start/index.md" "Quick Start Guide"
create_file "docs/$lang/getting-started/quick-start/first-connection.md" "Your First Connection"
create_file "docs/$lang/getting-started/quick-start/basic-operations.md" "Basic Operations"
create_file "docs/$lang/getting-started/quick-start/gui-navigation.md" "GUI Navigation"

# Configuration subsection
create_file "docs/$lang/getting-started/configuration/index.md" "Configuration Basics"
create_file "docs/$lang/getting-started/configuration/config-file-structure.md" "Configuration File Structure"
create_file "docs/$lang/getting-started/configuration/environment-variables.md" "Environment Variables"
create_file "docs/$lang/getting-started/configuration/cli-flags.md" "CLI Flags"

# Core Components section
create_file "docs/$lang/core-components/index.md" "Core Components"

# Architecture subsection
create_file "docs/$lang/core-components/architecture/index.md" "Architecture Overview"
create_file "docs/$lang/core-components/architecture/component-diagram.md" "Component Diagram"
create_file "docs/$lang/core-components/architecture/data-flow.md" "Data Flow"

# CLI subsection
create_file "docs/$lang/core-components/cli/index.md" "Command Line Interface"
create_file "docs/$lang/core-components/cli/command-reference.md" "Command Reference"
create_file "docs/$lang/core-components/cli/database-management.md" "Database Management Commands"
create_file "docs/$lang/core-components/cli/plugin-management.md" "Plugin Management Commands"
create_file "docs/$lang/core-components/cli/configuration-options.md" "Configuration Options"

# GUI subsection
create_file "docs/$lang/core-components/gui/index.md" "Graphical User Interface"
create_file "docs/$lang/core-components/gui/layout-navigation.md" "Layout and Navigation"
create_file "docs/$lang/core-components/gui/workspace-management.md" "Workspace Management"
create_file "docs/$lang/core-components/gui/customization.md" "Customization Options"
create_file "docs/$lang/core-components/gui/keyboard-shortcuts.md" "Keyboard Shortcuts"

# API subsection
create_file "docs/$lang/core-components/api/index.md" "API"
create_file "docs/$lang/core-components/api/api-reference.md" "API Reference"
create_file "docs/$lang/core-components/api/authentication-security.md" "Authentication and Security"
create_file "docs/$lang/core-components/api/integration-examples.md" "Integration Examples"

# Database Management section
create_file "docs/$lang/database-management/index.md" "Database Management"

# Connection Management subsection
create_file "docs/$lang/database-management/connection-management/index.md" "Connection Management"
create_file "docs/$lang/database-management/connection-management/creating-connections.md" "Creating New Connections"
create_file "docs/$lang/database-management/connection-management/connection-settings.md" "Connection Settings"
create_file "docs/$lang/database-management/connection-management/connection-pooling.md" "Connection Pooling"
create_file "docs/$lang/database-management/connection-management/secure-connections.md" "Secure Connections"

# Query Interface subsection
create_file "docs/$lang/database-management/query-interface/index.md" "Query Interface"
create_file "docs/$lang/database-management/query-interface/query-editor.md" "Query Editor Features"
create_file "docs/$lang/database-management/query-interface/visual-query-builder.md" "Visual Query Builder"
create_file "docs/$lang/database-management/query-interface/query-history.md" "Query History and Favorites"
create_file "docs/$lang/database-management/query-interface/query-optimization.md" "Query Optimization"

# Schema Management subsection
create_file "docs/$lang/database-management/schema-management/index.md" "Schema Management"
create_file "docs/$lang/database-management/schema-management/browsing-objects.md" "Browsing Database Objects"
create_file "docs/$lang/database-management/schema-management/creating-modifying-objects.md" "Creating and Modifying Objects"
create_file "docs/$lang/database-management/schema-management/entity-relationship-diagrams.md" "Entity Relationship Diagrams"
create_file "docs/$lang/database-management/schema-management/schema-comparison.md" "Schema Comparison and Synchronization"

# Data Management subsection
create_file "docs/$lang/database-management/data-management/index.md" "Data Management"
create_file "docs/$lang/database-management/data-management/data-browsing-editing.md" "Data Browsing and Editing"
create_file "docs/$lang/database-management/data-management/import-export.md" "Import/Export Operations"
create_file "docs/$lang/database-management/data-management/batch-operations.md" "Batch Operations"
create_file "docs/$lang/database-management/data-management/data-masking.md" "Data Masking"

# Databases section
create_file "docs/$lang/databases/index.md" "Supported Databases"

# Relational Databases subsection
create_file "docs/$lang/databases/relational/index.md" "Relational Databases"
create_file "docs/$lang/databases/relational/postgresql.md" "PostgreSQL"
create_file "docs/$lang/databases/relational/mysql.md" "MySQL"
create_file "docs/$lang/databases/relational/sqlite.md" "SQLite"
create_file "docs/$lang/databases/relational/sqlserver.md" "SQL Server"
create_file "docs/$lang/databases/relational/oracle.md" "Oracle Database"
create_file "docs/$lang/databases/relational/mariadb.md" "MariaDB"
create_file "docs/$lang/databases/relational/firebird.md" "Firebird"
create_file "docs/$lang/databases/relational/libsql.md" "LibSQL"

# NoSQL Databases subsection
create_file "docs/$lang/databases/nosql/index.md" "NoSQL Databases"
create_file "docs/$lang/databases/nosql/mongodb.md" "MongoDB"
create_file "docs/$lang/databases/nosql/redis.md" "Redis"
create_file "docs/$lang/databases/nosql/cassandra.md" "Cassandra"

# Distributed Databases subsection
create_file "docs/$lang/databases/distributed/index.md" "Distributed Databases"
create_file "docs/$lang/databases/distributed/cockroachdb.md" "CockroachDB"

# Analytical Databases subsection
create_file "docs/$lang/databases/analytical/index.md" "Analytical Databases"
create_file "docs/$lang/databases/analytical/duckdb.md" "DuckDB"
create_file "docs/$lang/databases/analytical/clickhouse.md" "ClickHouse"

create_file "docs/$lang/databases/database-specific-features.md" "Database-Specific Features"

# Plugin System section
create_file "docs/$lang/plugin-system/index.md" "Plugin System"

# Plugin Architecture subsection
create_file "docs/$lang/plugin-system/architecture/index.md" "Plugin Architecture"
create_file "docs/$lang/plugin-system/architecture/wasm-runtime.md" "WebAssembly Runtime"
create_file "docs/$lang/plugin-system/architecture/plugin-interfaces.md" "Plugin Interfaces"
create_file "docs/$lang/plugin-system/architecture/security-model.md" "Security Model"

# Using Plugins subsection
create_file "docs/$lang/plugin-system/using-plugins/index.md" "Using Plugins"
create_file "docs/$lang/plugin-system/using-plugins/discovering-plugins.md" "Discovering Plugins"
create_file "docs/$lang/plugin-system/using-plugins/installing-plugins.md" "Installing Plugins"
create_file "docs/$lang/plugin-system/using-plugins/configuring-plugins.md" "Configuring Plugins"
create_file "docs/$lang/plugin-system/using-plugins/updates-versioning.md" "Updates and Versioning"

# Developing Plugins subsection
create_file "docs/$lang/plugin-system/developing-plugins/index.md" "Developing Plugins"
create_file "docs/$lang/plugin-system/developing-plugins/development-lifecycle.md" "Plugin Development Lifecycle"
create_file "docs/$lang/plugin-system/developing-plugins/sql-plugin-development.md" "SQL Database Plugin Development"
create_file "docs/$lang/plugin-system/developing-plugins/nosql-plugin-development.md" "NoSQL Database Plugin Development"
create_file "docs/$lang/plugin-system/developing-plugins/theme-plugin-development.md" "Theme Plugin Development"
create_file "docs/$lang/plugin-system/developing-plugins/ai-plugin-development.md" "AI Integration Plugin Development"
create_file "docs/$lang/plugin-system/developing-plugins/automation-plugin-development.md" "Automation Plugin Development"
create_file "docs/$lang/plugin-system/developing-plugins/publishing-plugins.md" "Publishing Plugins"

# Advanced Features section
create_file "docs/$lang/advanced-features/index.md" "Advanced Features"

# Data Pipelines subsection
create_file "docs/$lang/advanced-features/data-pipelines/index.md" "Data Pipelines"
create_file "docs/$lang/advanced-features/data-pipelines/export-formats.md" "Export Formats"
create_file "docs/$lang/advanced-features/data-pipelines/import-formats.md" "Import Formats"
create_file "docs/$lang/advanced-features/data-pipelines/custom-transformations.md" "Custom Data Transformations"

# Backup and Recovery subsection
create_file "docs/$lang/advanced-features/backup-recovery/index.md" "Backup and Recovery"
create_file "docs/$lang/advanced-features/backup-recovery/backup-strategies.md" "Backup Strategies"
create_file "docs/$lang/advanced-features/backup-recovery/scheduled-backups.md" "Scheduled Backups"
create_file "docs/$lang/advanced-features/backup-recovery/point-in-time-recovery.md" "Point-in-Time Recovery"
create_file "docs/$lang/advanced-features/backup-recovery/disaster-recovery.md" "Disaster Recovery Planning"

# Performance Monitoring subsection
create_file "docs/$lang/advanced-features/performance-monitoring/index.md" "Performance Monitoring"
create_file "docs/$lang/advanced-features/performance-monitoring/database-performance.md" "Database Performance Metrics"
create_file "docs/$lang/advanced-features/performance-monitoring/query-performance.md" "Query Performance Analysis"
create_file "docs/$lang/advanced-features/performance-monitoring/performance-dashboards.md" "Performance Dashboards"
create_file "docs/$lang/advanced-features/performance-monitoring/alerts-notifications.md" "Alerts and Notifications"

# Security subsection
create_file "docs/$lang/advanced-features/security/index.md" "Security and Access Control"
create_file "docs/$lang/advanced-features/security/user-management.md" "User Management"
create_file "docs/$lang/advanced-features/security/role-based-access.md" "Role-Based Access Control"
create_file "docs/$lang/advanced-features/security/data-masking.md" "Dynamic Data Masking"
create_file "docs/$lang/advanced-features/security/row-level-security.md" "Row-Level Security"
create_file "docs/$lang/advanced-features/security/audit-logging.md" "Audit Logging"

# Task Automation subsection
create_file "docs/$lang/advanced-features/task-automation/index.md" "Task Automation"
create_file "docs/$lang/advanced-features/task-automation/job-scheduling.md" "Job Scheduling"
create_file "docs/$lang/advanced-features/task-automation/triggers-events.md" "Triggers and Events"
create_file "docs/$lang/advanced-features/task-automation/workflow-automation.md" "Workflow Automation"

# Integration section
create_file "docs/$lang/integration/index.md" "Integration and Extensibility"

# API Integration subsection
create_file "docs/$lang/integration/api-integration/index.md" "API Integration"
create_file "docs/$lang/integration/api-integration/rest-api.md" "REST API Reference"
create_file "docs/$lang/integration/api-integration/graphql-api.md" "GraphQL API Reference"
create_file "docs/$lang/integration/api-integration/websocket-api.md" "WebSocket API Reference"

# Service Integration subsection
create_file "docs/$lang/integration/service-integration/index.md" "Service Integration"
create_file "docs/$lang/integration/service-integration/cicd-integration.md" "CI/CD Pipeline Integration"
create_file "docs/$lang/integration/service-integration/version-control.md" "Version Control Integration"
create_file "docs/$lang/integration/service-integration/third-party-tools.md" "Third-Party Tool Integration"

# Custom Extensions subsection
create_file "docs/$lang/integration/custom-extensions/index.md" "Custom Extensions"
create_file "docs/$lang/integration/custom-extensions/extension-points.md" "Extension Points"
create_file "docs/$lang/integration/custom-extensions/custom-dashboards.md" "Building Custom Dashboards"
create_file "docs/$lang/integration/custom-extensions/data-visualizations.md" "Custom Data Visualizations"

# Developer Reference section
create_file "docs/$lang/developer-reference/index.md" "Developer Reference"

# Architecture Deep Dive subsection
create_file "docs/$lang/developer-reference/architecture-deep-dive/index.md" "Architecture Deep Dive"
create_file "docs/$lang/developer-reference/architecture-deep-dive/core-components.md" "Core Components"
create_file "docs/$lang/developer-reference/architecture-deep-dive/data-flow.md" "Data Flow"
create_file "docs/$lang/developer-reference/architecture-deep-dive/plugin-system.md" "Plugin System"

# API Reference subsection
create_file "docs/$lang/developer-reference/api-reference/index.md" "API Reference"
create_file "docs/$lang/developer-reference/api-reference/internal-apis.md" "Internal APIs"
create_file "docs/$lang/developer-reference/api-reference/public-apis.md" "Public APIs"
create_file "docs/$lang/developer-reference/api-reference/plugin-apis.md" "Plugin APIs"

# WIT Reference subsection
create_file "docs/$lang/developer-reference/wit-reference/index.md" "WebAssembly Interface Types (WIT) Reference"
create_file "docs/$lang/developer-reference/wit-reference/polyclient-sql-wit.md" "polyclient-sql.wit"
create_file "docs/$lang/developer-reference/wit-reference/custom-wit-definitions.md" "Custom WIT Definitions"

# Building from Source subsection
create_file "docs/$lang/developer-reference/building-from-source/index.md" "Building from Source"
create_file "docs/$lang/developer-reference/building-from-source/dev-environment-setup.md" "Development Environment Setup"
create_file "docs/$lang/developer-reference/building-from-source/build-process.md" "Build Process"
create_file "docs/$lang/developer-reference/building-from-source/testing-framework.md" "Testing Framework"

# Contributing section
create_file "docs/$lang/contributing/index.md" "Contributing to PolyClient"

# Guidelines subsection
create_file "docs/$lang/contributing/guidelines/index.md" "Contribution Guidelines"
create_file "docs/$lang/contributing/guidelines/code-of-conduct.md" "Code of Conduct"
create_file "docs/$lang/contributing/guidelines/development-workflow.md" "Development Workflow"
create_file "docs/$lang/contributing/guidelines/pull-request-process.md" "Pull Request Process"

# Development Environment subsection
create_file "docs/$lang/contributing/development-environment/index.md" "Development Environment"
create_file "docs/$lang/contributing/development-environment/local-setup.md" "Local Development Setup"
create_file "docs/$lang/contributing/development-environment/remote-setup.md" "Remote Development Setup"
create_file "docs/$lang/contributing/development-environment/development-tools.md" "Development Tools"

# Testing subsection
create_file "docs/$lang/contributing/testing/index.md" "Testing Guidelines"
create_file "docs/$lang/contributing/testing/unit-testing.md" "Unit Testing"
create_file "docs/$lang/contributing/testing/integration-testing.md" "Integration Testing"
create_file "docs/$lang/contributing/testing/end-to-end-testing.md" "End-to-End Testing"
create_file "docs/$lang/contributing/testing/performance-testing.md" "Performance Testing"

# Documentation subsection
create_file "docs/$lang/contributing/documentation/index.md" "Documentation Guidelines"
create_file "docs/$lang/contributing/documentation/code-documentation.md" "Code Documentation"
create_file "docs/$lang/contributing/documentation/api-documentation.md" "API Documentation"
create_file "docs/$lang/contributing/documentation/user-documentation.md" "User Documentation"

# Troubleshooting section
create_file "docs/$lang/troubleshooting/index.md" "Troubleshooting and Support"
create_file "docs/$lang/troubleshooting/common-issues.md" "Common Issues and Solutions"
create_file "docs/$lang/troubleshooting/logging-diagnostics.md" "Logging and Diagnostics"
create_file "docs/$lang/troubleshooting/community-support.md" "Community Support"
create_file "docs/$lang/troubleshooting/bug-reporting.md" "Bug Reporting"

# Benchmarks section
create_file "docs/$lang/benchmarks/index.md" "Performance Benchmarks"
create_file "docs/$lang/benchmarks/methodology.md" "Benchmark Methodology"
create_file "docs/$lang/benchmarks/connection-performance.md" "Database Connection Performance"
create_file "docs/$lang/benchmarks/query-performance.md" "Query Performance"
create_file "docs/$lang/benchmarks/gui-performance.md" "GUI Performance"
create_file "docs/$lang/benchmarks/plugin-performance.md" "Plugin Performance"
create_file "docs/$lang/benchmarks/comparative-benchmarks.md" "Comparative Benchmarks"

# Releases section
create_file "docs/$lang/releases/index.md" "Release Notes and Roadmap"
create_file "docs/$lang/releases/version-history.md" "Version History"
create_file "docs/$lang/releases/current-release.md" "Current Release Features"
create_file "docs/$lang/releases/known-issues.md" "Known Issues"
create_file "docs/$lang/releases/upcoming-features.md" "Upcoming Features"
create_file "docs/$lang/releases/roadmap.md" "Long-term Roadmap"

# Legal section
create_file "docs/$lang/legal/index.md" "Legal Information"

# License subsection
create_file "docs/$lang/legal/license/index.md" "License Information"
create_file "docs/$lang/legal/license/gpl-license.md" "GPL-3.0-or-later with Plugin Exception"
create_file "docs/$lang/legal/license/apache-license.md" "Plugin Apache-2.0 License"
create_file "docs/$lang/legal/license/third-party-licenses.md" "Third-Party Licenses"

create_file "docs/$lang/legal/privacy-policy.md" "Privacy Policy"
create_file "docs/$lang/legal/security-policy.md" "Security Policy"

# Appendices section
create_file "docs/$lang/appendices/index.md" "Appendices"
create_file "docs/$lang/appendices/command-reference.md" "Command Reference"
create_file "docs/$lang/appendices/configuration-reference.md" "Configuration Reference"
create_file "docs/$lang/appendices/api-reference.md" "API Reference"
create_file "docs/$lang/appendices/environment-variables.md" "Environment Variables"
create_file "docs/$lang/appendices/glossary.md" "Glossary"
create_file "docs/$lang/appendices/resource-links.md" "Resource Links"

echo "Documentation structure created successfully!"
echo "Total files created: $(find docs -type f | wc -l)"
echo "Total directories created: $(find docs -type d | wc -l)"
