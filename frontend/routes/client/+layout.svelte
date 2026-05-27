<script lang="ts">
  import SiteHeader from '../SiteHeader.svelte';
  import ClientSideMenu from './ClientSideMenu.svelte';

  let { children } = $props();
</script>

<main class="client-page">
  <SiteHeader forceCompact />

  <div class="client-shell">
    <ClientSideMenu />

    <section class="client-content">
      {@render children?.()}
    </section>
  </div>
</main>

<style>
  .client-page {
    min-height: 100vh;
    padding-top: 62px;
    background: #eef7ff;
    color: #102a57;
  }

  .client-shell {
    display: grid;
    grid-template-columns: 250px minmax(0, 1fr);
    min-height: calc(100vh - 62px);
  }

  .client-content {
    width: min(980px, calc(100% - 48px));
    margin: 0 auto;
    padding: 36px 0 56px;
  }

  @media (max-width: 748px) {
    .client-shell {
      grid-template-columns: 1fr;
    }

    .client-content {
      width: min(100% - 28px, 680px);
      padding: 22px 0 40px;
    }
  }
</style>
