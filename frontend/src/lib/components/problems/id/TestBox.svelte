<script lang="ts">
	import { TestStatusReporter, type TestRunOutput } from '$lib/TestStatusReporter';

	let { problemId, testCaseIds, code }: { problemId: string; testCaseIds: number[]; code: string } =
		$props();

	let testStatusReporter = new TestStatusReporter(testCaseIds);
	let testStates = $state(testStatusReporter.GetTestStatuses());
	let isLoading = $state(false);

	const runTests = async () => {
		isLoading = true;
		try {
			const testRunOutput = await fetch('http://localhost:8080/validate', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					problemId,
					code,
					language: 'go'
				})
			});
			if (!testRunOutput.ok) {
				throw new Error(`HTTP error with status ${testRunOutput.status}`);
			}
			const response: TestRunOutput = await testRunOutput.json();
			testStatusReporter.UpdateTestStatuses(response);
			testStates = testStatusReporter.GetTestStatuses();
		} catch (err) {
			console.log(err);
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
		<button onclick={runTests}> Run tests </button>
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
