<main>
  <Card>
    <span slot="header">Public Key</span>
    <div slot="content" class="description">
      <span>
        The ticketsbot.net public key is used to sign any data exported. For users who just want to access their
        data, this does not matter. However, for people hosting their own instance of the bot, the public key
        can be used to verify the authenticity of data imported by users. In other words, it can be used to
        ensure that data was truly exported from the official instance of the bot, and has not been tampered
        with after the export.
      </span>

      <span>
        We have also provided some Go libraries for validating data, which can be found on Github.
        <ul>
          <li><a href="https://github.com/TicketsBot/export/tree/master/pkg/validator" class="link">
            github.com/TicketsBot/export/pkg/validator</a></li>
          <li><a href="https://github.com/TicketsBot/export/tree/master/example" class="link">
            github.com/TicketsBot/export/example</a></li>
        </ul>
      </span>

      <div>
        <span>
        The public key is an Ed25519 key, and in PEM format is as follows:
        </span>

        <pre><code>{publicKey}</code></pre>
      </div>
    </div>
  </Card>
</main>

<style>
  main {
    display: flex;
    flex-direction: row;
    gap: 2%;
    padding: 3%;
  }

  .description {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  pre {
    background-color: var(--text);
    color: var(--bg);

    padding: 0.5em;
    border-radius: 0.25em;
    white-space: pre-wrap;
  }
</style>

<script>
  import Card from "../../../lib/components/Card.svelte";
  import {onMount} from "svelte";
  import {client} from "$lib/axios.js";

  let publicKey = "Loading...";

  async function loadPublicKey() {
    const res = await client.get("/keys/signing");
    publicKey = res.data;
  }

  onMount(() => {
    loadPublicKey();
  })
</script>
