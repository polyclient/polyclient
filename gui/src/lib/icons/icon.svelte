<script lang="ts">
	import Icon from '@iconify/svelte';
	import { getIconSource } from './registry.svelte';
	import { type IconName, isIconSourceComponent, isIconSourceIconify, isIconSourceSvg } from './types.ts';

	type Props = {
		name: IconName;
		class?: string;
	};

	let { name, class: className, ...rest }: Props = $props();

	const source = getIconSource(name);

	function extractIconify(icon: string) {
		return icon.replace('iconify:', '');
	}
</script>

{#if source}
	{#if isIconSourceComponent(source)}
		{@const Component = source}
		<Component
			class={className}
			{...rest}
		/>
	{/if}

	{#if isIconSourceSvg(source)}
		<span
			class={className}
			{...rest}
		>
			{@html source}
		</span>
	{/if}

	{#if isIconSourceIconify(source)}
		<Icon
			icon={extractIconify(source)}
			class={className}
			{...rest}
		/>
	{/if}
{/if}
