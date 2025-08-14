<script lang="ts">
	import { page } from '$app/state'
	import type { User } from '$lib/api/users'
	import { handleFrontendError } from '$lib/helpers'
	import { onMount, onDestroy } from 'svelte'

	let { user }: { user?: User } = $props()
	let visibleUserInfo: boolean = $state(false)
	let wrapper: HTMLElement | null = $state(null)

	function toggleUserInfo() {
		visibleUserInfo = !visibleUserInfo
	}

	function handleClickOutside(event: MouseEvent) {
		if (wrapper && event.target instanceof Node && !wrapper.contains(event.target)) {
			visibleUserInfo = false
		}
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside)
	})

	async function logout() {
		try {
			await fetch('/api/oauth/logout', { method: 'POST' })
			window.location.reload()
		} catch (err) {
			if (err instanceof Error) {
				handleFrontendError('Error logging out, try again', err)
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
		<button class="btn_wrapper" onclick={toggleUserInfo} bind:this={wrapper}>
			{#if user.profilePic}
				<img src={user.profilePic} alt="user avatar" class="profile-pic" />
			{:else}
				<img
					src="https://digitalhealthskills.com/wp-content/uploads/2022/11/3da39-no-user-image-icon-27.png"
					alt="user avatar"
					class="profile-pic"
				/>
			{/if}
			{#if visibleUserInfo}
				<article id="logged_in_popup">
					<p>Logged in as:</p>
					<p class="username">{user.name}</p>
				</article>
			{/if}
		</button>
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

	.btn_wrapper {
		position: relative;
		background: none;
		border: none;
		padding: 0;
		margin: 0;
	}

	#logged_in_popup {
		position: absolute;
		top: 80%;
		right: 10%;
		padding: 10px 20px;
		width: 100px;
		background-color: rgba(0, 0, 0, 0.7);
		display: flex;
		flex-direction: column;
		align-items: baseline;
		gap: 2px;
	}

	.username {
		font-weight: 600;
	}

	#logged_in_popup p {
		font-size: 11pt;
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

	.profile-pic {
		border-radius: 50%;
		border: 1px solid white;
		width: 40px;
		height: 40px;
	}
</style>
