<script lang="ts">
	import {
		PUBLIC_FRONTEND_BASE_URL_BROWSER,
		PUBLIC_GOOGLE_OAUTH_CLIENT_ID
	} from '$env/static/public'
	import { page } from '$app/state'
	const redirectTo = page.url.searchParams.get('redirectTo') ?? '/'

	const loginViaGithub = () => {
		window.location.href = `/api/oauth/github?redirectTo=${encodeURIComponent(redirectTo)}`
	}
</script>

<section>
	<article class="page-title">
		<a class="btn" href={redirectTo}>
			<img class="img_icon" src="/arrow_back.svg" alt="return back arrow" />
		</a>
		<h1 class="inter">Guide The AI</h1>
	</article>
	<article class="login_box">
		<article class="problem_article_header">
			<h2 class="inter problem_article_header_text">Login</h2>
			<p class="inter login_box_p">via any of these accounts</p>
		</article>
		<article class="login_buttons">
			<div>
				<div
					id="g_id_onload"
					data-client_id={PUBLIC_GOOGLE_OAUTH_CLIENT_ID}
					data-context="signin"
					data-login_uri="{PUBLIC_FRONTEND_BASE_URL_BROWSER}/api/oauth/google?redirectTo={redirectTo}"
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
			<button class="github_login" onclick={loginViaGithub}>
				<img src="/github-mark.svg" alt="github logo" />Sign in with GitHub
			</button>
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
		flex-direction: column;
	}

	.page-title {
		display: flex;
		gap: 10px;
		margin-bottom: 20px;
		align-items: center;
		justify-content: baseline;
	}

	h1 {
		color: white;
		font-size: 22pt;
	}

	.btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.img_icon {
		width: 100%;
		height: 100%;
	}

	.login_box {
		background-color: var(--foreground);
		max-width: 400px;
		margin: 0 auto;
		box-sizing: border-box;
		padding: 40px 40px;
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
		padding: 6px 10px;
		margin-top: 8px;
		border-radius: 4px;
		border: 0.6px solid rgb(218, 220, 224);
		font-size: 14px;
		letter-spacing: 0.25px;
		font-family: 'Google Sans', arial, sans-serif;
		font-weight: 500;
		color: #3c4043f2;

		display: flex;
		align-items: center;
		gap: 10px;
	}

	.github_login img {
		width: 18px;
	}

	.github_login:active {
		background-color: rgb(236, 243, 254);
		border: 2px solid rgb(0, 99, 155);
	}

	@media (min-width: 768px) {
		.login_box {
			padding: 40px 60px;
		}
	}
</style>
