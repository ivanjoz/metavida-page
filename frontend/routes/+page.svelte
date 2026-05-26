<script lang="ts">
  import {
    ArrowRight,
    CalendarDays,
    FlaskConical,
    HandHeart,
    HeartPulse,
    Landmark,
    Menu,
    MessageCircle,
    ShieldCheck,
    UsersRound,
    X
  } from 'lucide-svelte';

  let mobileMenuOpen = $state(false);

  const navItems = [
    { label: 'Inicio', href: '#inicio' },
    { label: 'Encuestas', href: '#encuestas' },
    { label: 'Noticias', href: '#noticias' },
    { label: 'Admin', href: '/admin' }
  ];

  const pillars = [
    {
      title: 'Impulsar',
      text: 'La investigación y la innovación de oncología',
      icon: FlaskConical
    },
    {
      title: 'Crear',
      text: 'Redes de soporte para pacientes y familias',
      icon: UsersRound
    },
    {
      title: 'Incentivar',
      text: 'La cooperación con organismos internacionales',
      icon: HandHeart
    },
    {
      title: 'Desarrollar',
      text: 'Programas de educación, concientización y nutrición oncológica',
      icon: MessageCircle
    }
  ];

  const surveys = [
    {
      label: 'Encuesta',
      title: 'Toxicidad Financiera y Calidad de Vida',
      text:
        'Participe en nuestra encuesta sobre toxicidad financiera y calidad de vida. Los costos asociados al cuidado de la salud pueden generar una carga económica y emocional significativa para pacientes y familias.',
      image: '/images/metavida/encuesta_toxi1.jpg',
      href: 'https://metavida.life/encuestatf/'
    },
    {
      label: 'Encuesta',
      title: 'Conociendo Nuestra Realidad',
      text:
        'Esta encuesta busca conocer la situación real que enfrentan los pacientes en el Perú, especialmente en el acceso a tratamientos, historias clínicas y apoyo recibido.',
      image: '/images/metavida/encuesta_paciente.jpg',
      href: 'https://metavida.life/encuestarp/'
    }
  ];

  const posts = [
    {
      title: 'Avances en Inmunoterapia: Nueva Esperanza Contra el Cáncer de Cabeza y Cuello',
      image: '/images/metavida/uero111-890x660.jpg',
      href: 'https://www.metavida.life/whats-the-reason-so-many-older-adults-arent-active/noticias/'
    },
    {
      title: 'Avances Globales en la Investigación del Cáncer: Nuevas Terapias y Desafíos Persistentes',
      image: '/images/metavida/dem111-890x660.jpg',
      href: 'https://www.metavida.life/the-most-important-ventilator-equipment-available/noticias/'
    },
    {
      title: 'El Poder de la Anticipación: Cómo el Tamizaje y la Prevención Están Cambiando la Lucha Contra el Cáncer',
      image: '/images/metavida/prevencion-890x660.jpg',
      href: 'https://www.metavida.life/blood-cancers-early-signs-symptoms-institute/noticias/'
    }
  ];
</script>

<svelte:head>
  <title>MetaVida</title>
  <meta
    name="description"
    content="MetaVida une esfuerzos por el bienestar, la salud y el acompañamiento de pacientes."
  />
</svelte:head>

