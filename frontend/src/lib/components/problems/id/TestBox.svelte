<script lang="ts">
	import { TestStatusReporter, type TestRunOutput } from '$lib/TestStatusReporter'
	import { handleFrontendError } from '$lib/helpers'
	import SingleTestCase from './SingleTestCase.svelte'
	import { type TestCase } from '$lib/api/problems'
	import LoadingSpinner from '$lib/components/helpers/LoadingSpinner.svelte'
	import { enhance } from '$app/forms'
	import type { SubmitFunction } from '@sveltejs/kit'

	let {
		problemId,
		testCases,
		code,
		markProblemCompletedFunc
	}: {
		problemId: number
		testCases: TestCase[]
		code: string
		markProblemCompletedFunc: () => Promise<void>
	} = $props()

	let testStatusReporter = new TestStatusReporter(testCases)
	let testStates = $state(testStatusReporter.GetTestStatuses())
	let isLoading = $state(false)

	async function updateTests(testRunOutput: TestRunOutput) {
		testStatusReporter.UpdateTestStatuses(testRunOutput)
		testStates = testStatusReporter.GetTestStatuses()
		if (testStatusReporter.AllTestsSuccessful()) {
			await markProblemCompletedFunc()
		}
	}

	const handleRunTests: SubmitFunction = () => {
		isLoading = true
		return async ({ update, result }) => {
			try {
				await update()
				if (result.type === 'success' && result.data?.response) {
					const testRunOutput = result.data.response
					updateTests(testRunOutput)
				} else if (result.type === 'failure') {
					throw Error(result.data?.message || 'Unknown server error occurred')
				} else {
					throw Error('Could not run tests')
				}
			} catch (err) {
				if (err instanceof Error) {
					handleFrontendError('Error running tests, try again later', err)
					return
				}
			} finally {
				isLoading = false
			}
		}
	}
</script>

<article class="problem_article_box">
	<header class="problem_article_header">
		<h2 class="inter problem_article_header_text">Tests</h2>
	</header>
	{#each testStates as test}
		<SingleTestCase {test}></SingleTestCase>
	{/each}
	<footer>
		<form method="POST" action="?/runTests" use:enhance={handleRunTests}>
			<input type="hidden" name="code" value={code} />
			<input type="hidden" name="language" value="go" />
			<input type="hidden" name="problemId" value={problemId} />

			<button type="submit" class="inter" disabled={isLoading}>
				{#if isLoading}
					<LoadingSpinner></LoadingSpinner>
				{:else}
					Run tests
				{/if}
			</button>
		</form>
	</footer>
</article>

<style>
	button {
		background-color: black;
		border: none;
		width: 100px;
		height: 40px;
		color: white;
		font-size: 12pt;
		border-radius: 5px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 700;
	}

	footer {
		display: flex;
		justify-content: end;
	}
</style>
