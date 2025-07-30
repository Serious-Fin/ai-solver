<script lang="ts">
	import { validate } from '$lib/api/validate';
	import { TestStatusReporter } from '$lib/TestStatusReporter';
	import { handleError } from '$lib/helpers';
	import SingleTestCase from './SingleTestCase.svelte';
	import type { TestCase } from '$lib/api/problems';
	import LoadingSpinner from '$lib/components/helpers/LoadingSpinner.svelte';

	let { problemId, testCases, code }: { problemId: string; testCases: TestCase[]; code: string } =
		$props();

	let testStatusReporter = new TestStatusReporter(testCases);
	let testStates = $state(testStatusReporter.GetTestStatuses());
	let isLoading = $state(false);

	const handleRunTests = async () => {
		isLoading = true;
		try {
			const testRunOutput = await validate({
				problemId,
				code,
				language: 'go'
			});
			console.log(testRunOutput);
			testStatusReporter.UpdateTestStatuses(testRunOutput);
			testStates = testStatusReporter.GetTestStatuses();
		} catch (err) {
			if (err instanceof Error) {
				handleError('Error running tests, try again later', err);
				return;
			}
		} finally {
			isLoading = false;
		}
	};
</script>

<article>
	<header>
		<h2 class="inter_700">Tests</h2>
	</header>
	{#each testStates as test}
		<SingleTestCase {test}></SingleTestCase>
	{/each}
	<footer>
		<button class="inter_700" onclick={handleRunTests} disabled={isLoading}>
			{#if isLoading}
				<LoadingSpinner></LoadingSpinner>
			{:else}
				Run tests
			{/if}
		</button>
	</footer>
</article>

<style>
	article {
		background-color: #e9e9e9;
		margin-top: 32px;
		padding: 32px 24px;
		box-sizing: border-box;
		max-width: 100%;
	}

	header {
		margin-bottom: 24px;
	}

	header h2 {
		color: rgba(0, 0, 0, 0.7);
		font-size: 18pt;
	}

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
	}

	footer {
		display: flex;
		justify-content: end;
	}

	.inter_700 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
	}
</style>
