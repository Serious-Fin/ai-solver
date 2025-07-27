<script lang="ts">
	import { v4 as uuidv4 } from 'uuid';
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';
	import LoadingSpinner from '$lib/components/helpers/LoadingSpinner.svelte';

	let { code, updateCode }: { code: string; updateCode: (newCode: string) => void } = $props();

	let sessionId: string = uuidv4();
	let isLoading = $state(false);

	const handleQueryAgent: SubmitFunction = () => {
		isLoading = true;
		return async ({ update, result }) => {
			try {
				await update();
				isLoading = false;
				if (result.type === 'success' && result.data?.response) {
					code = result.data.response;
					updateCode(code);
				} else {
					throw Error('could not query agent');
				}
			} catch (err) {
				// TODO: error table here
			}
		};
	};
</script>

<article>
	<header>
		<h2>Chat box</h2>
	</header>
	<form method="POST" action="?/query" use:enhance={handleQueryAgent}>
		<textarea class="chat_box_input inter-400" name="input" id="input"></textarea>

		<input type="hidden" name="code" value={code} />
		<input type="hidden" name="language" value="go" />
		<input type="hidden" name="sessionId" value={sessionId} />

		<footer class="block_footer send_box_footer">
			<select name="agent" id="agent">
				<option value="chatgpt">ChatGPT</option>
				<option value="gemini">Gemini</option>
			</select>

			<button class="query_btn" disabled={isLoading}>
				{#if isLoading}
					<LoadingSpinner />
				{:else}
					<img
						class="img_icon"
						src="/send-symbol.svg"
						alt="a paper plane icon symbolizing send action"
					/>
				{/if}
			</button>
		</footer>
	</form>
</article>

<style>
	article {
		background-color: #e9e9e9;
		border: 1px solid rgba(0, 0, 0, 0.8);
		border-radius: 10px;
		margin-top: 16px;
		padding: 8px 16px;
		box-sizing: border-box;
		max-width: 100%;
	}

	header {
		margin-bottom: 8px;
		display: flex;
		justify-content: space-between;
	}

	header h2 {
		color: rgba(0, 0, 0, 0.7);
		font-weight: 700;
		font-size: 14pt;

		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
	}

	.chat_box_input {
		width: 100%;
		max-width: 100%;
		min-height: 90px;
		font-size: 10pt;
	}

	.inter-400 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 400;
		font-style: normal;
	}

	.block_footer {
		margin-top: 8px;
		display: flex;
	}

	.send_box_footer {
		justify-content: space-between;
		align-items: start;
	}

	.query_btn {
		width: 50px;
		height: 50px;
		border-radius: 10px;
		border: 2px solid black;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: black;
	}

	.img_icon {
		width: 80%;
		height: 80%;
	}
</style>
