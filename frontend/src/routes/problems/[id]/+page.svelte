<script lang="ts">
	import DescriptionBox from '$lib/components/problems/id/DescriptionBox.svelte'
	import CodeBox from '$lib/components/problems/id/CodeBox.svelte'
	import ChatBox from '$lib/components/problems/id/ChatBox.svelte'
	import TestBox from '$lib/components/problems/id/TestBox.svelte'

	import type { PageProps } from './$types'
	import { markProblemCompleted, type TestCase } from '$lib/api/problems'
	import UserBox from '$lib/components/UserBox.svelte'
	import { handleFrontendError, showSuccess, showWarning } from '$lib/helpers'
	let { data }: PageProps = $props()

	let problemId: string = data.problem.id
	let title: string = data.problem.title
	let description: string = data.problem.description ?? ''
	let testCases: TestCase[] = data.problem.testCases ?? []
	let code: string = $state(data.problem.goPlaceholder ?? '')
	let isCompleted: boolean = data.problem.isCompleted
	const user = data.user

	function updateCode(newCode: string) {
		code = newCode
	}

	const markProblemCompletedFunc = async () => {
		if (user) {
			try {
				await markProblemCompleted(problemId, user.id)
				showSuccess('Problem completed!')
			} catch (err) {
				if (err instanceof Error) {
					handleFrontendError('Error saving progress, try again', err)
				}
			}
		} else {
			showWarning('Problem completed but log in to save progress')
		}
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
		{:else}
			<p></p>
		{/if}
		<UserBox {user}></UserBox>
	</article>

	<DescriptionBox {description}></DescriptionBox>

	<CodeBox {code}></CodeBox>

	<ChatBox {code} {updateCode}></ChatBox>

	<TestBox {problemId} {testCases} {code} {markProblemCompletedFunc}></TestBox>
</section>

<!-- TODO: after completing the problem, display completed icon immediately-->

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
		grid-template-columns: min-content auto 32px min-content;
		align-items: center;
		gap: 10px;
		padding: 0 5px;
	}

	h1 {
		color: #ffffff;
		font-size: 18pt;
		font-weight: 600;
	}

	.btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.img_icon {
		width: 100%;
		height: 100%;
	}

	.completed_icon {
		width: 20px;
		height: 20px;
	}

	@media (min-width: 768px) {
		h1 {
			font-size: 20pt;
		}

		article {
			padding: 0 20px;
		}
	}

	@media (min-width: 1024px) {
		h1 {
			font-size: 22pt;
		}
	}
</style>
