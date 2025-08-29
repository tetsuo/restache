const PlayExampleClassName = {
  PLAY_CONTAINER: '.js-exampleContainer',
  EXAMPLE_INPUT: '.Playground-Code',
  EXAMPLE_OUTPUT: '.Playground-Output',
  RUN_BUTTON: '.Playground-RunButton',
};

export class PlaygroundExampleController {
  constructor(exampleEl, run) {
    this.exampleEl = exampleEl;
    this.runButtonEl = exampleEl.querySelector(PlayExampleClassName.RUN_BUTTON);
    this.inputEl = this.makeTextArea(exampleEl.querySelector(PlayExampleClassName.EXAMPLE_INPUT));
    this.outputEl = exampleEl.querySelector(PlayExampleClassName.EXAMPLE_OUTPUT);
    this.run = run;

    if (!this.inputEl) {
      return;
    }

    this.resize();
    this.inputEl.addEventListener('keyup', () => this.resize());
    this.inputEl.addEventListener('keydown', e => this.onKeydown(e));

    this.runButtonEl?.addEventListener('click', () => this.handleRunButtonClick());

    if (exampleEl.hasAttribute('data-autorun')) {
      this.handleRunButtonClick();
    }
    // This needs to be recalculated when the font size changes on breakpoint
    const media = window.matchMedia('(min-width: 65rem)');
    media.addEventListener('change', () => this.resize());
  }

  makeTextArea(el) {
    const t = document.createElement('textarea');
    t.classList.add('Playground-Code', 'code');
    t.spellcheck = false;
    t.value = el?.textContent.trim() ?? '';
    el?.parentElement?.replaceChild(t, el);
    return t;
  }

  resize() {
    if (this.inputEl?.value) {
      const numLineBreaks = (this.inputEl.value.match(/\n/g) || []).length;
      this.inputEl.style.height = 'auto';
      // Force reflow to fix iOS Safari bug where scrollHeight isn't updated after
      // font-size changes (from Safari settings)
      void this.inputEl.offsetHeight;

      if (numLineBreaks >= 1) {
        this.inputEl.style.height = this.inputEl.scrollHeight + 2 + 'px';
        this.inputEl.style.overflow = "auto";
      } else {
        this.inputEl.style.height = "51px";
        this.inputEl.style.overflow = "hidden";
      }
    }
  }

  onKeydown(e) {
    if (e.key === 'Tab') {
      document.execCommand('insertText', false, '\t');
      e.preventDefault();
    }

    if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) {
      e.preventDefault();
      this.handleRunButtonClick();
    }

    if (e.key === 'Escape') {
      e.preventDefault();
      this.inputEl.blur();
    }
  }

  setOutputText(output) {
    if (this.outputEl) {
      this.outputEl.textContent = output;
    }
  }

  appendToOutputText(output) {
    if (this.outputEl) {
      this.outputEl.textContent += output;
    }
  }

  handleRunButtonClick() {
    const code = this.inputEl?.value ?? '';
    try {
      const result = this.run(code);
      if (result.error) {
        this.setOutputText('Error: ' + result.error);
      } else {
        this.setOutputText(result.output);
      }
    } catch (err) {
      this.setOutputText('Exception: ' + err.message);
    }
  }
}

export function initPlaygrounds(runner) {
  for (const el of document.querySelectorAll(PlayExampleClassName.PLAY_CONTAINER)) {
    new PlaygroundExampleController(el, runner);
  }
}
