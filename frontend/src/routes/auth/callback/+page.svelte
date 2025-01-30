<main>
    {#if !error}
        <span>Signing you in...</span>
    {:else}
        <span>{error}</span>
    {/if}
</main>

<style>
    main {
        display: flex;
        justify-content: center;
        align-items: center;
        padding-top: 10%;
        font-size: 1.25rem;
    }
</style>

<script>
    import {onMount} from "svelte";

    import {PUBLIC_BACKEND_URI} from "$env/static/public";
    import {goto} from "$app/navigation";

    export let data;

    let error;

    async function exchange(code) {
        const res = await fetch(`${PUBLIC_BACKEND_URI}/auth/exchange`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({code})
        });

        if (res.ok) {
            const {token} = await res.json();
            localStorage.setItem("token", token);

            goto("/app");
        } else {
            console.error("Failed to exchange code for token");

            const {error} = await res.json();
            alert(error);
        }
    }

    async function loadGuilds(code) {
        const res = await fetch(`${PUBLIC_BACKEND_URI}/auth/guilds`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({code})
        });

        if (res.ok) {
            const guilds = await res.json();
            localStorage.setItem("guilds", JSON.stringify(guilds));

            return true;
        } else {
            console.error("Failed to load guilds");

            const {err} = await res.json();
            error = err;

            return false;
        }
    }

    onMount(async () => {
        if (data.state === "guilds") {
            if (await loadGuilds(data.code)) {
                window.close();
            }
        } else {
            await exchange(data.code);
        }
    });
</script>