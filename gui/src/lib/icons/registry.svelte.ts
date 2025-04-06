import { defaultIcons } from './defaults.ts';
import type { IconComponent, IconName, IconSet } from './types.ts';

export const registry = $state<Partial<IconSet>>({ ...defaultIcons });

/**
 * registerIconSet registers a new partial icon set, merging it with the existing registry.
 * Icons from the new set will override existing icons with the same name. This allows
 * plugins to replace default icons.
 *
 * @param iconSet A partial object mapping icon names to an `IconComponent`.
 */
export function registerIconSet(iconSet: Partial<IconSet>) {
	for (const [key, component] of Object.entries(iconSet)) {
		registry[key as IconName] = component;
	}
}

/**
 * getIconComponent returns the icon component for a given icon name.
 *
 * @param name The name of the icon to get the component for.
 * @returns The icon component for the given name.
 */
export function getIconComponent(name: IconName): IconComponent | undefined {
	return registry[name];
}
