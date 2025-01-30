<main>
    <div class="new-task">
        <Card>
            <span slot="header">Export / Remove Data</span>
            <span slot="content" class="description">
                <div>
                    <span>
                        You can request to export or remove your data from our servers at any time. Please select an
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
                                {#if request.status === "complete"}
                                    <div class="status-row">
                                        <span class="status complete">Complete</span>

                                        {#if request.is_export}
                                            {#if new Date(request.download_expires) > new Date()}
                                                <span><a href="{request.download_url}" class="download">Download</a>
                                                (Expires in {formatExpiry(new Date(request.download_expires))})</span>
                                            {:else}
                                                <span>Download link has expired</span>
                                            {/if}
                                        {/if}
                                    </div>
                                {:else}
                                    <span class="status in-progress">In Progress</span>
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

    .download {
        color: #3472f7;
        cursor: pointer;
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
    import {PUBLIC_BACKEND_URI} from "$env/static/public";
    import {onMount} from "svelte";

    const REQUEST_NAMES = {
      guild_transcripts: "Export Server Transcripts"
    };

    let requests = [];

    async function loadRequests(code) {
      const res = await fetch(`${PUBLIC_BACKEND_URI}/requests`, {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${localStorage.getItem("token")}`
        }
      });
      if (res.ok) {
        requests = await res.json();
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
