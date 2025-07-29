<script lang="ts">
	import { validate } from '$lib/api/validate';
	import { TestStatusReporter } from '$lib/TestStatusReporter';
	import { handleError } from '$lib/helpers';
	import SingleTestCase from './SingleTestCase.svelte';
	import type { TestCase } from '$lib/api/problems';

	let { problemId, testCases, code }: { problemId: string; testCases: TestCase[]; code: string } =
		$props();

	const testCaseIds = testCases.map((testCase) => testCase.id);
	let testStatusReporter = new TestStatusReporter(testCaseIds);
	let testStates = $state(testStatusReporter.GetTestStatuses());
	testStatusReporter.UpdateTestStatuses({
		succeededTests: [0, 1, 2],
		failedTests: [
			{ id: 3, want: '[1, 2, 3]', got: '[3, 2, 1]', message: 'wrong output' },
			{ id: 4, want: '[7, 8]', got: '[]', message: 'wrong output' },
			{ id: 5, want: '"foo"', got: '"bar"', message: 'wrong output' }
		]
	});
	let isLoading = $state(false);

	const handleRunTests = async () => {
		isLoading = true;
		try {
			const testRunOutput = await validate({
				problemId,
				code,
				language: 'go'
			});
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
		<button onclick={handleRunTests}> Run tests </button>
	</footer>
</article>

<style>
	article {
		background-color: #e9e9e9;
		border: 1px solid rgba(0, 0, 0, 0.8);
		border-radius: 5px;
		margin-top: 16px;
		padding: 24px 32px;
		box-sizing: border-box;
		max-width: 100%;
	}

	header {
		margin-bottom: 16px;
	}

	header h2 {
		color: rgba(0, 0, 0, 0.7);
		font-size: 18pt;
	}

	.inter_700 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
	}
</style>
