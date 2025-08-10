<section>
	<article class="login_box">
		<article class="problem_article_header">
			<h1 class="inter problem_article_header_text">Login</h1>
			<p class="inter login_box_p">via any of these accounts</p>
		</article>
		<article>
			<div>
				<div
					id="g_id_onload"
					data-client_id="YOUR_GOOGLE_CLIENT_ID"
					data-login_uri="https://your.domain/your_login_endpoint"
					data-auto_prompt="false"
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
		max-width: 350px;
		margin: 0 auto;
		box-sizing: border-box;
		padding: 30px 40px;
		border-radius: 7px;
	}

	.login_box_p {
		color: rgba(0, 0, 0, 0.75);
	}
</style>
