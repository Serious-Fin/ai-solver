<script lang="ts">
	import { page } from '$app/state'
	import type { User } from '$lib/api/users'
	import { getUsernameFromEmail, handleError } from '$lib/helpers'

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
		window.location.href = `/login?redirectTo=${encodeURIComponent(page.url.pathname)}`
	}
</script>

{#if !user}
	<article class="inter">
		<button onclick={login} class="login">Log in</button>
	</article>
{:else}
	<article class="inter">
		<p class="part_email">{getUsernameFromEmail(user.email)}</p>
		<p class="full_email">{user.email}</p>
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

	.login {
		font-size: 12pt;
		padding: 4px 8px;
		width: 64px;
	}

	img {
		width: 30px;
		height: 30px;
	}

	.part_email {
		display: none;
		font-size: 11pt;
	}

	.full_email {
		display: none;
		font-size: 11pt;
	}

	@media (min-width: 768px) {
		.part_email {
			display: block;
		}

		.full_email {
			display: none;
		}
	}

	@media (min-width: 1024px) {
		.part_email {
			display: none;
		}

		.full_email {
			display: block;
		}
	}
</style>
