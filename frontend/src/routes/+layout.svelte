<script lang="ts">
	import { Toaster } from 'svelte-sonner'
	import '../global.css'
</script>

<Toaster position="bottom-center" richColors duration={2000} />
<slot />
