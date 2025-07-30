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

const arrayStartRegex = /\[\][a-zA-Z0-9_]*{/g
const arrayEndRegex = /}/g
export function printHumanReadable(input: string): string {
    input = input.replaceAll(arrayStartRegex, "[")
    input = input.replaceAll(arrayEndRegex, "]")
    return input
}
