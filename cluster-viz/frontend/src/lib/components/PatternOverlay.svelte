<script lang="ts">
  import type { ScenarioMetadata, ScenarioPatterns, Technology } from '../types/scenario';

  export let visible: boolean = true;
  export let metadata: ScenarioMetadata;
  export let patterns: ScenarioPatterns;
  export let technologies: Technology[];

  let dontShowAgain = false;

  function handleClose() {
    visible = false;
    if (dontShowAgain) {
      localStorage.setItem('hidePatternOverlay', 'true');
    }
  }

  function handleStartExploring() {
    handleClose();
  }
</script>

{#if visible}
  <div class="overlay-backdrop" on:click={handleClose} role="button" tabindex="0">
    <div class="overlay-content" on:click|stopPropagation role="dialog">
      <div class="overlay-header">
        <h1 class="overlay-title">{metadata.title}</h1>
        <button class="close-button" on:click={handleClose} aria-label="Close">
          ✕
        </button>
      </div>

      <div class="overlay-body">
        <p class="description">{metadata.description}</p>

        <div class="section">
          <h2 class="section-title">📚 Patterns You'll Learn</h2>
          <div class="patterns-grid">
            {#each patterns.primary as pattern}
              <div class="pattern-card primary">
                <div class="pattern-header">
                  <span class="pattern-icon">{pattern.icon}</span>
                  <h3 class="pattern-name">{pattern.name}</h3>
                </div>
                <p class="pattern-description">{pattern.description}</p>
                {#if pattern.benefits && pattern.benefits.length > 0}
                  <ul class="pattern-benefits">
                    {#each pattern.benefits as benefit}
                      <li>{benefit}</li>
                    {/each}
                  </ul>
                {/if}
              </div>
            {/each}
            {#each patterns.secondary as pattern}
              <div class="pattern-card secondary">
                <div class="pattern-header">
                  <span class="pattern-icon">{pattern.icon}</span>
                  <h3 class="pattern-name">{pattern.name}</h3>
                </div>
                <p class="pattern-description">{pattern.description}</p>
                {#if pattern.benefits && pattern.benefits.length > 0}
                  <ul class="pattern-benefits">
                    {#each pattern.benefits as benefit}
                      <li>{benefit}</li>
                    {/each}
                  </ul>
                {/if}
              </div>
            {/each}
          </div>
        </div>

        <div class="section">
          <h2 class="section-title">🛠️ Technologies</h2>
          <div class="tech-grid">
            {#each technologies as tech}
              <div class="tech-item">
                <strong>{tech.name}</strong>
                <span class="tech-role">{tech.role}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>

      <div class="overlay-footer">
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={dontShowAgain} />
          Don't show this again
        </label>
        <button class="start-button" on:click={handleStartExploring}>
          Start Exploring →
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .overlay-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 2rem;
  }

  .overlay-content {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
    border: 1px solid rgba(71, 85, 105, 0.5);
    border-radius: 16px;
    max-width: 800px;
    width: 100%;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  }

  .overlay-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 2rem 2rem 1rem;
    border-bottom: 1px solid rgba(71, 85, 105, 0.3);
  }

  .overlay-title {
    font-size: 1.75rem;
    font-weight: 700;
    color: #e2e8f0;
    margin: 0;
    background: linear-gradient(135deg, #60a5fa 0%, #a78bfa 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .close-button {
    background: none;
    border: none;
    color: #94a3b8;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0.5rem;
    line-height: 1;
    transition: color 0.2s;
  }

  .close-button:hover {
    color: #e2e8f0;
  }

  .overlay-body {
    padding: 2rem;
  }

  .description {
    font-size: 1.0625rem;
    line-height: 1.6;
    color: #cbd5e1;
    margin-bottom: 2rem;
  }

  .section {
    margin-bottom: 2rem;
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: #e2e8f0;
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .patterns-grid {
    display: grid;
    gap: 1rem;
  }

  .pattern-card {
    background: rgba(30, 41, 59, 0.6);
    border: 1px solid rgba(71, 85, 105, 0.5);
    border-radius: 12px;
    padding: 1.25rem;
    transition: all 0.3s;
  }

  .pattern-card.primary {
    border-left: 3px solid #f97316;
  }

  .pattern-card.secondary {
    border-left: 3px solid #a855f7;
  }

  .pattern-card:hover {
    border-color: rgba(96, 165, 250, 0.5);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  .pattern-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
  }

  .pattern-icon {
    font-size: 1.5rem;
  }

  .pattern-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #e2e8f0;
    margin: 0;
  }

  .pattern-description {
    font-size: 0.9375rem;
    line-height: 1.5;
    color: #cbd5e1;
    margin-bottom: 0.75rem;
  }

  .pattern-benefits {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .pattern-benefits li {
    font-size: 0.875rem;
    color: #94a3b8;
    padding-left: 1.25rem;
    position: relative;
  }

  .pattern-benefits li::before {
    content: "✓";
    position: absolute;
    left: 0;
    color: #10b981;
  }

  .tech-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 0.75rem;
  }

  .tech-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding: 0.75rem 1rem;
    background: rgba(30, 41, 59, 0.4);
    border: 1px solid rgba(71, 85, 105, 0.3);
    border-radius: 8px;
    font-size: 0.9375rem;
  }

  .tech-item strong {
    color: #e2e8f0;
  }

  .tech-role {
    color: #94a3b8;
    font-size: 0.8125rem;
  }

  .overlay-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem 2rem;
    border-top: 1px solid rgba(71, 85, 105, 0.3);
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: #94a3b8;
    cursor: pointer;
  }

  .checkbox-label input[type="checkbox"] {
    cursor: pointer;
  }

  .start-button {
    background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
    color: white;
    border: none;
    border-radius: 8px;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  }

  .start-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
  }

  /* Scrollbar styling */
  .overlay-content::-webkit-scrollbar {
    width: 8px;
  }

  .overlay-content::-webkit-scrollbar-track {
    background: rgba(15, 23, 42, 0.5);
  }

  .overlay-content::-webkit-scrollbar-thumb {
    background: rgba(71, 85, 105, 0.5);
    border-radius: 4px;
  }

  .overlay-content::-webkit-scrollbar-thumb:hover {
    background: rgba(71, 85, 105, 0.8);
  }
</style>
