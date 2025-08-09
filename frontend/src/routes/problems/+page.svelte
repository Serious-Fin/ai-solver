<script lang="ts">
	import { getDifficultyName } from '$lib/helpers';
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

<section>
	<article>
		<h1 class="inter-600">All problems</h1>
	</article>
	<article>
		<ul>
			{#each data.problems as problem}
				<a href="problems/{problem.id}">
					<li>
						<p class="inter-400">{problem.title}</p>
						<p class="inter-600 difficulty {getDifficultyName(problem.difficulty)}">
							{getDifficultyName(problem.difficulty)}
						</p>
						{#if problem.isCompleted}
							<img
								class="completed_icon"
								src="/done-symbol.svg"
								alt="exercise already completed check"
							/>
						{/if}
					</li>
				</a>
			{/each}
		</ul>
	</article>
</section>

<style>
	section {
		background-color: #161c2e;
		width: 100vw;
		height: auto;
		min-height: 100vh;
		padding: 32px 0;
		box-sizing: border-box;
	}

	h1 {
		color: #ffffff;
		font-size: 18pt;
		margin-bottom: 20px;
		margin-left: 20px;
	}

	.inter-600 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 600;
		font-style: normal;
	}

	.inter-400 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 400;
		font-style: normal;
	}

	.difficulty {
		font-size: 10pt;
	}

	.easy {
		color: green;
	}

	.medium {
		color: orange;
	}

	.hard {
		color: red;
	}

	ul {
		list-style-type: none;
		margin: 0;
		padding: 0;
	}

	li {
		background-color: #e9e9e9;
		padding: 20px 20px;
		margin-bottom: 30px;
		box-sizing: border-box;

		font-size: 12pt;
		color: black;

		display: grid;
		grid-template-columns: auto 60px 30px;
		align-items: center;
		gap: 15px;
	}

	a {
		text-decoration: none;
	}
</style>
