<script lang="ts">
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

<ul>
	{#each data.problems as problem}
		<li><a href="problems/{problem.id}">{problem.title}</a></li>
	{/each}
</ul>

<!-- each problem should have name, difficulty, completed or not -->
