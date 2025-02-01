<main>
    <a class="sign-in" href="{OAUTH_URL}">
        <img src="/img/discord-mark-black.png" class="discord-logo" alt="Discord logo"/>
        <span>Sign in with Discord</span>
    </a>

    <span class="explanation">
        Why do I need to sign in?
        <Tooltip>
            <span slot="content">
                <i class="fas fa-question-circle explanation-icon"></i>
            </span>

            <span slot="tip">
                You need to sign in so that we can verify your identity before giving you access to your data.

                When you sign in to this service, all session data is stored locally in your browser, and is never stored on our server.
            </span>
        </Tooltip>
    </span>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 50px;
        padding-top: 10%;
    }

    .sign-in {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px;
        user-select: none;
    }

    .sign-in > * {
        cursor: pointer;
    }

    .discord-logo {
        width: 50px;
    }

    .explanation {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 5px;
        opacity: 85%;
        font-size: 0.8em;
    }

    .explanation-icon {
        cursor: pointer;
    }
</style>

<script>
    import {
        PUBLIC_ROOT_URL,
        PUBLIC_APPLICATION_ID
    } from "$env/static/public";

    import Tooltip from "$lib/components/Tooltip.svelte";

    const REDIRECT_URI = encodeURIComponent(`${PUBLIC_ROOT_URL}/auth/callback`);
    const OAUTH_URL = `https://discord.com/oauth2/authorize?response_type=code&redirect_uri=${REDIRECT_URI}&scope=identify+guilds&client_id=${PUBLIC_APPLICATION_ID}`;
</script>