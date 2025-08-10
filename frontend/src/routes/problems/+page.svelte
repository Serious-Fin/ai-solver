<script lang="ts">
	import { getDifficultyName } from '$lib/helpers'
	import type { PageProps } from './$types'
	let { data }: PageProps = $props()
	console.log(data.sessionId)
</script>

<section>
	<article>
		<h1 class="inter">All problems</h1>
		<h1>{data.sessionId}</h1>
	</article>
	<article>
		<ul>
			{#each data.problems as problem}
				<a href="problems/{problem.id}">
					<li>
						<p class="inter title">{problem.title}</p>
						<p class="inter difficulty {getDifficultyName(problem.difficulty)}">
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
		background-color: var(--background);
		width: 100vw;
		height: auto;
		min-height: 100vh;
		padding: 32px 0;
		box-sizing: border-box;
	}

	h1 {
		color: #ffffff;
		font-size: 18pt;
		margin-bottom: 32px;
		margin-left: 20px;
		font-weight: 600;
	}

	.title {
		font-weight: 400;
	}

	.difficulty {
		font-size: 10pt;
		font-weight: 600;
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
		padding: 0;
		margin: 0;
	}

	li {
		background-color: var(--foreground);
		padding: 20px 20px;
		margin-bottom: 30px;
		box-sizing: border-box;
		max-width: 768px;

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

	@media (min-width: 768px) {
		li {
			margin-left: 30px;
		}
	}
</style>
