<script lang="ts">
  import { scenarios, currentScenarioId, switchScenario } from '../stores/scenario';
  import type { ScenarioSummary } from '../types/scenario';

  let isOpen = false;

  function toggleDropdown() {
    isOpen = !isOpen;
  }

  function handleSelect(scenario: ScenarioSummary) {
    switchScenario(scenario.id);
    isOpen = false;
  }

  function handleClickOutside(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.scenario-switcher')) {
      isOpen = false;
    }
  }

  $: currentScenario = $scenarios.find(s => s.id === $currentScenarioId);
</script>

<svelte:window on:click={handleClickOutside} />

<div class="scenario-switcher">
  <button class="switcher-button" on:click|stopPropagation={toggleDropdown}>
    <div class="button-content">
      <span class="button-label">Scenario:</span>
      <span class="button-value">{currentScenario?.title || 'Loading...'}</span>
    </div>
    <svg class="chevron" class:open={isOpen} width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
      <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
    </svg>
  </button>

  {#if isOpen}
    <div class="dropdown-menu">
      {#each $scenarios as scenario (scenario.id)}
        <button
          class="dropdown-item"
          class:active={scenario.id === $currentScenarioId}
          on:click|stopPropagation={() => handleSelect(scenario)}
        >
          <div class="item-content">
            <h4 class="item-title">{scenario.title}</h4>
            <p class="item-description">{scenario.description}</p>
            {#if scenario.patterns.length > 0}
              <div class="item-patterns">
                {#each scenario.patterns as pattern}
                  <span class="pattern-tag">{pattern}</span>
                {/each}
              </div>
            {/if}
          </div>
          {#if scenario.id === $currentScenarioId}
            <svg class="check-icon" width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          {/if}
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .scenario-switcher {
    position: relative;
  }

  .switcher-button {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: rgba(30, 41, 59, 0.8);
    border: 1px solid rgba(71, 85, 105, 0.5);
    border-radius: 8px;
    padding: 0.625rem 1rem;
    color: #e2e8f0;
    cursor: pointer;
    transition: all 0.2s;
    min-width: 300px;
  }

  .switcher-button:hover {
    border-color: rgba(96, 165, 250, 0.5);
    background: rgba(30, 41, 59, 0.95);
  }

  .button-content {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    flex: 1;
  }

  .button-label {
    font-size: 0.75rem;
    color: #94a3b8;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .button-value {
    font-size: 0.9375rem;
    font-weight: 600;
    color: #e2e8f0;
  }

  .chevron {
    color: #94a3b8;
    transition: transform 0.2s;
  }

  .chevron.open {
    transform: rotate(180deg);
  }

  .dropdown-menu {
    position: absolute;
    top: calc(100% + 0.5rem);
    left: 0;
    right: 0;
    background: rgba(30, 41, 59, 0.95);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(71, 85, 105, 0.5);
    border-radius: 12px;
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5);
    overflow: hidden;
    z-index: 100;
    max-height: 500px;
    overflow-y: auto;
  }

  .dropdown-item {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    width: 100%;
    padding: 1rem;
    background: none;
    border: none;
    border-bottom: 1px solid rgba(71, 85, 105, 0.2);
    color: #e2e8f0;
    cursor: pointer;
    text-align: left;
    transition: all 0.2s;
  }

  .dropdown-item:last-child {
    border-bottom: none;
  }

  .dropdown-item:hover {
    background: rgba(59, 130, 246, 0.1);
  }

  .dropdown-item.active {
    background: rgba(59, 130, 246, 0.15);
  }

  .item-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .item-title {
    font-size: 0.9375rem;
    font-weight: 600;
    color: #e2e8f0;
    margin: 0;
  }

  .item-description {
    font-size: 0.8125rem;
    color: #94a3b8;
    margin: 0;
    line-height: 1.4;
  }

  .item-patterns {
    display: flex;
    flex-wrap: wrap;
    gap: 0.375rem;
    margin-top: 0.25rem;
  }

  .pattern-tag {
    font-size: 0.6875rem;
    padding: 0.125rem 0.5rem;
    background: rgba(168, 85, 247, 0.2);
    color: #c4b5fd;
    border-radius: 4px;
    border: 1px solid rgba(168, 85, 247, 0.3);
  }

  .check-icon {
    color: #10b981;
    flex-shrink: 0;
  }

  /* Scrollbar styling */
  .dropdown-menu::-webkit-scrollbar {
    width: 6px;
  }

  .dropdown-menu::-webkit-scrollbar-track {
    background: rgba(15, 23, 42, 0.5);
  }

  .dropdown-menu::-webkit-scrollbar-thumb {
    background: rgba(71, 85, 105, 0.5);
    border-radius: 3px;
  }

  .dropdown-menu::-webkit-scrollbar-thumb:hover {
    background: rgba(71, 85, 105, 0.8);
  }
</style>
