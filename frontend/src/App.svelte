<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import {
    Toggle,
    ActivateForMinutes,
    Deactivate,
    GetStatus,
    GetRemainingTime,
    QuitApp,
    OpenURL,
    RestoreMenuPosition,
    SetLaunchAtStartup,
    IsLaunchAtStartup,
  } from "../wailsjs/go/main/App.js";
  import {
    EventsOn,
    WindowSetSize,
    WindowCenter,
  } from "../wailsjs/runtime/runtime.js";
  import {
    Moon,
    Sun,
    Clock,
    Power,
    LogOut,
    Info,
    Coffee,
    Mail,
    ArrowLeft,
    Square,
    CheckSquare,
  } from "lucide-svelte";

  let active = false;
  let errorMsg = "";
  let remainingSeconds = 0;
  let timerInterval: number | undefined;
  let view: "menu" | "about" = "menu";
  let launchAtStartup = false;

  onMount(async () => {
    const status = await GetStatus();
    active = status;

    if (status) {
      startRemainingTimer();
    }

    const startup = await IsLaunchAtStartup();
    launchAtStartup = startup;

    EventsOn("timer-expired", () => {
      active = false;
      clearInterval(timerInterval);
      remainingSeconds = 0;
    });
  });

  onDestroy(() => {
    clearInterval(timerInterval);
  });

  function startRemainingTimer() {
    clearInterval(timerInterval);
    timerInterval = window.setInterval(async () => {
      const secs = await GetRemainingTime();
      remainingSeconds = secs;
      if (secs <= 0) {
        active = false;
        clearInterval(timerInterval);
      }
    }, 1000);
  }

  async function handleActivateIndefinite() {
    errorMsg = "";
    try {
      const result = await Toggle();
      active = result;
      if (!result) remainingSeconds = 0;
    } catch (e) {
      errorMsg = String(e);
    }
  }

  async function handleTimer(minutes: number) {
    errorMsg = "";
    try {
      const result = await ActivateForMinutes(minutes);
      active = result;
      if (result) startRemainingTimer();
    } catch (e) {
      errorMsg = String(e);
    }
  }

  async function handleDeactivate() {
    errorMsg = "";
    try {
      await Deactivate();
      active = false;
      remainingSeconds = 0;
      clearInterval(timerInterval);
    } catch (e) {
      errorMsg = String(e);
    }
  }

  function handleQuit() {
    QuitApp();
  }

  function handleCoffee() {
    OpenURL("https://paypal.me/JOS32994");
  }

  async function toggleLaunchAtStartup() {
    const result = await SetLaunchAtStartup(!launchAtStartup);
    launchAtStartup = result;
  }

  function handleEmail() {
    OpenURL("mailto:joseespana94@gmail.com");
  }

  function formatTime(seconds: number): string {
    if (seconds <= 0) return "";
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}:${s.toString().padStart(2, "0")}`;
  }

  function openAbout() {
    view = "about";
    WindowSetSize(460, 460);
    WindowCenter();
  }

  function closeAbout() {
    view = "menu";
    WindowSetSize(240, 290);
    RestoreMenuPosition();
  }

  $: titleText = active ? "Wakemify ✓" : "Wakemify";
  $: timerText = formatTime(remainingSeconds);
</script>

{#if view === "about"}
  <div class="about-view">
    <div class="about-header">
      <button class="back-btn" on:click={closeAbout}>
        <ArrowLeft size={14} stroke-width={2} />
      </button>
      <span class="about-title">About Wakemify</span>
    </div>

    <div class="about-body">
      <div class="about-desc">
        <strong>Wakemify v1.0.0</strong><br /><br />
        Alternativa moderna, minimalista y ligera para mantener tu Mac despierta.<br
        />
        Integrada mediante IOKit de macOS desde tu barra de menú.<br /><br />
        ⚡️ Prevención de sueño nativa<br />
        ⏱️ Temporizadores: 30m, 1h, 2h y 5h<br />
        🌙 Interfaz oscura nativa<br /><br />
        Wails + Svelte • GNU GPLv3
      </div>

      <button class="about-btn" on:click={handleCoffee}>
        <Coffee size={16} stroke="#d4a373" stroke-width={2} />
        <span>Regálame un café</span>
      </button>

      <button class="about-btn" on:click={handleEmail}>
        <Mail size={16} stroke="#888" stroke-width={2} />
        <span>Envíame un correo</span>
      </button>
      <span class="about-email">joseespana94@gmail.com</span>

      <div class="startup-row">
        <div class="startup-toggle" on:click={toggleLaunchAtStartup}>
          {#if launchAtStartup}
            <CheckSquare size={16} stroke="#d4a373" stroke-width={2} />
          {:else}
            <Square size={16} stroke="#666" stroke-width={2} />
          {/if}
          <span>Iniciar con el sistema</span>
        </div>
      </div>
    </div>
  </div>
{:else}
  <div class="menu">
    <div class="menu-title">
      <span>{titleText}</span>
      {#if timerText}
        <span class="timer-text">{timerText}</span>
      {/if}
    </div>

    <div class="menu-separator"></div>

    <button
      class="menu-item"
      disabled={active}
      on:click={handleActivateIndefinite}
    >
      <Sun size={14} stroke={active ? "#555" : "#d4a373"} stroke-width={2.5} />
      <span>Activar (Indefinido)</span>
    </button>

    <button
      class="menu-item"
      disabled={active}
      on:click={() => handleTimer(30)}
    >
      <Clock size={14} stroke={active ? "#555" : "#ccc"} stroke-width={2} />
      <span>30 minutos</span>
    </button>

    <button
      class="menu-item"
      disabled={active}
      on:click={() => handleTimer(60)}
    >
      <Clock size={14} stroke={active ? "#555" : "#ccc"} stroke-width={2} />
      <span>1 hora</span>
    </button>

    <button
      class="menu-item"
      disabled={active}
      on:click={() => handleTimer(120)}
    >
      <Clock size={14} stroke={active ? "#555" : "#ccc"} stroke-width={2} />
      <span>2 horas</span>
    </button>

    <button
      class="menu-item"
      disabled={active}
      on:click={() => handleTimer(300)}
    >
      <Clock size={14} stroke={active ? "#555" : "#ccc"} stroke-width={2} />
      <span>5 horas</span>
    </button>

    <div class="menu-separator"></div>

    <button
      class="menu-item deactivate"
      disabled={!active}
      on:click={handleDeactivate}
    >
      <Power size={14} stroke={active ? "#c0392b" : "#555"} stroke-width={2} />
      <span>Desactivar</span>
    </button>

    <div class="menu-separator"></div>

    <button class="menu-item" on:click={openAbout}>
      <Info size={14} stroke="#888" stroke-width={2} />
      <span>About Wakemify</span>
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
{/if}

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
    padding: 5px 14px;
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

  .about-view {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: #1e1e1e;
    padding: 4px 0;
    box-sizing: border-box;
  }

  .about-header {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 6px 14px;
    -webkit-app-region: drag;
  }

  .back-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
    color: #fff;
    -webkit-app-region: no-drag;
  }

  .back-btn:hover {
    background: #333;
  }

  .about-title {
    font-size: 13px;
    font-weight: 600;
    color: #aaa;
  }

  .about-body {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    padding: 14px 14px;
    flex: 1;
  }

  .about-desc {
    color: #888;
    font-size: 12px;
    text-align: center;
    line-height: 1.6;
    margin: 0;
  }

  .about-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    width: 180px;
    padding: 5px 10px;
    background: #2a2a2a;
    border: 1px solid #333;
    border-radius: 6px;
    color: #ccc;
    font-size: 12px;
    cursor: pointer;
    transition: background 0.1s;
  }

  .about-btn:hover {
    background: #333;
    border-color: #555;
  }

  .about-email {
    font-size: 11px;
    color: #666;
    margin-top: -6px;
  }

  .startup-row {
    display: flex;
    justify-content: flex-end;
    width: 100%;
    margin-top: auto;
  }

  .startup-toggle {
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    color: #999;
    font-size: 12px;
    padding: 4px 8px;
    border-radius: 4px;
    user-select: none;
  }

  .startup-toggle:hover {
    background: #2a2a2a;
    color: #ccc;
  }
</style>
