<button on:click {disabled}>
    {#if icon}
        <i class="fas {icon}" />
    {/if}

    <slot />
</button>

<style>
    button {
        display: flex;

        justify-content: center;
        align-items: center;

        gap: 4px;

        text-align: center;

        color: white;
        background-color: var(--primary);
        border-color: var(--primary);
        outline: none;
        border-style: solid;
        border-width: 2px;
        border-radius: .25rem;
        font-family: 'Poppins', sans-serif;
        padding: var(--padding, 2px 5px);

        font-size: var(--font-size, 1rem);

        cursor: pointer;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.2);
    }

    button:not(:disabled):active {
        background-color: #0062cc;
        border-color: #0062cc;
        transform: translateY(1px);
    }

    button:disabled {
        background-color: #6c757d;
        border-color: #6c757d;
        cursor: default;
    }
</style>

<script>
    export let icon;
    export let disabled = false;
</script>