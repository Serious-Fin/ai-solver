<script lang="ts">
	import { PUBLIC_GOOGLE_OAUTH_CLIENT_ID } from '$env/static/public'
	import { page } from '$app/state'
	const redirectTo = page.url.searchParams.get('redirectTo')
</script>

<section>
	<article class="login_box">
		<article class="problem_article_header">
			<h1 class="inter problem_article_header_text">Login</h1>
			<p class="inter login_box_p">via any of these accounts</p>
		</article>
		<article class="login_buttons">
			<div>
				<div
					id="g_id_onload"
					data-client_id={PUBLIC_GOOGLE_OAUTH_CLIENT_ID}
					data-context="signin"
					data-login_uri="http://localhost:5173/api/oauth/google?redirectTo={redirectTo}"
				></div>
				<div
					class="g_id_signin"
					data-type="standard"
					data-size="medium"
					data-theme="outline"
					data-text="sign_in_with"
					data-shape="rectangular"
					data-logo_alignment="center"
				></div>
			</div>
			<div class="github_login">
				<img src="/github-mark.svg" alt="github logo" />Sign in with GitHub
			</div>
		</article>
	</article>
</section>

<style>
	section {
		background-color: var(--background);
		height: 100vh;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.login_box {
		background-color: var(--foreground);
		max-width: 400px;
		margin: 0 auto;
		box-sizing: border-box;
		padding: 40px 60px;
		border-radius: 7px;
	}

	.login_buttons {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.login_box_p {
		color: rgba(0, 0, 0, 0.75);
	}

	.github_login {
		background-color: white;
		padding: 6px 8px;
		margin-top: 8px;
		border-radius: 4px;
		border: 0.6px solid rgb(218, 220, 224);
		font-size: 14px;
		letter-spacing: 0.25px;
		font-family: 'Google Sans', arial, sans-serif;
		font-weight: 510;
		color: #3c4043f2;

		display: flex;
		align-items: center;
		gap: 6px;
	}

	.github_login img {
		width: 18px;
	}

	.github_login:active {
		background-color: rgb(236, 243, 254);
		border: 2px solid rgb(0, 99, 155);
	}
</style>
