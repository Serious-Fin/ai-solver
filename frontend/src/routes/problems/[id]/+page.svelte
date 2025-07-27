<script lang="ts">
	import DescriptionBox from '$lib/components/problems/id/DescriptionBox.svelte';
	import CodeBox from '$lib/components/problems/id/CodeBox.svelte';
	import ChatBox from '$lib/components/problems/id/ChatBox.svelte';
	import TestBox from '$lib/components/problems/id/TestBox.svelte';

	import type { PageProps } from './$types';
	let { data }: PageProps = $props();

	let problemId: string = data.problem.id;
	let title: string = data.problem.title;
	let description: string = data.problem.description ?? '';
	let testCaseIds: number[] = data.problem.testCaseIds ?? [];
	let code: string = $state(data.problem.goPlaceholder ?? '');

	function updateCode(newCode: string) {
		code = newCode;
	}
</script>

<section>
	<article>
		<a class="back btn" href="/problems">
			<img class="img_icon" src="/back-arrow.svg" alt="return back arrow" />
		</a>
		<h1 class="inter-600">{title}</h1>
		<img class="completed_icon" src="/done-symbol.svg" alt="exercise already completed check" />
	</article>

	<DescriptionBox {description}></DescriptionBox>

	<CodeBox {code}></CodeBox>

	<ChatBox {code} {updateCode}></ChatBox>

	<TestBox {problemId} {testCaseIds} {code}></TestBox>
</section>

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
	section {
		background-color: #3b2645;
		width: 100vw;
		height: auto;
		min-height: 100vh;
		padding: 16px;
		box-sizing: border-box;
	}

	article {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	h1 {
		color: #ffffff;
		font-size: 18pt;
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

	.back {
		background-color: white;
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
