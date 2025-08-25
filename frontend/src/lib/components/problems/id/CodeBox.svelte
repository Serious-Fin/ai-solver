<script lang="ts">
	import Prism from 'prismjs'
	import 'prismjs/components/prism-go'

	let { code, buildError }: { code: string; buildError: string } = $props()

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
	<p class="build-error inter">{buildError}</p>
</article>

<style>
	.code {
		border: 1px solid rgba(0, 0, 0, 0.2);
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.15);
		font-size: 11pt;
		border-radius: 7px;
	}

	.build-error {
		color: red;
		white-space: pre-wrap;
		line-height: 1.4;
		font-size: 11pt;
	}

	@media (min-width: 768px) {
		code {
			font-size: 12pt;
		}
	}
</style>
