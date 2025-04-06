<script lang="ts">
	import { getIconComponent } from './registry.svelte';
	import type { IconName } from './types.ts';

	type Props = {
		name: IconName;
		class?: string;
	};

	let { name, class: className }: Props = $props();

	const Component = getIconComponent(name);
</script>

{#if Component}
	<Component class={className} />
{/if}
