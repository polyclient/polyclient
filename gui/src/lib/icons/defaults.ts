import { Home } from '@lucide/svelte';
import { type IconSet, iconsNames } from './types.ts';

const exampleRawSvg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle></svg>`;

export const defaultIcons: IconSet = {
	[iconsNames.railDatabase]: 'iconify:mdi:database',
	[iconsNames.railExplorer]: 'iconify:mdi:database',
	[iconsNames.railHistory]: 'iconify:mdi:database',
	[iconsNames.railPlugins]: 'iconify:mdi:database',
	[iconsNames.railAssistant]: 'iconify:mdi:database',
	[iconsNames.railSettings]: 'iconify:mdi:database',

	[iconsNames.dbSchema]: 'iconify:mdi:database',
	[iconsNames.dbTable]: 'iconify:mdi:database',
	[iconsNames.dbColumn]: 'iconify:mdi:database',
	[iconsNames.dbColumnPk]: 'iconify:mdi:database',
	[iconsNames.dbColumnFk]: 'iconify:mdi:database',
	[iconsNames.dbView]: 'iconify:mdi:database',
	[iconsNames.dbFunction]: 'iconify:mdi:database',
	[iconsNames.dbProcedure]: 'iconify:mdi:database',
	[iconsNames.dbTrigger]: 'iconify:mdi:database',
	[iconsNames.dbConnection]: 'iconify:mdi:database',
	[iconsNames.dbFolder]: 'iconify:mdi:database',

	[iconsNames.uiChevronRight]: 'iconify:mdi:database',
	[iconsNames.uiChevronDown]: 'iconify:mdi:database',
	[iconsNames.uiChevronLeft]: 'iconify:mdi:database',
	[iconsNames.uiChevronUp]: 'iconify:mdi:database',
	[iconsNames.uiArrowRight]: 'iconify:mdi:database',
	[iconsNames.uiArrowLeft]: 'iconify:mdi:database',
	[iconsNames.uiArrowUp]: 'iconify:mdi:database',
	[iconsNames.uiArrowDown]: 'iconify:mdi:database',
	[iconsNames.uiNext]: exampleRawSvg,
	[iconsNames.uiPrevious]: Home,
	[iconsNames.uiPlus]: 'iconify:mdi:database',
	[iconsNames.uiMinus]: 'iconify:mdi:database',
	[iconsNames.uiInfo]: 'iconify:mdi:database',
	[iconsNames.uiError]: 'iconify:mdi:database',
	[iconsNames.uiWarning]: 'iconify:mdi:database',
	[iconsNames.uiSuccess]: 'iconify:mdi:database',
};
