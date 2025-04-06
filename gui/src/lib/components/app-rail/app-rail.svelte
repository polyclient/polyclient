<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Separator from '$lib/components/ui/separator';
	import { Icon } from '$lib/icons';
</script>

<aside class="flex items-center justify-between border">
	<div class="flex items-center">
		<Button
			variant="ghost"
			size="sm"
			class="rounded-none"
		>
			<Icon name="rail:history" />
			Query History
		</Button>
		<Button
			variant="ghost"
			size="sm"
			class="rounded-none"
		>
			<Icon name="rail:assistant" />
			AI Assistant
		</Button>
		<Button
			variant="ghost"
			size="sm"
			class="rounded-none"
		>
			<Icon name="rail:plugins" />
			Plugins
		</Button>
	</div>

	<div class="flex items-center">
		<Button
			variant="ghost"
			size="sm"
			class="rounded-none"
		>
			<Icon name="rail:notifications" />
			Notifications
		</Button>
		<Button
			variant="ghost"
			size="sm"
			class="rounded-none"
		>
			<Icon name="rail:settings" />
			Settings
		</Button>
	</div>
</aside>
