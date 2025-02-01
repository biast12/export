<nav>
    <div>
        <img class="logo" src="/img/logo.png" alt="Logo" />
        <span class="brand">Data Management Self Service</span>
        <div class="vertical-line"></div>
        <a href="/app">Home</a>
        <div class="vertical-line"></div>
        <a href="/app/public-key">Download Public Key</a>
    </div>
    <div>
        <a href="/auth/sign-out">
            <i class="fas fa-sign-out" />
            Sign Out
        </a>
    </div>
</nav>

<slot />

<style>
    nav {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        align-items: center;
        padding: 5px 10px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        user-select: none;
    }

    nav > div {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 10px;
    }

    nav a, nav span {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 5px;
        font-size: 18px;
        user-select: auto;
    }

    .logo {
        width: 40px;
        height: 40px;
        vertical-align: middle;
        border-radius: 15%;
    }

    .brand {
        font-weight: 500;
    }

    .vertical-line {
        width: 1px;
        content: '';
        height: 22px;
        background-color: var(--text);
        opacity: 0.3;
    }

    @media screen and (max-width: 576px) {
        .brand, .vertical-line:first-of-type {
            display: none;
        }
    }
</style>