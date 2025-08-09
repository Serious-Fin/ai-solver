<script lang="ts">
	import DescriptionBox from '$lib/components/problems/id/DescriptionBox.svelte';
	import CodeBox from '$lib/components/problems/id/CodeBox.svelte';
	import ChatBox from '$lib/components/problems/id/ChatBox.svelte';
	import TestBox from '$lib/components/problems/id/TestBox.svelte';

	import type { PageProps } from './$types';
	import type { TestCase } from '$lib/api/problems';
	let { data }: PageProps = $props();

	let problemId: string = data.problem.id;
	let title: string = data.problem.title;
	let description: string = data.problem.description ?? '';
	let testCases: TestCase[] = data.problem.testCases ?? [];
	let code: string = $state(data.problem.goPlaceholder ?? '');
	let isCompleted: boolean = data.problem.isCompleted;

	function updateCode(newCode: string) {
		code = newCode;
	}
</script>

<section>
	<article>
		<a class="btn" href="/problems">
			<img class="img_icon" src="/arrow_back.svg" alt="return back arrow" />
		</a>
		<h1 class="inter">{title}</h1>
		{#if isCompleted}
			<img class="completed_icon" src="/done-symbol.svg" alt="exercise already completed check" />
		{/if}
	</article>

	<DescriptionBox {description}></DescriptionBox>

	<CodeBox {code}></CodeBox>

	<ChatBox {code} {updateCode}></ChatBox>

	<TestBox {problemId} {testCases} {code}></TestBox>
</section>

<style>
	section {
		background-color: var(--background);
		width: 100vw;
		height: auto;
		min-height: 100vh;
		padding: 32px 0;
		box-sizing: border-box;
	}

	article {
		display: grid;
		grid-template-columns: 50px auto 50px;
		align-items: center;
		gap: 10px;
		padding-left: 10px;
	}

	h1 {
		color: #ffffff;
		font-size: 18pt;
		font-weight: 600;
	}

	.btn {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.img_icon {
		width: 80%;
		height: 80%;
	}

	.completed_icon {
		width: 20px;
		height: 20px;
	}

	@media (min-width: 768px) {
		h1 {
			font-size: 20pt;
		}
	}

	@media (min-width: 1024px) {
		h1 {
			font-size: 22pt;
		}
	}
</style>
