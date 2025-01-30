<slot/>

<svelte:head>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.2/css/all.min.css" integrity="sha512-SnH5WK+bZxgPHs44uWIX+LLJAJ9/2PkPKZ5QiAj6Ta86w+fsb2TkcmfRyVX3pBnMFcV7oQPJkl9QevSCWr3W6A==" crossorigin="anonymous" referrerpolicy="no-referrer" />
</svelte:head>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');

    :root {
        --bg: #f0f0f0;
        --text: #303030;
        --primary: #3472f7;
        --form-border: rgb(118, 118, 118);
    }

    :global(body) {
        background-color: var(--bg);
        margin: 0;
        padding: 0;
        height: 100%;

        font-family: 'Poppins', sans-serif;
        color: var(--text);
    }

    :global(a:link, a:visited, a:hover, a:active) {
        text-decoration: none;
        color: var(--text);
    }

    :global(a.link:link, a.link:visited, a.link:hover, a.link:active) {
        color: var(--primary);
    }

    :global(select, input) {
        padding: 0.3rem;
        border: 1px solid var(--form-border);
        border-radius: 0.25rem;
        font-size: 1rem;
    }

    :global(select:focus, input:focus) {
        outline: none;
    }

    :global(option.header) {
        font-weight: bold;
    }

    :global(label) {
        font-size: 0.8rem;
        opacity: 0.85;
    }
</style>