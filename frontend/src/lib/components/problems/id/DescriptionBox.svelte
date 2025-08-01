<script lang="ts">
	import { marked } from 'marked';
	let { description }: { description: string } = $props();
</script>

<article>
	<header>
		<h2>Description</h2>
	</header>
	<div>
		{@html marked.parse(description, { breaks: true })}
	</div>
</article>

<style>
	article {
		background-color: #e9e9e9;
		margin-top: 32px;
		padding: 32px 24px;
		box-sizing: border-box;
		max-width: 100%;
	}

	header {
		margin-bottom: 24px;
		display: flex;
		justify-content: space-between;
	}

	header h2 {
		color: rgba(0, 0, 0, 0.7);

		font-size: 20pt;
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
	}

	div {
		font-family:
			-apple-system,
			BlinkMacSystemFont,
			Segoe UI,
			Helvetica,
			Arial,
			sans-serif,
			Apple Color Emoji,
			Segoe UI Emoji;

		line-height: 1.5;
		font-size: 12pt;
	}

	@media (min-width: 768px) {
		article {
			padding: 36px 32px;
		}
	}
</style>
