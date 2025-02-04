<main>
    <div class="wrapper">
        <Card>
            <span slot="header">Export Server Transcripts</span>
            <div slot="content" class="content">
                <span>
                    Transcripts of tickets created in your server will be exported as
                    <a href="https://en.wikipedia.org/wiki/JSON" class="link">JSON</a> files. This is a machine-readable
                    format, meaning that they can easily be imported by other instances of the bot.

                    You must be the <b>owner</b> of the server to export transcripts. Administrator or other permissions
                    are not sufficient.
                </span>

                <span>In order to export your transcripts, please provide us with some additional details.</span>

                <form on:submit|preventDefault={createRequest}>
                    <GuildSelector onlyOwned bind:guildId />

                    <div class="button-wrapper">
                        <Button icon="fa-paper-plane" --font-size="1rem" --padding="5px 10px"
                                disabled={guildId === "" || guildId.length < 17 || guildId.length > 21}>Submit</Button>
                    </div>
                </form>
            </div>
        </Card>
    </div>
</main>

<style>
    main {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        padding: 3% 0;
    }

    .wrapper {
        width: 50%;
        min-width: 600px;
        max-width: 95%;
    }

    .content {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        padding-bottom: 3px;
    }

    form {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    @media screen and (max-width: 1000px) {
        .wrapper {
            min-width: unset;
            width: 90%;
        }
    }

    .button-wrapper {
        display: flex;
        justify-content: flex-end;
    }
</style>

<script>
    import Card from "$lib/components/Card.svelte";
    import GuildSelector from "$lib/includes/GuildSelector.svelte";
    import Button from "$lib/components/Button.svelte";
    import {goto} from "$app/navigation";
    import {client} from "$lib/axios.js";

    let guildId = "";

    async function createRequest() {
      const res = await client.post('/requests', {
        request_type: "guild_transcripts",
        guild_id: guildId
      });

      if (res.status === 201) {
        goto("/app?request_created=true");
      } else {
        alert(res.data.error || "Unknown error occurred.");
      }
    }
</script>