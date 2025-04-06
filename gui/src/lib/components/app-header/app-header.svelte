<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import { Icon } from '$lib/icons';
</script>

<header class="h-full flex justify-between border-b px-1">
	<div class="flex items-center gap-1">
		<Button
			variant="ghost"
			size="icon"
		>
			<Icon name="ui:previous" />
		</Button>
		<Button
			variant="ghost"
			size="icon"
		>
			<Icon name="ui:next" />
		</Button>
	</div>
</header>
