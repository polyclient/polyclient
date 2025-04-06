import type { Component } from 'svelte';

export const iconsNames = {
	railDatabase: 'rail:database',
	railExplorer: 'rail:explorer',
	railHistory: 'rail:history',
	railPlugins: 'rail:plugins',
	railAssistant: 'rail:assistant',
	railSettings: 'rail:settings',

	dbSchema: 'db:schema',
	dbTable: 'db:table',
	dbColumn: 'db:column',
	dbColumnPk: 'db:column-pk',
	dbColumnFk: 'db:column-fk',
	dbView: 'db:view',
	dbFunction: 'db:function',
	dbProcedure: 'db:procedure',
	dbTrigger: 'db:trigger',
	dbConnection: 'db:connection',
	dbFolder: 'db:folder',

	uiChevronRight: 'ui:chevron-right',
	uiChevronLeft: 'ui:chevron-left',
	uiChevronUp: 'ui:chevron-up',
	uiChevronDown: 'ui:chevron-down',
	uiArrowRight: 'ui:arrow-right',
	uiArrowLeft: 'ui:arrow-left',
	uiArrowUp: 'ui:arrow-up',
	uiArrowDown: 'ui:arrow-down',
	uiNext: 'ui:next',
	uiPrevious: 'ui:previous',
	uiPlus: 'ui:plus',
	uiMinus: 'ui:minus',
	uiInfo: 'ui:info',
	uiError: 'ui:error',
	uiWarning: 'ui:warning',
	uiSuccess: 'ui:success',
} as const;

export type IconName = (typeof iconsNames)[keyof typeof iconsNames];

export type IconComponent = Component;

export type IconSet = Record<IconName, IconComponent | undefined>;
