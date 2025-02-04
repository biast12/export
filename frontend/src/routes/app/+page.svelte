<main>
    <div class="new-task">
        <Card>
            <span slot="header">Export Data</span>
            <span slot="content" class="description">
                <div>
                    <span>
                        You can request to export your data from our servers at any time. Please select an
                        action below to get started.
                    </span>

                    <select class="action-selector" on:input={navigate}>
                        <option disabled selected>Select a data action...</option>
                        <option class="header" disabled>Transcripts</option>
                        <option value="export-guild">Export Server Transcripts</option>
                    </select>
                </div>
            </span>
        </Card>
    </div>

    <div class="previous-tasks">
        <Card>
            <span slot="header">Previous Requests</span>
            <div slot="content">
                {#each requests as request}
                    <div class="task">
                        <div>
                            <div>
                                <span class="name">{REQUEST_NAMES[request.type]}</span>
                                {#if request.status === "completed"}
                                    <div class="status-row">
                                        <span class="status complete">Complete</span>

                                      {#if new Date(request.artifact_expires_at) > new Date()}
                                                <span>
                                                  <a href="{request.download_url}" class="download"class:downloading={downloadingArtifacts[request.id]}
                                                     on:click={() => downloadArtifact(request.id)}>
                                                    {#if downloadingArtifacts[request.id]}
                                                      Downloading...
                                                    {:else}
                                                      <i class="fa-solid fa-download"></i>
                                                      Download
                                                    {/if}
                                                  </a>
                                                  (Expires in {formatExpiry(new Date(request.artifact_expires_at))})
                                                </span>
                                      {:else}
                                        <span>Download link has expired</span>
                                      {/if}
                                    </div>
                                {:else if request.status === "queued"}
                                    <span class="status in-progress">In Progress</span>
                                {:else if request.status === "failed"}
                                    <span class="status failed">Failed</span>
                                {:else}
                                  <span class="status">Unknown</span>
                                {/if}

                                {#if request.guild_id}
                                    <span>Server ID: <em>{request.guild_id}</em></span>
                                {/if}

                                <span>Requested on {new Date(request.created_at).toLocaleDateString()}</span>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </Card>
    </div>
</main>

<style>
    main {
        display: flex;
        flex-direction: row;
        gap: 2%;
        padding: 3%;
    }

    .new-task {
        flex: 1;
    }

    .previous-tasks {
        width: 30%;
    }

    .action-selector {
        width: 30%;
        min-width: 500px;
    }

    .description {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .description > *:first-child {
        display: flex;
        flex-direction: column;
    }

    .task {
        display: flex;
        flex-direction: column;
        user-select: none;
    }

    .task > div {
        display: flex;
        flex-direction: row;
    }

    .task > div > div {
        display: flex;
        flex-direction: column;
    }

    .task > div > div:first-child {
        flex: 1;
    }

    .task > div > div:last-child {
        justify-content: flex-end;
        margin-bottom: 3px;
    }

    .task:not(:first-child) {
        margin-top: 1px;
    }

    .task:not(:last-child) {
        margin-bottom: 2px;
    }

    .task:not(:last-child)::after {
        content: "";
        display: block;
        width: 100%;
        height: 1px;
        background-color: var(--text);
        opacity: 0.5;
    }

    .task .name {
        font-weight: 500;
    }

    .task .status-row {
        display: flex;
        flex-direction: row;
        gap: 5px;
    }

    .task .status-row > *:not(:last-child)::after {
        content: "-";
        color: var(--text);
        margin-left: 5px;
        opacity: 0.5;
    }

    .status.complete {
        color: green;
    }

    .status.in-progress {
        color: darkorange;
    }

    .status.failed {
      color: darkred;;
    }

    .download {
        color: #3472f7;
        cursor: pointer;
    }

    .download.downloading {
        color: #727272;
        cursor: not-allowed;
    }

    @media screen and (max-width: 1000px) {
        main {
            flex-direction: column;
            gap: 2rem;
        }

        .previous-tasks {
            width: 100%;
        }
    }

    @media screen and (max-width: 800px) {
        .action-selector {
            width: 100%;
            min-width: unset;
        }
    }
</style>

<script>
    import Card from "$lib/components/Card.svelte";
    import {goto} from "$app/navigation";
    import {onMount} from "svelte";
    import {client} from "$lib/axios.js";

    const REQUEST_NAMES = {
      guild_transcripts: "Export Server Transcripts"
    };

    let requests = [];
    let downloadingArtifacts = {};

    async function loadRequests() {
      const res = await client.get("/requests");
      if (res.status === 200) {
        requests = res.data;
      }
    }

    async function downloadArtifact(requestId) {
      if (downloadingArtifacts[requestId]) {
        return;
      }

      downloadingArtifacts[requestId] = true;

      const res = await client({
        url: `/requests/${requestId}/artifact`,
        method: "GET",
        responseType: "blob"
      });

      if (res.status === 200) {
        const href = URL.createObjectURL(res.data);
        const link = document.createElement('a');
        link.href = href;
        link.setAttribute('download', `export-${requestId}.zip`);

        document.body.appendChild(link);
        link.click();

        document.body.removeChild(link);
        URL.revokeObjectURL(href);
      } else {
        const json = await res.json();
        alert(json.error);
      }
    }

    function navigate(e) {
        switch (e.target.value) {
            case "export-guild":
                goto("/app/export/guild-transcripts");
                break;
        }
    }

    function formatExpiry(date) {
        const interval = date - new Date();

        const days = Math.floor(interval / (1000 * 60 * 60 * 24));
        if (days > 0) {
            return `${days} day${days > 1 ? "s" : ""}`;
        }

        const hours = Math.floor(interval / (1000 * 60 * 60));
        if (hours > 0) {
            return `${hours} hour${hours > 1 ? "s" : ""}`;
        }

        const minutes = Math.floor(interval / (1000 * 60));
        return `${minutes} minute${minutes > 1 ? "s" : ""}`;
    }

    onMount(() => {
        loadRequests();
    });
</script>
