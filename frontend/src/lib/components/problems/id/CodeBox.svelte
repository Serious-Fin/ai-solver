<script lang="ts">
	import Prism from 'prismjs'
	import 'prismjs/components/prism-go'

	let { code }: { code: string } = $props()

	let codeElement: HTMLElement

	$effect(() => {
		codeElement.innerHTML = code
		Prism.highlightElement(codeElement)
	})
</script>

<article class="problem_article_box">
	<header class="problem_article_header">
		<h2 class="inter problem_article_header_text">Code</h2>
	</header>
	<pre class="code"><code class="language-go" bind:this={codeElement}>{code}</code></pre>
</article>

<style>
	.code {
		border: 1px solid rgba(0, 0, 0, 0.2);
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.15);
		font-size: 11pt;
		border-radius: 7px;
	}

	@media (min-width: 768px) {
		code {
			font-size: 12pt;
		}
	}
</style>
