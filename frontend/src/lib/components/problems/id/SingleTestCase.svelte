<script lang="ts">
	import type { SingleTestStatus } from '$lib/TestStatusReporter';
	import { TestStatus } from '$lib/TestStatusReporter';

	let { test }: { test: SingleTestStatus } = $props();
</script>

<div
	class="test-box {test.status === TestStatus.PASS
		? 'success'
		: test.status === TestStatus.FAIL
			? 'fail'
			: ''}"
>
	<header>
		<div class="title">
			<p class="test-box_header inter_700">Test {test.id}</p>
			{#if test.status === TestStatus.UNKNOWN}
				<img
					src="/test-not-run-icon.png"
					alt="question mark in circle"
					class="test-box_status-img"
				/>
			{:else if test.status === TestStatus.PASS}
				<img
					src="/test-success-icon.png"
					alt="green circle with white checkmark in the middle"
					class="test-box_status-img"
				/>
			{:else if test.status === TestStatus.FAIL}
				<img
					src="/test-fail-icon.png"
					alt="red circle with white X mark in the middle"
					class="test-box_status-img"
				/>
			{/if}
		</div>
		<button class="drop-down-arrow"><img src="/arrow_drop_down.png" alt="drop down arror" /></button
		>
	</header>
	{#if test.status === TestStatus.FAIL}
		<div class="body">
			<p>Got: {test.got}</p>
			<p>Want: {test.want}</p>
		</div>
	{/if}
</div>

<style>
	.test-box {
		width: 100%;
		background-color: #dcdcdc;
		border-radius: 5px;
		margin-bottom: 18px;
		box-sizing: border-box;
		padding: 14px 16px;
		box-shadow: 4px 4px 4px 0 rgba(0, 0, 0, 0.196);
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.title {
		display: flex;
		gap: 16px;
	}

	.success {
		background-color: #cdd7cb;
	}

	.fail {
		background-color: #d7cccb;
	}

	.test-box_header {
		color: #292929;
		font-size: 14pt;
	}

	.test-box_status-img {
		width: 20px;
		height: 20px;
	}

	.drop-down-arrow {
		border: none;
		background: none;
	}

	.inter_700 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
	}
</style>
