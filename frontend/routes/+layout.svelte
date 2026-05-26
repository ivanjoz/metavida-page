<script lang="ts">
  import './tailwind.css';
  import './app.css';

  let { children } = $props();
</script>

<svelte:head>
  <title>Metavida</title>
  <meta
    name="description"
    content="CRM y pagina publica del proyecto Metavida."
  />
</svelte:head>

{@render children?.()}
