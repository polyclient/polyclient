import { polyclientDefault } from './icon-sets/polyclient-default.ts';
import type { IconName, IconSet, IconSource } from './types.ts';

export const registry = $state<Partial<IconSet>>({ ...polyclientDefault });

/**
 * registerIconSet registers a new partial icon set, merging it with the existing registry.
 * Icons from the new set will override existing icons with the same name. This allows
 * plugins to replace default icons.
 *
 * @param iconSet A partial object mapping icon names to an `IconSource`.
 */
export function registerIconSet(iconSet: Partial<IconSet>) {
	for (const [key, source] of Object.entries(iconSet)) {
		registry[key as IconName] = source;
	}
}

/**
 * getIconSource returns the icon source for a given icon name.
 *
 * @param name The name of the icon to get the source for.
 * @returns The icon source for the given name.
 */
export function getIconSource(name: IconName): IconSource | undefined {
	return registry[name];
}
