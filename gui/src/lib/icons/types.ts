import type { Component } from 'svelte';

export const iconsNames = {
	railHistory: 'rail:history',
	railAssistant: 'rail:assistant',
	railPlugins: 'rail:plugins',
	railNotifications: 'rail:notifications',
	railSettings: 'rail:settings',

	dbDatabase: 'db:database',
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

	uiChevronLeft: 'ui:chevron-left',
	uiChevronRight: 'ui:chevron-right',
	uiChevronUp: 'ui:chevron-up',
	uiChevronDown: 'ui:chevron-down',
	uiArrowLeft: 'ui:arrow-left',
	uiArrowRight: 'ui:arrow-right',
	uiArrowUp: 'ui:arrow-up',
	uiArrowDown: 'ui:arrow-down',
	uiPrevious: 'ui:previous',
	uiNext: 'ui:next',
	uiPlus: 'ui:plus',
	uiMinus: 'ui:minus',
	uiInfo: 'ui:info',
	uiError: 'ui:error',
	uiWarning: 'ui:warning',
	uiSuccess: 'ui:success',
} as const;

export type IconName = (typeof iconsNames)[keyof typeof iconsNames];

export type IconSourceComponent = Component;

export type IconSourceSvg = `<svg${string}></svg>`;

export type IconSourceIconify = string;

export type IconSource = IconSourceComponent | IconSourceSvg | IconSourceIconify;

export type IconSet = Record<IconName, IconSource | undefined>;

/**
 * Type guard to check if an icon source is a Svelte component.
 */
export function isIconSourceComponent(source: IconSource): source is IconSourceComponent {
	return typeof source === 'function';
}

/**
 * Type guard to check if an icon source is a raw SVG string.
 */
export function isIconSourceSvg(source: IconSource): source is IconSourceSvg {
	return typeof source === 'string' && source.trimStart().startsWith('<svg');
}

/**
 * Type guard to check if an icon source is an Iconify icon.
 */
export function isIconSourceIconify(source: IconSource): source is IconSourceIconify {
	return typeof source === 'string' && source.trimStart().startsWith('iconify:');
}
