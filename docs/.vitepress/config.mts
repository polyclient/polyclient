import { DefaultTheme, defineConfig } from "vitepress";
import { resolve } from "path";
import { readFileSync } from "fs";

const versionFilePath = resolve(__dirname, "../../version.txt");
const version = readFileSync(versionFilePath, "utf-8").trim();

const sidebarGuide = [
	{
		text: "Introduction",
		link: "/guide/",
	},
	{
		text: "Installation",
		items: [
			{ text: "System requirements", link: "/guide/installation/system-requirements" },
			{ text: "Linux", link: "/guide/installation/linux" },
			{ text: "macOS", link: "/guide/installation/macos" },
			{ text: "Windows", link: "/guide/installation/windows" },
			{ text: "Docker", link: "/guide/installation/docker" },
		],
	},
	{
		text: "Quick start",
		items: [
			{ text: "First connection", link: "/guide/quick-start/first-connection" },
			{ text: "Basic operations", link: "/guide/quick-start/basic-operations" },
		],
	},
	{
		text: "Configuration",
		items: [
			{ text: "Overview", link: "/guide/configuration/" },
			{ text: "Global settings", link: "/guide/configuration/global-settings" },
			{ text: "Database connections", link: "/guide/configuration/database-connections" },
			{ text: "Authentication", link: "/guide/configuration/authentication" },
			{ text: "Plugins", link: "/guide/configuration/plugins" },
			{ text: "CLI", link: "/guide/configuration/cli" },
			{ text: "GUI", link: "/guide/configuration/gui" },
			{ text: "Environment variables", link: "/guide/configuration/environment-variables" },
			{ text: "Configuration file", link: "/guide/configuration/configuration-file" },
			{ text: "Troubleshooting", link: "/guide/configuration/troubleshooting" },
		],
	},
] satisfies DefaultTheme.NavItem[];

const sidebarCLI = [
	{ text: "Overview", link: "/cli/" },
	{ text: "Command reference", link: "/cli/command-reference" },
	{ text: "Database management", link: "/cli/database-management" },
	{ text: "Plugin management", link: "/cli/plugin-management" },
	{
		text: "Configuration",
		items: [
			{ text: "CLI flags", link: "/cli/configuration/cli-flags" },
			{ text: "Config file structure", link: "/cli/configuration/config-file-structure" },
			{ text: "Environment variables", link: "/cli/configuration/environment-variables" },
		],
	},
	{ text: "Troubleshooting", link: "/cli/troubleshooting" },
] satisfies DefaultTheme.NavItem[];

const sidebarGUI = [
	{ text: "Overview", link: "/gui/" },
	{ text: "Layout & navigation", link: "/gui/layout-navigation" },
	{ text: "Workspace management", link: "/gui/workspace-management" },
	{ text: "Keyboard shortcuts", link: "/gui/keyboard-shortcuts" },
	{ text: "Customization", link: "/gui/customization" },
	{
		text: "Query interface",
		items: [
			{ text: "Query editor", link: "/gui/query-interface/query-editor" },
			{ text: "Visual query builder", link: "/gui/query-interface/visual-query-builder" },
			{ text: "Query history", link: "/gui/query-interface/query-history" },
			{ text: "Query optimization", link: "/gui/query-interface/query-optimization" },
		],
	},
	{
		text: "Schema management",
		items: [
			{ text: "Browsing objects", link: "/gui/schema-management/browsing-objects" },
			{ text: "Creating & modifying objects", link: "/gui/schema-management/creating-modifying-objects" },
			{ text: "Entity relationship diagrams", link: "/gui/schema-management/entity-relationship-diagrams" },
			{ text: "Schema comparison", link: "/gui/schema-management/schema-comparison" },
		],
	},
	{ text: "Troubleshooting", link: "/gui/troubleshooting" },
] satisfies DefaultTheme.NavItem[];

const sidebarAPI = [
	{ text: "Overview", link: "/api/" },
	{ text: "API reference", link: "/api/api-reference" },
	{ text: "Authentication & security", link: "/api/authentication-security" },
	{
		text: "Endpoints",
		items: [
			{ text: "Connection", link: "/api/endpoints/connection" },
			{ text: "Query", link: "/api/endpoints/query" },
			{ text: "Schema", link: "/api/endpoints/schema" },
			{ text: "User", link: "/api/endpoints/user" },
		],
	},
	{ text: "Integration examples", link: "/api/integration-examples" },
	{ text: "Troubleshooting", link: "/api/troubleshooting" },
] satisfies DefaultTheme.NavItem[];

const sidebarPlugins = [
	{ text: "Overview", link: "/plugins/" },
	{
		text: "Using plugins",
		items: [
			{ text: "Discovering plugins", link: "/plugins/using-plugins/discovering-plugins" },
			{ text: "Installing plugins", link: "/plugins/using-plugins/installing-plugins" },
			{ text: "Configuring plugins", link: "/plugins/using-plugins/configuring-plugins" },
			{ text: "Updates & versioning", link: "/plugins/using-plugins/updates-versioning" },
		],
	},
	{
		text: "Developing plugins",
		items: [
			{ text: "SQL plugin development", link: "/plugins/developing-plugins/sql-plugin-development" },
			{ text: "NoSQL plugin development", link: "/plugins/developing-plugins/nosql-plugin-development" },
			{ text: "Automation plugin development", link: "/plugins/developing-plugins/automation-plugin-development" },
			{ text: "Theme plugin development", link: "/plugins/developing-plugins/theme-plugin-development" },
			{ text: "Publishing plugins", link: "/plugins/developing-plugins/publishing-plugins" },
		],
	},
	{
		text: "Architecture",
		items: [
			{ text: "Plugin interfaces", link: "/plugins/architecture/plugin-interfaces" },
			{ text: "Wasm runtime", link: "/plugins/architecture/wasm-runtime" },
			{ text: "Security model", link: "/plugins/architecture/security-model" },
		],
	},
] satisfies DefaultTheme.NavItem[];

// https://vitepress.dev/reference/site-config
export default defineConfig({
	lang: "en-US",
	title: "PolyClient",
	description: "Unified database management platform.",
	lastUpdated: true,
	cleanUrls: true,
	ignoreDeadLinks: "localhostLinks",
	themeConfig: {
		nav: [
			{ text: "Guide", link: "/guide/" },
			{ text: "CLI", link: "/cli" },
			{ text: "GUI", link: "/gui" },
			{ text: "API", link: "/api" },
			{ text: "Plugins", link: "/plugins" },
			{
				text: `v${version}`,
				items: [
					{
						text: "Releases",
						link: "https://github.com/polyclient/polyclient/releases",
					},
					{
						text: "Contributing",
						link: "https://github.com/polyclient/polyclient/blob/main/CONTRIBUTING.md",
					},
				],
			},
		],
		search: {
			provider: "local",
		},
		sidebar: {
			"/guide/": sidebarGuide,
			"/cli/": sidebarCLI,
			"/gui/": sidebarGUI,
			"/api/": sidebarAPI,
			"/plugins/": sidebarPlugins,
		},
		editLink: {
			pattern: "https://github.com/polyclient/polyclient/blob/master/docs/:path",
			text: "Suggest changes to this page",
		},
		socialLinks: [{ icon: "github", link: "https://github.com/polyclient/polyclient" }],
	},
});
