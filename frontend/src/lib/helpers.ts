import { toast } from 'svelte-sonner';
import { browser } from '$app/environment';

export function handleError(msgToUser: string, err: Error) {
    toast.error(msgToUser)

    if (browser) {
        fetch('/api/log-error', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                message: err.message,
            })
        }).catch(logErr => {
            console.error('Failed to log error to server:', logErr);
        });
    }
}

