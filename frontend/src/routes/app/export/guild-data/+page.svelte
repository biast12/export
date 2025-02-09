<main>
    <div class="wrapper">
        <Card>
            <span slot="header">Export Server Data</span>
            <div slot="content" class="content">
                <span>
                    All database data associated to your server (e.g. tickets, ticket panels, etc.), excluding
                    transcripts will be exported. This data will be provided to you in a JSON format. These file can
                    then potentially be used to import the data into another bot instance.

                    You must be the <b>owner</b> of the server to export data. Administrator or other permissions are
                    not sufficient.
                </span>

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
        request_type: "guild_data",
        guild_id: guildId
      });

      if (res.status === 201) {
        goto("/app?request_created=true");
      } else {
        alert(res.data.error || "Unknown error occurred.");
      }
    }
</script>