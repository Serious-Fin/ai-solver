<script lang="ts">
	import { onMount } from 'svelte';
	import { enhance } from '$app/forms';
	import Prism from 'prismjs';
	import 'prismjs/components/prism-go';
	import 'prismjs/components/prism-c.js';
	import { marked } from 'marked';
	import { v4 as uuidv4 } from 'uuid';
	import { TestStatusReporter, type TestRunOutput } from '$lib/TestStatusReporter';

	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
	let sessionId: string = $state(uuidv4());

	onMount(async () => {
		// @ts-ignore
		await import('prismjs/components/prism-cpp.js');
		Prism.highlightAll();
	});

	let code = $state(data.problem.goPlaceholder ?? '');
	let isLoading = $state(false);
	let codeElement: HTMLElement;
	let testStatusReporter = new TestStatusReporter(data.problem.testCaseIds ?? []);
	let testStates = $state(testStatusReporter.GetTestStatuses());

	$effect(() => {
		codeElement.innerHTML = code;
		Prism.highlightElement(codeElement);
	});

	const handleReturnedCode = () => {
		return async ({ result }) => {
			console.log(result);
			if (result.type === 'success') {
				code = result.data.response;
			}
		};
	};

	const runTests = async () => {
		isLoading = true;
		try {
			const testRunOutput = await fetch('http://localhost:8080/validate', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					problemId: data.problem.id,
					code,
					language: 'go'
				})
			});
			if (!testRunOutput.ok) {
				throw new Error(`HTTP error with status ${testRunOutput.status}`);
			}
			const response: TestRunOutput = await testRunOutput.json();
			testStatusReporter.UpdateTestStatuses(response);
			testStates = testStatusReporter.GetTestStatuses();
		} catch (err) {
			console.log(err);
		} finally {
			isLoading = false;
		}
	};
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
		<div class="classy-font">
			{@html marked.parse(data.problem.description ?? '')}
		</div>
	</article>

	<article class="block">
		<header class="block_header">
			<h2 class="block_header_text inter-700">Code</h2>
		</header>
		<pre class="code"><code class="language-go" bind:this={codeElement}>{code}</code></pre>
	</article>

	<article class="block">
		<header class="block_header">
			<h2 class="block_header_text inter-700">Chat box</h2>
		</header>
		<form method="POST" action="?/query" use:enhance={handleReturnedCode}>
			<textarea class="chat_box_input inter-400" name="input" id="input"></textarea>

			<input type="hidden" name="code" value={code} />
			<input type="hidden" name="language" value="go" />
			<input type="hidden" name="sessionId" value={sessionId} />

			<footer class="block_footer send_box_footer">
				<select name="agent" id="agent">
					<option value="chatgpt">ChatGPT</option>
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

	<article class="block">
		<header class="block_header">
			<h2 class="block_header_text inter-700">Tests</h2>
		</header>
		{#each testStates as testState}
			<div>
				<p>Test {testState.id}</p>
				<p>Status: {testState.status}</p>
			</div>
		{/each}
		<footer>
			<button onclick={runTests}> Run tests </button>
		</footer>
	</article>
</section>

<!-- TODO: extract code for different components to different files -->
<!-- TODO: make model dropdown actually work -->
<!-- TODO: add loading for "send" button while waiting for response -->
<!-- TODO: add error pop-out for when something fails -->
<!-- TODO: make test design in figma -->
<!-- TODO: add images instead of status numbers for tests -->
<!-- TODO: add fail reason explanation -->
<!-- TODO: modernize front-end design in figma for mobile -->
<!-- TODO: implement front-end design in code for mobile -->
<!-- TODO: add front-end design for tablets in figma -->
<!-- TODO: implement front-end design for tablets in code -->
<!-- TODO: add front-end design for PCs in figma -->
<!-- TODO: implement front-end design for PCs in code -->

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

	.classy-font {
		font-family:
			-apple-system,
			BlinkMacSystemFont,
			Segoe UI,
			Helvetica,
			Arial,
			sans-serif,
			Apple Color Emoji,
			Segoe UI Emoji;
	}

	@media (min-width: 760px) {
		#main-area {
			background-color: green;
		}
	}
</style>
