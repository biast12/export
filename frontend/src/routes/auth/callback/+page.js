import {error} from "@sveltejs/kit";

export function load({ params, url }) {
    const code = url.searchParams.get('code');
    if (!code) {
        error(400, 'Missing code: try signing in again');
    }

    const state = url.searchParams.get('state');

    return { code, state };
}