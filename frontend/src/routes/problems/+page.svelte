<script lang="ts">
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

<section>
	<ul>
		{#each data.problems as problem}
			<li>
				<a href="problems/{problem.id}">Title: {problem.title} Difficulty: {problem.difficulty}</a>
			</li>
		{/each}
	</ul>
</section>

<!-- each problem should have name, difficulty, completed or not -->

<style>
	section {
		background-color: #161c2e;
		width: 100vw;
		height: auto;
		min-height: 100vh;
		padding: 32px 0;
		box-sizing: border-box;
	}

	ul {
		list-style-type: none;
		margin: 0;
		padding: 0;
	}

	li {
		border: 1px solid green;
	}
</style>
