<main>
    <span>Signing you out...</span>
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

    onMount(() => {
        window.localStorage.clear();

        setTimeout(() => {
            window.location.href = 'https://ticketsbot.net';
        }, 1000);
    })
</script>