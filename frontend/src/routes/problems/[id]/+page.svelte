<script lang="ts">
	import { onMount } from 'svelte';
	import Prism from 'prismjs';
	import 'prismjs/components/prism-go';
	import { marked } from 'marked';

	import type { PageProps } from './$types';
	let { data }: PageProps = $props();

	onMount(() => {
		Prism.highlightAll();
	});

	let code = $state(data.problem.goCode);
</script>

<section id="main-area">
	<article id="header">
		<a class="back btn" href="/problems">
			<img class="img_icon" src="/back-arrow.svg" alt="return back arrow" />
		</a>
		<h1 id="main_header" class="inter-600">{data.problem.title}</h1>
		<img class="completed_icon" src="/done-symbol.svg" alt="exercise already completed check" />
	</article>

	<article class="block">
		<header class="block_header">
			<h2 class="block_header_text inter-700">Description</h2>
		</header>
		{@html marked.parse(data.problem.description ?? '')}
	</article>

	<article class="block">
		<header class="block_header">
			<h2 class="block_header_text inter-700">Code</h2>

			<select name="programming-language" id="programming-language">
				<option value="go">Go</option>
				<option value="cpp">C++</option>
				<option value="csharp">C#</option>
			</select>
		</header>
		<pre class="code"><code class="language-go">{code}</code></pre>
	</article>

	<article class="block">
		<header class="block_header">
			<h2 class="block_header_text inter-700">Chat box</h2>
		</header>
		<form action="#">
			<textarea class="chat_box_input inter-400"></textarea>

			<footer class="block_footer send_box_footer">
				<select name="assistant" id="assistant">
					<option value="gpt-4_1">GPT-4.1</option>
					<option value="gpt-4_1_mini">GPT-4.1 Mini</option>
					<option value="gemini_2_5_flash">Gemini 2.5 Flash</option>
				</select>

				<button class="send btn">
					<img
						class="img_icon"
						src="/send-symbol.svg"
						alt="a paper plane icon symbolizing send action"
					/>
				</button>
			</footer>
		</form>
	</article>
</section>

<style>
	#main-area {
		background-color: #3b2645;
		width: 100vw;
		height: auto;
		min-height: 100vh;
		padding: 16px;
		box-sizing: border-box;
	}

	#header {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	#main_header {
		color: #ffffff;
		font-size: 18pt;
	}

	.block {
		background-color: #e9e9e9;
		border: 1px solid rgba(0, 0, 0, 0.8);
		border-radius: 10px;
		margin-top: 16px;
		padding: 8px 16px;
		box-sizing: border-box;
		max-width: 100%;
	}

	.chat_box_input {
		width: 100%;
		max-width: 100%;
		min-height: 90px;
		font-size: 10pt;
	}

	.code {
		border: 1px solid rgba(0, 0, 0, 0.2);
		box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.15);
		font-size: 10pt;
	}

	.btn {
		width: 50px;
		height: 50px;
		border-radius: 10px;
		border: 2px solid black;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.img_icon {
		width: 80%;
		height: 80%;
	}

	.send {
		background-color: black;
	}

	.back {
		background-color: white;
	}

	.block_header {
		margin-bottom: 8px;
		display: flex;
		justify-content: space-between;
	}

	.block_footer {
		margin-top: 8px;
		display: flex;
	}

	.send_box_footer {
		justify-content: space-between;
		align-items: start;
	}

	.block_description {
		max-width: 30ch;
		font-size: 11pt;
		line-height: 1.6;
	}

	.block_header_text {
		color: rgba(0, 0, 0, 0.7);
		font-weight: 700;
		font-size: 14pt;
	}

	.inter-700 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
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

	.completed_icon {
		width: 13px;
		height: 13px;
		opacity: 80%;
	}

	@media (min-width: 760px) {
		#main-area {
			background-color: green;
		}
	}
</style>
