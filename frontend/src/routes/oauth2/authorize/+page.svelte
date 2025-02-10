<main>
  {#if error}
    <h2>Error</h2>
    <span>{error}</span>
  {:else if ready}
    <h2>Grant Data Access</h2>
    <p>Grant access to <b>{oauthClient.label}</b> to export your data</p>

    <p>This will allow the developer of <b>{oauthClient.label}</b> to:</p>
    <ul>
      {#each scopes as scope}
        <li>{VALID_SCOPES[scope]}</li>
      {/each}
    </ul>

    <p>For the server <b>{guild.name}</b> (<i>{guild.id}</i>)</p>

    <Button on:click={authorize}>Authorize</Button>
  {/if}
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding-top: 10%;
    font-size: 1.25rem;
    text-align: center;
  }
</style>

<script>
  import {page} from '$app/state';
  import {onMount} from "svelte";
  import {client} from "$lib/axios.js";
  import Button from "$lib/components/Button.svelte";

  const VALID_SCOPES = {
    'guild_transcripts': 'Generate & download transcript exports',
    'guild_data': 'Generate & download server data exports',
  };

  const getQueryParam = (key) => page.url.searchParams.get(key);

  let error;
  let ready = false;

  const clientId = getQueryParam('client_id');
  const guildId = getQueryParam('guild_id');
  const scopes = getQueryParam('scopes')?.split(' ');
  const redirectUri = getQueryParam('redirect_uri');
  const state = getQueryParam('state');

  let oauthClient;
  let guild;

  async function authorize() {
    const newScopes = scopes.map(s => `${guildId}:${s}`).join(' ');

    const res = await client.post(`/oauth/authorize?response_type=code&client_id=${clientId}&scope=${newScopes}&redirect_uri=${redirectUri}`);
    if (res.status !== 200) {
      error = res.data.error;
      return;
    }

    const url = new URL(redirectUri);
    url.search = '';
    url.searchParams.set('code', res.data.code);
    if (state) url.searchParams.set('state', state);

    window.location.href = url.toString();
  }

  async function load() {
    await client.get(`/auth/check_token`);

    if (!clientId) {
      error = 'Missing client_id parameter';
      return;
    }

    if (!guildId) {
      error = 'Missing guild_id parameter';
      return;
    }

    if (!redirectUri) {
      error = 'Missing redirect_uri parameter';
      return;
    }

    if (!scopes || scopes.length === 0 || scopes.find(s => !VALID_SCOPES[s])) {
      error = 'Invalid or missing scopes parameter';
      return;
    }

    const guilds = window.localStorage.getItem("guilds");
    if (!guilds) {
      error = 'No guilds found';
      return;
    }

    guild = JSON.parse(guilds).find(g => g.id === guildId);
    if (!guild) {
      error = 'Server not found: you can only export data from servers which you are the owner of';
      return;
    }

    const res = await client.get(`/oauth2/client/${clientId}`);
    if (res.status !== 200) {
      error = res.data.error;
      return;
    }

    oauthClient = res.data;
    ready = true;
  }

  onMount(load);
</script>