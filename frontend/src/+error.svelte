<script lang="ts">
	import { page } from '$app/state'
</script>

{#if !page.error}
	<h1>Unknown error occurred</h1>
{:else}
	<h1>{page.error.message}</h1>
{/if}
