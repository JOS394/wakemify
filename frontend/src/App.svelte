<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { Toggle, ActivateForMinutes, Deactivate, GetStatus, GetRemainingTime, QuitApp } from '../wailsjs/go/main/App.js'
  import { EventsOn } from '../wailsjs/runtime/runtime.js'
  import { Moon, Zap, Clock, Power, LogOut } from 'lucide-svelte'

  let active = false
  let errorMsg = ''
  let remainingSeconds = 0
  let timerInterval: number | undefined

  onMount(async () => {
    const status = await GetStatus()
    active = status

    if (status) {
      startRemainingTimer()
    }

    EventsOn('timer-expired', () => {
      active = false
      clearInterval(timerInterval)
      remainingSeconds = 0
    })
  })

  onDestroy(() => {
    clearInterval(timerInterval)
  })

  function startRemainingTimer() {
    clearInterval(timerInterval)
    timerInterval = window.setInterval(async () => {
      const secs = await GetRemainingTime()
      remainingSeconds = secs
      if (secs <= 0) {
        active = false
        clearInterval(timerInterval)
      }
    }, 1000)
  }

  async function handleActivateIndefinite() {
    errorMsg = ''
    try {
      const result = await Toggle()
      active = result
      if (!result) remainingSeconds = 0
    } catch (e) {
      errorMsg = String(e)
    }
  }

  async function handleTimer(minutes: number) {
    errorMsg = ''
    try {
      const result = await ActivateForMinutes(minutes)
      active = result
      if (result) startRemainingTimer()
    } catch (e) {
      errorMsg = String(e)
    }
  }

  async function handleDeactivate() {
    errorMsg = ''
    try {
      await Deactivate()
      active = false
      remainingSeconds = 0
      clearInterval(timerInterval)
    } catch (e) {
      errorMsg = String(e)
    }
  }

  function handleQuit() {
    QuitApp()
  }

  function formatTime(seconds: number): string {
    if (seconds <= 0) return ''
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    return `${m}:${s.toString().padStart(2, '0')}`
  }

  $: titleText = active ? 'Wakemify ✓' : 'Wakemify'
  $: timerText = formatTime(remainingSeconds)
</script>

<div class="menu">
  <div class="menu-title">
    <span>{titleText}</span>
    {#if timerText}
      <span class="timer-text">{timerText}</span>
    {/if}
  </div>

  <div class="menu-separator"></div>

  <button class="menu-item" disabled={active} on:click={handleActivateIndefinite}>
    <Zap size={14} stroke={active ? '#555' : '#d4a373'} stroke-width={2.5} />
    <span>Activar (Indefinido)</span>
  </button>

  <button class="menu-item" disabled={active} on:click={() => handleTimer(30)}>
    <Clock size={14} stroke={active ? '#555' : '#ccc'} stroke-width={2} />
    <span>30 minutos</span>
  </button>

  <button class="menu-item" disabled={active} on:click={() => handleTimer(60)}>
    <Clock size={14} stroke={active ? '#555' : '#ccc'} stroke-width={2} />
    <span>1 hora</span>
  </button>

  <button class="menu-item" disabled={active} on:click={() => handleTimer(120)}>
    <Clock size={14} stroke={active ? '#555' : '#ccc'} stroke-width={2} />
    <span>2 horas</span>
  </button>

  <button class="menu-item" disabled={active} on:click={() => handleTimer(300)}>
    <Clock size={14} stroke={active ? '#555' : '#ccc'} stroke-width={2} />
    <span>5 horas</span>
  </button>

  <div class="menu-separator"></div>

  <button class="menu-item deactivate" disabled={!active} on:click={handleDeactivate}>
    <Power size={14} stroke={active ? '#c0392b' : '#555'} stroke-width={2} />
    <span>Desactivar</span>
  </button>

  <div class="menu-separator"></div>

  <button class="menu-item quit" on:click={handleQuit}>
    <LogOut size={14} stroke="#888" stroke-width={2} />
    <span>Salir</span>
  </button>

  {#if errorMsg}
    <div class="menu-separator"></div>
    <div class="menu-item error-text">{errorMsg}</div>
  {/if}
</div>

<style>
  .menu {
    width: 100%;
    height: 100vh;
    background: #1e1e1e;
    padding: 4px 0;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
  }

  .menu-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 14px;
    font-size: 13px;
    font-weight: 600;
    color: #aaa;
    -webkit-app-region: drag;
  }

  .timer-text {
    font-size: 11px;
    color: #d4a373;
    font-variant-numeric: tabular-nums;
  }

  .menu-separator {
    height: 1px;
    background: #333;
    margin: 4px 10px;
  }

  .menu-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 7px 14px;
    background: transparent;
    border: none;
    color: #ccc;
    font-size: 13px;
    cursor: pointer;
    text-align: left;
    transition: background 0.1s;
  }

  .menu-item:hover:not(:disabled) {
    background: #333;
  }

  .menu-item:disabled {
    color: #555;
    cursor: default;
  }

  .menu-item.deactivate:not(:disabled) {
    color: #c0392b;
  }

  .menu-item.quit:hover {
    background: rgba(255, 255, 255, 0.05);
  }

  .error-text {
    color: #e74c3c;
    font-size: 11px;
    cursor: default;
  }
</style>
