<script lang="ts">
	import { v4 as uuidv4 } from 'uuid';
	import { enhance } from '$app/forms';
	import type { SubmitFunction } from '@sveltejs/kit';
	import LoadingSpinner from '$lib/components/helpers/LoadingSpinner.svelte';
	import { handleError } from '$lib/helpers';

	let { code, updateCode }: { code: string; updateCode: (newCode: string) => void } = $props();

	let sessionId: string = uuidv4();
	let isLoading = $state(false);

	const handleQueryAgent: SubmitFunction = () => {
		isLoading = true;
		return async ({ update, result }) => {
			try {
				await update();
				if (result.type === 'success' && result.data?.response) {
					code = result.data.response;
					updateCode(code);
				} else if (result.type === 'failure') {
					throw Error(result.data?.message || 'Unknown server error occurred');
				} else {
					throw Error('Could not query agent');
				}
			} catch (err) {
				if (err instanceof Error) {
					handleError('Error sending message, try again later', err);
					return;
				}
			} finally {
				isLoading = false;
			}
		};
	};
</script>

<article class="problem_article_box">
	<header class="problem_article_header">
		<h2 class="inter problem_article_header_text">Chat box</h2>
	</header>
	<form method="POST" action="?/query" use:enhance={handleQueryAgent}>
		<textarea class="chat_box_input inter" name="input" id="input"></textarea>

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
	header {
		display: flex;
		justify-content: space-between;
	}

	.chat_box_input {
		width: 100%;
		max-width: 100%;
		min-height: 130px;
		font-size: 12pt;
		box-sizing: border-box;
		border-radius: 7px;
		padding: 12px 12px;
		font-weight: 400;
	}

	#agent {
		width: 110px;
		height: 40px;
		font-size: 12pt;
		border-radius: 7px;
		padding: 0 4px;
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

	@media (min-width: 768px) {
		#agent {
			font-size: 13pt;
		}

		.query_btn {
			width: 60px;
			height: 60px;
		}
	}
</style>
