<script lang="ts">
  import { Home, Settings, Users } from 'lucide-svelte';

  let { children } = $props();

  const nav = [
    { href: '/admin', label: 'Panel', icon: Home },
    { href: '/admin/leads', label: 'Leads', icon: Users },
    { href: '/admin/settings', label: 'Ajustes', icon: Settings }
  ];
</script>

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
  </aside>

  <main class="admin-content">
    {@render children?.()}
  </main>
</div>

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
</style>
