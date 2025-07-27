<script lang="ts">
	import { validate } from '$lib/api/validate';
	import { TestStatusReporter } from '$lib/TestStatusReporter';
	import { handleError } from '$lib/helpers';

	let { problemId, testCaseIds, code }: { problemId: string; testCaseIds: number[]; code: string } =
		$props();

	let testStatusReporter = new TestStatusReporter(testCaseIds);
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
		<h2>Tests</h2>
	</header>
	{#each testStates as testState}
		<div>
			<p>Test {testState.id}</p>
			<p>Status: {testState.status}</p>
		</div>
	{/each}
	<footer>
		<button onclick={handleRunTests}> Run tests </button>
	</footer>
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
</style>
