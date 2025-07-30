<script lang="ts">
	import type { SingleTestStatus } from '$lib/TestStatusReporter';
	import { TestStatus } from '$lib/TestStatusReporter';
	import { printHumanReadable } from '$lib/helpers';

	let { test }: { test: SingleTestStatus } = $props();
	let isExpanded = $state(false);

	function toggleDropdown() {
		isExpanded = !isExpanded;
	}
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
		<button class="drop-down-arrow" onclick={toggleDropdown}>
			{#if isExpanded}
				<img src="/arrow_drop_up.png" alt="up arrow" />
			{:else}
				<img src="/arrow_drop_down.png" alt="drop down arrow" />
			{/if}
		</button>
	</header>
	<div class="body {isExpanded ? '' : 'invisible'}">
		{#each test.inputs as input, index}
			<p><b>Input {index + 1}:</b> {printHumanReadable(input)}</p>
		{/each}
		<p><b>Expect:</b> {printHumanReadable(test.output)}</p>
		{#if test.status === TestStatus.FAIL}
			<p><b>Got:</b> {printHumanReadable(test.got)}</p>
		{/if}
	</div>
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

	.body p {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 400;
		font-style: normal;
		line-height: 1.6;
		font-size: 11pt;
	}

	.invisible {
		display: none;
	}

	.test-box_header {
		color: #292929;
		font-size: 12pt;
	}

	.test-box_status-img {
		width: 16px;
		height: 16px;
	}

	.drop-down-arrow {
		border: none;
		background: none;
		width: 70px;
	}

	.inter_700 {
		font-family: 'Inter', sans-serif;
		font-optical-sizing: auto;
		font-weight: 700;
		font-style: normal;
	}
</style>
