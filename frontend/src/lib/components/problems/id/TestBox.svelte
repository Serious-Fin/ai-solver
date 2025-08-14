<script lang="ts">
	import { validate } from '$lib/api/validate'
	import { TestStatusReporter } from '$lib/TestStatusReporter'
	import { handleFrontendError } from '$lib/helpers'
	import SingleTestCase from './SingleTestCase.svelte'
	import { type TestCase } from '$lib/api/problems'
	import LoadingSpinner from '$lib/components/helpers/LoadingSpinner.svelte'

	let {
		problemId,
		testCases,
		code,
		markProblemCompletedFunc
	}: {
		problemId: string
		testCases: TestCase[]
		code: string
		markProblemCompletedFunc: () => Promise<void>
	} = $props()

	let testStatusReporter = new TestStatusReporter(testCases)
	let testStates = $state(testStatusReporter.GetTestStatuses())
	let isLoading = $state(false)

	const handleRunTests = async () => {
		isLoading = true
		try {
			const testRunOutput = await validate({
				problemId,
				code,
				language: 'go'
			})
			testStatusReporter.UpdateTestStatuses(testRunOutput)
			testStates = testStatusReporter.GetTestStatuses()
			if (testStatusReporter.AllTestsSuccessful()) {
				await markProblemCompletedFunc()
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
</script>

<article class="problem_article_box">
	<header class="problem_article_header">
		<h2 class="inter problem_article_header_text">Tests</h2>
	</header>
	{#each testStates as test}
		<SingleTestCase {test}></SingleTestCase>
	{/each}
	<footer>
		<button class="inter" onclick={handleRunTests} disabled={isLoading}>
			{#if isLoading}
				<LoadingSpinner></LoadingSpinner>
			{:else}
				Run tests
			{/if}
		</button>
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
