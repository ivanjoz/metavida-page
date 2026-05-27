<script lang="ts">
  import { page } from '$app/state';
  import { ClipboardList, Settings } from 'lucide-svelte';

  const clientMenuItems = [
    { label: 'Mis Consultas', href: '/client', icon: ClipboardList },
    { label: 'Configuración', href: '/client/configuracion', icon: Settings }
  ];

  function isCurrentClientSection(href: string) {
    // Match the exact client landing route separately so nested pages do not highlight two entries.
    return href === '/client' ? page.url.pathname === href : page.url.pathname.startsWith(href);
  }

  function logClientNavigation(label: string, href: string) {
    // Keep client-area navigation traceable while the section is still pre-alpha.
    console.debug('Client side menu navigation requested', { label, href });
  }
</script>

<aside class="client-menu">
  <div class="client-menu__title text-lg font-bold leading-tight">
    <span>Mi Cuenta</span>
  </div>

  <nav class="client-menu__nav" aria-label="Menu de cliente">
    {#each clientMenuItems as item}
      {@const Icon = item.icon}
      {@const active = isCurrentClientSection(item.href)}
      <a
        class="text-sm font-semibold leading-[1.2]"
        class:active
        href={item.href}
        aria-current={active ? 'page' : undefined}
        onclick={() => logClientNavigation(item.label, item.href)}
      >
        <Icon size={18} strokeWidth={1.9} />
        <span>{item.label}</span>
      </a>
    {/each}
  </nav>
</aside>

<style>
  .client-menu {
    min-height: calc(100vh - 62px);
    border-right: 1px solid #dbe7f2;
    background: #ffffff;
    padding: 28px 18px;
    box-shadow: 12px 0 32px rgba(16, 42, 87, 0.05);
  }

  .client-menu__title {
    margin-bottom: 18px;
    padding: 0 10px 16px;
    border-bottom: 1px solid #edf3f9;
    color: #075f80;
  }

  .client-menu__nav {
    display: grid;
    gap: 8px;
  }

  a {
    display: flex;
    align-items: center;
    gap: 10px;
    min-height: 44px;
    padding: 0 12px;
    border-radius: 6px;
    color: #2c4056;
    transition: background 0.18s ease, color 0.18s ease, transform 0.18s ease;
  }

  a:hover,
  a:focus-visible {
    background: #eef7ff;
    color: #075f80;
    outline: none;
  }

  a.active {
    background: #dff3ff;
    color: #075f80;
    box-shadow: inset 3px 0 0 #18c682;
  }

  @media (max-width: 748px) {
    .client-menu {
      min-height: auto;
      border-right: 0;
      border-bottom: 1px solid #dbe7f2;
      padding: 14px;
      box-shadow: 0 12px 24px rgba(16, 42, 87, 0.04);
    }

    .client-menu__title {
      display: none;
    }

    .client-menu__nav {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }
</style>
