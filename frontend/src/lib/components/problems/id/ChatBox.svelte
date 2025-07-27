<script lang="ts">
	import { v4 as uuidv4 } from 'uuid';
	import { enhance } from '$app/forms';

	let { code }: { code: string } = $props();

	let sessionId: string = uuidv4();

	const handleReturnedCode = () => {
		return async ({ result }) => {
			console.log(result);
			if (result.type === 'success') {
				code = result.data.response;
			}
		};
	};
</script>

<article>
	<header>
		<h2>Chat box</h2>
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

	.send {
		background-color: black;
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
</style>
