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

	function updateCode(newCode: string) {
		code = newCode;
	}
</script>

<section>
	<article>
		<a class="btn" href="/problems">
			<img class="img_icon" src="/arrow_back.svg" alt="return back arrow" />
		</a>
		<h1 class="inter-600">{title}</h1>
		<img class="completed_icon" src="/done-symbol.svg" alt="exercise already completed check" />
	</article>

	<DescriptionBox {description}></DescriptionBox>

	<CodeBox {code}></CodeBox>

	<ChatBox {code} {updateCode}></ChatBox>

	<TestBox {problemId} {testCases} {code}></TestBox>
</section>

<!-- TODO: modernize front-end design in figma for mobile -->
<!-- clicking on status message shows "not run yet" etc -->
<!-- change question mark with minus -->
<!-- use a different markdown parser because this one skips empty lines -->
<!-- TODO: implement front-end design in code for mobile -->
<!-- TODO: add front-end design for tablets in figma -->
<!-- TODO: implement front-end design for tablets in code -->
<!-- TODO: add front-end design for PCs in figma -->
<!-- TODO: implement front-end design for PCs in code -->

<style>
	section {
		background-color: #161c2e;
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

	.inter-600 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 600;
		font-style: normal;
	}

	.completed_icon {
		width: 13px;
		height: 13px;
		opacity: 80%;
	}

	@media (min-width: 760px) {
		section {
			background-color: green;
		}
	}
</style>
