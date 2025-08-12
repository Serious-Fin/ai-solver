<script lang="ts">
	import type { User } from '$lib/api/users'
	import { handleError } from '$lib/helpers'

	let { user }: { user?: User } = $props()

	async function logout() {
		try {
			await fetch('/api/oauth/logout', { method: 'POST' })
			window.location.reload()
		} catch (err) {
			if (err instanceof Error) {
				handleError('Error logging out, try again', err)
				return
			}
		}
	}

	async function login() {
		window.location.href = '/login'
	}
</script>

{#if !user}
	<article class="inter">
		<button onclick={login} class="login">Log in</button>
	</article>
{:else}
	<article class="inter">
		<p class="logout_text">{user.email}</p>
		<button onclick={logout} class="logout"><img src="/logout.svg" alt="logout" /></button>
	</article>
{/if}

<style>
	article {
		width: fit-content;
		display: flex;
		align-items: center;
		gap: 20px;
	}

	p {
		color: white;
	}

	button {
		box-sizing: border-box;
	}

	.logout {
		width: fit-content;
		height: auto;
		padding: 4px;
	}

	.logout_text {
		font-size: 11pt;
	}

	.login {
		font-size: 12pt;
		padding: 4px 8px;
	}

	img {
		width: 30px;
		height: 30px;
	}
</style>
