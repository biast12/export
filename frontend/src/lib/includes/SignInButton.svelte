<div class="sign-in" on:click={open}>
    <img src="/img/discord-mark-black.png" class="discord-logo" alt="Discord logo"/>
    <span>Sign in with Discord</span>
</div>

<style>
    div {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px;
        user-select: none;
    }

    div > * {
        cursor: pointer;
    }

    .discord-logo {
        width: 50px;
    }
</style>

<script>
    import {PUBLIC_APPLICATION_ID, PUBLIC_ROOT_URL} from "$env/static/public";

    export let scopes = "identify";
    export let state = ";"

    const REDIRECT_URI = encodeURIComponent(`${PUBLIC_ROOT_URL}/auth/callback`);
    const OAUTH_URL = `https://discord.com/oauth2/authorize?response_type=code&redirect_uri=${REDIRECT_URI}&scope=${encodeURIComponent(scopes)}&client_id=${PUBLIC_APPLICATION_ID}&state=${state}`;

    function open() {
        window.open(OAUTH_URL, "newwindow", "width=500, height=700");
    }
</script>