<main class="min-h-screen bg-[#eef7ff] text-[#102a57]">
  <header class="fixed inset-x-0 top-0 z-40 border-b border-white/20 bg-[#0b2246]/82 backdrop-blur-md">
    <div class="mx-auto flex h-82 w-[min(1180px,calc(100%-32px))] items-center justify-between">
      <a href="#inicio" class="flex items-center gap-12" aria-label="MetaVida inicio">
        <img src="/images/metavida/logo.svg" alt="MetaVida" class="h-46 w-auto max-w-[178px]" />
      </a>

      <nav class="hidden items-center gap-30 text-sm font-semibold uppercase leading-[1.2] text-white/88 md:flex">
        {#each navItems as item}
          <a class="transition hover:text-[#70b8ff]" href={item.href}>{item.label}</a>
        {/each}
      </nav>

      <button
        class="grid h-42 w-42 place-items-center rounded-[6px] border border-white/30 text-white md:hidden"
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
          <a class="rounded-[6px] px-12 py-10 hover:bg-white/10" href={item.href} onclick={() => (mobileMenuOpen = false)}>
            {item.label}
          </a>
        {/each}
      </nav>
    {/if}
  </header>

  <section
    id="inicio"
    class="relative flex min-h-[760px] items-center overflow-hidden bg-[#0b2246] pt-82 text-white md:min-h-[820px]"
  >
    <img
      src="/images/metavida_fondo.avif"
      alt=""
      class="absolute inset-0 h-full w-full object-cover opacity-78"
    />
    <div class="absolute inset-0 bg-[linear-gradient(90deg,rgba(7,25,56,0.92),rgba(9,42,83,0.58)_52%,rgba(9,42,83,0.16))]"></div>
    <div class="absolute inset-0 bg-[linear-gradient(180deg,rgba(6,20,44,0.22),rgba(6,20,44,0)_40%,rgba(6,20,44,0.18))]"></div>

    <div class="relative mx-auto grid w-[min(1180px,calc(100%-32px))] gap-36 py-78">
      <p class="max-w-max rounded-full border border-white/28 bg-white/14 px-16 py-8 text-sm font-semibold uppercase leading-[1.2] text-[#d7ecff] shadow-[0_14px_34px_rgba(4,24,55,0.18)]">
        MetaVida
      </p>
      <h1 class="max-w-[920px] text-[44px] font-semibold leading-[1.03] text-white md:text-[76px] lg:text-[96px]">
        Uniendo esfuerzos por <span class="text-[#8cc9ff]">el bienestar y la salud.</span>
      </h1>
      <div class="flex max-w-[820px] flex-col gap-24 md:flex-row md:items-center md:justify-between">
        <p class="max-w-[590px] text-lg leading-[1.55] text-white/86 md:text-xl">
          Una comunidad enfocada en investigación, soporte a pacientes, educación y prevención oncológica en el Perú.
        </p>
        <a
          href="#encuestas"
          class="inline-flex h-54 shrink-0 items-center justify-center gap-10 rounded-full bg-[#2f8fe8] px-24 text-base font-bold text-white shadow-[0_18px_34px_rgba(12,83,160,0.28)] transition hover:bg-[#126fca]"
        >
          Participar
          <ArrowRight size={19} />
        </a>
      </div>
    </div>
  </section>

  <section class="relative z-10 -mt-54">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] overflow-hidden rounded-[28px] border border-white bg-white shadow-[0_24px_80px_rgba(18,64,111,0.14)] md:grid-cols-4">
      {#each pillars as pillar}
        {@const Icon = pillar.icon}
        <article class="flex min-h-152 gap-16 border-b border-[#dbe9f7] p-24 last:border-b-0 md:border-b-0 md:border-r md:last:border-r-0">
          <div class="grid h-48 w-48 shrink-0 place-items-center rounded-full bg-[#e8f4ff] text-[#176fbc]">
            <Icon size={24} />
          </div>
          <div>
            <h2 class="text-xl font-bold leading-[1.2] text-[#102a57]">{pillar.title}</h2>
            <p class="mt-8 text-sm leading-[1.45] text-[#71819b]">{pillar.text}</p>
          </div>
        </article>
      {/each}
    </div>
  </section>

  <section id="encuestas" class="py-92">
    <div class="mx-auto w-[min(1180px,calc(100%-32px))]">
      <div class="mb-44 flex flex-col gap-12 text-center md:items-center">
        <div>
          <p class="mx-auto max-w-max rounded-full bg-white px-22 py-9 text-sm font-bold uppercase leading-[1.2] text-[#2e7dd7] shadow-[0_10px_28px_rgba(21,72,128,0.08)]">Participacion ciudadana</p>
          <h2 class="mt-22 text-[36px] font-bold leading-tight text-[#102a57] md:text-[54px]">¡Su voz es esencial!</h2>
        </div>
        <p class="mx-auto max-w-[560px] text-base leading-[1.6] text-[#73819b]">
          Las encuestas ayudan a convertir experiencias reales en investigacion, evidencia y mejores redes de apoyo.
        </p>
      </div>

      <div class="grid gap-24 lg:grid-cols-2">
        {#each surveys as survey}
          <article class="grid overflow-hidden rounded-[30px] border border-white bg-white shadow-[0_22px_52px_rgba(21,72,128,0.10)] transition hover:-translate-y-2 hover:shadow-[0_28px_68px_rgba(21,72,128,0.16)] md:grid-cols-[260px_1fr]">
            <img src={survey.image} alt={survey.title} class="h-260 w-full object-cover md:h-full md:p-16 md:pr-0 md:[border-radius:30px]" />
            <div class="flex min-h-320 flex-col p-30">
              <p class="text-sm font-semibold leading-[1.2] text-[#7786a2]">{survey.label}</p>
              <h3 class="mt-10 text-2xl font-semibold leading-tight text-[#102a57]">{survey.title}</h3>
              <div class="my-22 h-1 bg-[#d8e5f4]"></div>
              <p class="flex-1 text-base leading-[1.6] text-[#73819b]">{survey.text}</p>
              <a
                href={survey.href}
                class="mt-22 inline-flex h-52 w-52 items-center justify-center rounded-full bg-[#f3f9ff] text-[#102a57] shadow-[0_12px_28px_rgba(21,72,128,0.12)] transition hover:bg-[#2f8fe8] hover:text-white"
                aria-label={`Ir a la encuesta: ${survey.title}`}
              >
                <ArrowRight size={20} />
              </a>
            </div>
          </article>
        {/each}
      </div>
    </div>
  </section>

  <section class="bg-[#0b2246] py-78 text-white">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-24 md:grid-cols-3">
      <div class="flex gap-16">
        <HeartPulse class="mt-4 text-[#8cc9ff]" size={30} />
        <div>
          <h2 class="text-2xl font-bold">Salud</h2>
          <p class="mt-10 leading-[1.6] text-white/72">Acompañamiento con enfoque humano para pacientes y familias.</p>
        </div>
      </div>
      <div class="flex gap-16">
        <ShieldCheck class="mt-4 text-[#6ed6ff]" size={30} />
        <div>
          <h2 class="text-2xl font-bold">Prevencion</h2>
          <p class="mt-10 leading-[1.6] text-white/72">Educación, concientización y acciones tempranas contra el cáncer.</p>
        </div>
      </div>
      <div class="flex gap-16">
        <Landmark class="mt-4 text-[#b7dcff]" size={30} />
        <div>
          <h2 class="text-2xl font-bold">Investigacion</h2>
          <p class="mt-10 leading-[1.6] text-white/72">Información social y científica para impulsar mejores decisiones.</p>
        </div>
      </div>
    </div>
  </section>

  <section id="noticias" class="py-86">
    <div class="mx-auto w-[min(1180px,calc(100%-32px))]">
      <div class="mb-34">
        <p class="max-w-max rounded-full bg-white px-22 py-9 text-sm font-bold uppercase leading-[1.2] text-[#2e7dd7] shadow-[0_10px_28px_rgba(21,72,128,0.08)]">Nuestro Blog</p>
        <h2 class="mt-18 text-[36px] font-bold leading-tight text-[#102a57] md:text-[54px]">Noticias</h2>
      </div>

      <div class="grid gap-22 md:grid-cols-3">
        {#each posts as post}
          <article class="overflow-hidden rounded-[26px] border border-white bg-white shadow-[0_18px_44px_rgba(21,72,128,0.10)]">
            <a href={post.href} class="block">
              <img src={post.image} alt="" class="aspect-[1.35] w-full object-cover" />
            </a>
            <div class="p-22">
              <div class="mb-14 flex flex-wrap gap-14 text-sm leading-[1.35] text-[#73819b]">
                <span class="inline-flex items-center gap-6"><CalendarDays size={16} /> mayo 31, 2025</span>
                <span>Metavida</span>
              </div>
              <h3 class="text-xl font-bold leading-[1.28]">
                <a href={post.href} class="transition hover:text-[#2f8fe8]">{post.title}</a>
              </h3>
            </div>
          </article>
        {/each}
      </div>
    </div>
  </section>

  <footer class="bg-[#0b2246] py-28 text-white">
    <div class="mx-auto flex w-[min(1180px,calc(100%-32px))] flex-col gap-16 text-sm text-white/72 md:flex-row md:items-center md:justify-between">
      <p>Copyright © 2026 MetaVida. All Rights Reserved.</p>
      <div class="flex gap-18">
        <a class="hover:text-white" href="/admin">Admin</a>
        <a class="hover:text-white" href="https://www.metavida.life/mision/">Mision</a>
      </div>
    </div>
  </footer>
</main>
