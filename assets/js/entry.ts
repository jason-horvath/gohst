import Alpine from "alpinejs";

// const vitePort = import.meta.env.VITE_PORT || 5173;

// const ws = new WebSocket(`ws://localhost:${vitePort}`);
// ws.onmessage = (event) => {
//   if (event.data === "full-reload") {
//     console.log("🔄 Vite restarted! Refreshing browser...");
//     location.reload(); // Force refresh
//   }
// };

// // Handle WebSocket errors silently
// ws.onerror = (error) => {
//   console.warn("⚠️ Vite WebSocket error (probably restarting):", error);
// };

// // Handle WebSocket closing gracefully
// ws.onclose = () => {
//   console.warn("⚠️ Vite WebSocket closed, likely because Vite is restarting.");
// };

declare global {
  interface Window {
    Alpine: typeof Alpine;
  }
}

window.Alpine = Alpine;
Alpine.start();

console.log("Alpine.js loaded");
