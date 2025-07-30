<script lang="ts">
	import Prism from 'prismjs';
	import 'prismjs/components/prism-go';

	let { code }: { code: string } = $props();

	let codeElement: HTMLElement;

	$effect(() => {
		codeElement.innerHTML = code;
		Prism.highlightElement(codeElement);
	});
</script>

<article>
	<header>
		<h2>Code</h2>
	</header>
	<pre class="code"><code class="language-go" bind:this={codeElement}>{code}</code></pre>
</article>

<style>
	article {
		background-color: #e9e9e9;
		margin-top: 32px;
		padding: 32px 16px;
		box-sizing: border-box;
		max-width: 100%;
	}

	header {
		margin-bottom: 24px;
	}

	header h2 {
		color: rgba(0, 0, 0, 0.7);
		font-weight: 700;
		font-size: 20pt;

		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
	}

	.code {
		border: 1px solid rgba(0, 0, 0, 0.2);
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.15);
		font-size: 11pt;
		border-radius: 7px;
	}
</style>
