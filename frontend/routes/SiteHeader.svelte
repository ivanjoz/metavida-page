<script lang="ts">
  import { onMount } from 'svelte';
  import { base } from '$app/paths';
  import { Menu, X } from 'lucide-svelte';
  import { navItems } from '$routes/page.content';

  let { forceCompact = false } = $props<{ forceCompact?: boolean }>();

  // Header owns mobile state so every public page gets the same navigation behavior.
  let mobileMenuOpen = $state(false);
  let compactHeader = $state(false);
  let headerIsCompact = $derived(forceCompact || compactHeader);

  function syncHeaderHeight() {
    const shouldUseCompactHeader = window.scrollY > 24;

    if (compactHeader === shouldUseCompactHeader) return;

    // Log only state transitions so scroll debugging stays useful without flooding the console.
    console.debug('Public header compact state changed', {
      compactHeader: shouldUseCompactHeader,
      scrollY: window.scrollY
    });

    compactHeader = shouldUseCompactHeader;
  }

  onMount(() => {
    if (forceCompact) {
      // Client pages keep the public header visible without using the tall landing-page height.
      console.debug('Public header forced into compact mode for embedded app area');
      return;
    }

    syncHeaderHeight();
    window.addEventListener('scroll', syncHeaderHeight, { passive: true });

    return () => window.removeEventListener('scroll', syncHeaderHeight);
  });
</script>

<header
  class="fixed inset-x-0 top-0 z-40 border-b border-white/20 bg-[#0b2246]/82 backdrop-blur-md transition-[background-color,box-shadow] duration-300 ease-out"
  class:shadow-[0_12px_34px_rgba(4,19,44,0.2)]={headerIsCompact}
>
  <div
    class="mx-auto flex w-[min(1180px,calc(100%-32px))] items-center justify-between transition-[height] duration-300 ease-out"
    class:h-82={!headerIsCompact}
    class:h-62={headerIsCompact}
  >
    <a href={`${base}/`} class="flex items-center gap-12" aria-label="MetaVida inicio">
      <img
        src={`${base}/images/metavida/logo.svg`}
        alt="MetaVida"
        class="w-auto max-w-[178px] transition-[height] duration-300 ease-out"
        class:h-46={!headerIsCompact}
        class:h-34={headerIsCompact}
      />
    </a>

    <nav class="hidden items-center gap-30 text-sm font-semibold uppercase leading-[1.2] text-white/88 md:flex">
      {#each navItems as item}
        <a
          class="group relative py-8 transition duration-200 hover:-translate-y-1 hover:text-[#70b8ff] focus-visible:text-[#70b8ff] focus-visible:outline-none"
          href={item.href}
        >
          {item.label}
          <span class="absolute inset-x-0 bottom-1 h-2 origin-left scale-x-0 rounded-full bg-[#70b8ff] transition-transform duration-200 group-hover:scale-x-100 group-focus-visible:scale-x-100"></span>
        </a>
      {/each}
    </nav>

    <button
      class="grid w-42 place-items-center rounded-[6px] border border-white/30 text-white transition-[height] duration-300 ease-out md:hidden"
      class:h-42={!headerIsCompact}
      class:h-36={headerIsCompact}
      type="button"
      aria-label="Abrir menu"
      onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
    >
      {#if mobileMenuOpen}
        <X size={21} />
      {:else}
        <Menu size={21} />
      {/if}
    </button>
  </div>

  {#if mobileMenuOpen}
    <nav class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-4 pb-18 text-sm font-semibold uppercase leading-[1.2] text-white md:hidden">
      {#each navItems as item}
        <a
          class="rounded-[6px] px-12 py-10 transition duration-200 hover:bg-white/12 hover:text-[#70b8ff] focus-visible:bg-white/12 focus-visible:text-[#70b8ff] focus-visible:outline-none"
          href={item.href}
          onclick={() => (mobileMenuOpen = false)}
        >
          {item.label}
        </a>
      {/each}
    </nav>
  {/if}
</header>
