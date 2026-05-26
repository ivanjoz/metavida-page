<script lang="ts">
  import { onMount } from 'svelte';
  import { Home, Settings, Users } from 'lucide-svelte';
  import Button from '$components/buttons/Button.svelte';
  import { checkIsLogin, Env } from '$core/env';

  let { children } = $props();
  let ready = $state(false);

  const nav = [
    { href: '/admin', label: 'Panel', icon: Home },
    { href: '/admin/leads', label: 'Leads', icon: Users },
    { href: '/admin/settings', label: 'Ajustes', icon: Settings }
  ];

  onMount(() => {
    if (checkIsLogin() !== 2) {
      window.location.href = '/login';
      return;
    }
    ready = true;
  });
</script>

{#if ready}
<div class="admin-layout">
  <aside class="admin-sidebar">
    <a class="brand" href="/">Metavida</a>
    <nav>
      {#each nav as item}
        {@const Icon = item.icon}
        <a href={item.href}>
          <Icon size={18} />
          {item.label}
        </a>
      {/each}
    </nav>

    <div class="sidebar-footer">
      <Button
        name="Salir"
        color="red"
        css="s1"
        onClick={() => Env.clearAccesos?.()}
      />
    </div>
  </aside>

  <main class="admin-content">
    {@render children?.()}
  </main>
</div>
{/if}

<style>
  .brand {
    display: block;
    margin-bottom: 28px;
    color: #15493e;
    font-size: 22px;
    font-weight: 800;
  }

  nav {
    display: grid;
    gap: 8px;
  }

  nav a {
    display: flex;
    align-items: center;
    gap: 10px;
    min-height: 42px;
    padding: 0 12px;
    border-radius: 8px;
    color: #24312e;
  }

  nav a:hover {
    background: #edf3e9;
  }

  .sidebar-footer {
    margin-top: 32px;
  }
</style>
