<section>
    <header>
        <div>
            <slot name="header" />
        </div>
    </header>
    <div class="content">
        <slot name="content" />
    </div>
</section>

<style>
    section {
        display: flex;
        flex-direction: column;

        width: 100%;

        /* box shadow, but none on top */
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2), 0 0 4px rgba(0, 0, 0, 0.2);
        border-radius: 5px;
    }

    header {
        font-weight: 500;
        font-size: 1.5em;
        user-select: none;
    }

    header::after {
        content: "";
        display: block;
        width: 100%;
        height: 1px;
        background-color: var(--text);
        opacity: 0.3;
    }

    header > div {
        padding: 0 8px;
    }

    .content {
        padding: 5px 8px;
    }
</style>