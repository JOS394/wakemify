# Wakemify

Keep your Mac awake — una app de menubar como Amphetamine o Caffeine, construida con Wails.

## Features

- Prevención de sueño del sistema vía IOKit (`PreventUserIdleSystemSleep`)
- Icono en la barra de menú (menubar) con indicador visual Sol/Luna
- Ventana popover al hacer click en el icono
- Activar indefinidamente o con temporizador: 30 min, 1 hora, 2 horas, 5 horas
- Timer countdown en tiempo real
- Interfaz oscura nativa con iconos Lucide
- About Wakemify con info del proyecto, donación y contacto
- Iniciar con el sistema (checkbox en About)
- Cierre automático al perder foco

## Stack

| Capa     | Tecnología                          |
| -------- | ----------------------------------- |
| Backend  | Go, Wails v2, CGO (IOKit)           |
| Frontend | Svelte 3 + TypeScript               |
| Menubar  | NSStatusItem nativo via CGO/ObjC    |
| Iconos   | Lucide Svelte (Moon / Sun)          |

## Requisitos

- macOS (Apple Silicon o Intel)
- Go 1.23+
- Node.js 18+
- Wails CLI v2

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Desarrollo

```bash
git clone <repo>
cd wakemify
wails dev
```

Esto inicia el servidor de desarrollo de Vite con hot-reload y compila la app en modo debug.

## Build producción

```bash
wails build
```

El binario empaquetado se genera en `build/bin/`.

## Estructura

```
wakemify/
├── main.go              # Entry point de Wails
├── app.go               # Lógica expuesta al frontend
├── power_darwin.go      # CGO + IOKit (prevención de sueño)
├── systray.go           # NSStatusItem nativo (ObjC vía CGO)
├── systray_cb.go        # Callbacks //export entre ObjC y Go
├── icons/               # Iconos de la menubar (Moon / Sun)
├── frontend/
│   ├── src/
│   │   ├── App.svelte   # UI tipo menú con Lucide icons
│   │   ├── main.ts
│   │   └── style.css
│   ├── index.html
│   └── package.json
└── build/
    └── darwin/
        └── Info.plist   # LSUIElement para app sin dock
```

## Uso

1. Al iniciar, la app aparece solo en la barra de menú (sin dock icon)
2. Click en el icono → ventana popover con las opciones
3. Elegir **Activar (Indefinido)** o un temporizador
4. El icono cambia a Sol ☀️ cuando está activo
5. Click en **Desactivar** para detener
6. **About Wakemify** → info, donación (PayPal), sugerencias (email)
7. **Salir** cierra la app y libera el assertion

## Contacto

- Correo: joseespana94@gmail.com
- PayPal: https://paypal.me/JOS32994

## Licencia

GNU GPLv3
