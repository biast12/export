<div class="item">
    <label for="needs-selector">Server Selector</label>
    <select name="needs-selector" bind:value={needsSelector} on:input={() => guildId = ""}>
        <option value="knows_id">I know the numeric server ID of the server I want to export my data from
        </option>
        <option value="wants_list" selected>I do not know the numeric server ID</option>
    </select>
</div>

<div class="item">
    {#if needsSelector === "knows_id"}
        <label for="server-id">Numeric Server ID</label>
        <input type="text" name="server-id" id="server-id" placeholder="Numeric Server ID" required
               on:input={handleServerIdTextUpdate} on:keypress={handleServerIdKeyPress}>

        <span class="validation-error" class:hide={guildId.length === 0 || (guildId.length >= 17 && guildId.length <= 21)}>
            Server ID must be between 17-21 digits.
        </span>
    {:else if needsSelector === "wants_list"}
        {#if guilds === null}
            <SignInButton scopes="identify guilds" state="guilds"/>
        {:else}
            <label for="guild-list">Server List</label>
            <div class="guild-list">
                {#each guilds.filter(g => !onlyOwned || g.owner === true) as guild}
                    {@const fileExtension = guild.icon?.startsWith("a_") ? "gif" : "webp"}
                    <div class="guild" class:active={guildId === guild.id} on:click={() => setActive(guild)}>
                        {#if guild.icon === null}
                            <img src="https://cdn.discordapp.com/embed/avatars/0.png" alt="Server Icon">
                        {:else}
                            <img src="https://cdn.discordapp.com/icons/{guild.id}/{guild.icon}.{fileExtension}" alt="Server Icon">
                        {/if}

                        <span>{guild.name}</span>
                    </div>
                {/each}
            </div>
        {/if}
    {/if}
</div>

<style>
    .item {
        display: flex;
        flex-direction: column;
        gap: 0;
    }

    .validation-error {
        font-size: 0.8rem;
        color: #e34242;
        opacity: 0.8;
    }

    .hide {
        display: none;
    }

    .guild-list {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        justify-content: space-evenly;
        row-gap: 6px;
    }

    .guild {
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 32%;
        box-sizing: content-box;
        cursor: pointer;
        border: 2px solid transparent;
        border-radius: 2px;
        user-select: none;
        box-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
    }

    .guild.active {
        border: 2px solid var(--primary);
    }

    .guild > img {
        width: 64px;
        height: 64px;
        border-radius: 50%;
        padding: 5px;
    }

    @media screen and (max-width: 768px) {
        .guild {
            width: 48%;
        }
    }

    @media screen and (max-width: 576px) {
        .guild {
            width: 100%;
        }
    }
</style>

<script>
    import SignInButton from "./SignInButton.svelte";
    import {onMount} from "svelte";

    export let guildId = "";
    export let onlyOwned = false;

    let guilds = null;

    let needsSelector = "wants_list";

    function handleServerIdKeyPress(e) {
        if (!/^\d$/.test(e.key)) {
            e.preventDefault();
        }
    }

    function handleServerIdTextUpdate(e) {
        guildId = e.target.value;
    }

    function setActive(guild) {
        guildId = guild.id;
    }

    onMount(() => {
        const storedGuilds = window.localStorage.getItem("guilds");
        if (guilds !== undefined) {
            guilds = JSON.parse(storedGuilds);
        }

        window.addEventListener("storage", (event) => {
            if (event.key === "guilds") {
                guilds = JSON.parse(event.newValue);
            }
        });
    })
</script>