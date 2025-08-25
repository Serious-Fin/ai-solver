<script lang="ts">
	import { Toaster } from 'svelte-sonner'
	import '../global.css'
</script>

<Toaster position="top-center" richColors duration={3600} />
<slot />
