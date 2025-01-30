<section>
    <slot name="content">
        <i class="fas fa-question-circle"></i>
    </slot>

    <div class="tip">
        <slot name="tip"></slot>
    </div>
</section>

<style>
    section > .tip {
        position: absolute;
        z-index: 1;

        visibility: hidden;

        width: max-content;
        max-width: 25%;
        border-radius: 3px;
        padding: 5px 10px;

        word-wrap: break-word;
        white-space: pre-line;

        background-color: #0f0f0f;
        color: var(--bg);
        text-align: center;
    }

    section:hover > .tip, .tip:hover {
        visibility: visible;
    }
</style>