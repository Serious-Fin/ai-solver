<script lang="ts">
	import { marked } from 'marked';
	let { description }: { description: string } = $props();
</script>

<article class="problem_article_box">
	<header class="problem_article_header">
		<h2 class="inter problem_article_header_text">Description</h2>
	</header>
	<div>
		{@html marked.parse(description, { breaks: true })}
	</div>
</article>

<style>
	header {
		display: flex;
		justify-content: space-between;
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
</style>
