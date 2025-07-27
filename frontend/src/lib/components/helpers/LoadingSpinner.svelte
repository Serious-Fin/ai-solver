<div class="spinner"></div>

<style>
	.spinner {
		border: 4px solid rgba(255, 255, 255, 0.3);
		border-top: 4px solid white;
		border-radius: 50%;
		width: 20px;
		height: 20px;
		animation: rotate 1s linear infinite;
	}

	@keyframes rotate {
		0% {
			transform: rotate(0deg);
		}

		50% {
			transform: rotate(180deg);
		}

		100% {
			transform: rotate(360deg);
		}
	}
</style>
