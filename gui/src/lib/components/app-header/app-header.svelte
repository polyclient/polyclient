<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import Icon from '@iconify/svelte';
</script>

<header class="h-[50px] flex justify-between border-b px-2">
	<div class="flex items-center gap-2">
		<Button
			variant="outline"
			size="icon"
			class="h-8"
		>
			<Icon icon="solar:sidebar-minimalistic-outline" />
		</Button>
		<Separator orientation="vertical" />
		<nav>
			<a href="/">Namespace</a> /
			<a href="/about">Database</a>
		</nav>
	</div>
</header>
