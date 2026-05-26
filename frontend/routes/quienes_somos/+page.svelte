<script lang="ts">
  import {
    Building2,
    HandHeart,
    HeartPulse,
    Lightbulb,
    Microscope,
    Scale,
    Sprout,
    UsersRound
  } from 'lucide-svelte';
  import SiteHeader from '../SiteHeader.svelte';

  // Keep the old-page mission goals as data so the markup stays focused on layout.
  const missionGoals = [
    {
      text: 'Promover la asistencia social a pacientes con enfermedades críticas, en especial a pacientes oncológicos.',
      icon: HeartPulse
    },
    {
      text: 'Gestionar con organismos públicos y privados, nacionales o extranjeros, iniciativas o proyectos que contribuyan al bienestar general de las personas.',
      icon: Building2
    },
    {
      text: 'Defender los derechos de los pacientes e impulsar proyectos de asistencia social, acceso a ensayos clínicos, terapias avanzadas, políticas de salud pública, reformas legislativas y programas de atención médica.',
      icon: Scale
    },
    {
      text: 'Promover una cultura de conciencia social, capacitar y facilitar herramientas informativas y de atención para una mejor calidad de vida.',
      icon: Lightbulb
    },
    {
      text: 'Impulsar la investigación y la innovación en enfermedades críticas, especialmente en oncología.',
      icon: Microscope
    },
    {
      text: 'Realizar actividades sociales, culturales, académicas, deportivas, religiosas y empresariales para recaudar fondos para la asociación.',
      icon: UsersRound
    },
    {
      text: 'Aceptar donaciones para financiar los proyectos aprobados por la asociación.',
      icon: HandHeart
    },
    {
      text: 'Adquirir bienes muebles e inmuebles en beneficio de la asociación y disponer de ellos acorde con los objetivos institucionales.',
      icon: Sprout
    }
  ];
</script>

<svelte:head>
  <title>Quiénes Somos | MetaVida</title>
  <meta
    name="description"
    content="Conoce la misión y visión de MetaVida, una asociación que coordina esfuerzos solidarios por la salud, el bienestar y la calidad de vida."
  />
</svelte:head>

<main class="min-h-screen bg-[#f4f9ff] text-[#102a57]">
  <SiteHeader />

  <section class="relative isolate overflow-hidden bg-[#0b2246] pt-82 text-white">
    <img src="/images/metavida_fondo.avif" alt="" class="absolute inset-0 -z-20 h-full w-full object-cover opacity-55" />
    <div class="absolute inset-0 -z-10 bg-[linear-gradient(110deg,rgba(7,25,56,0.96),rgba(7,59,91,0.78)_58%,rgba(16,125,143,0.44))]"></div>

    <div class="mx-auto grid min-h-[520px] w-[min(1180px,calc(100%-32px))] content-end py-78">
      <div class="max-w-[900px]">
        <p class="max-w-max rounded-full bg-[#18c682] px-18 py-8 text-sm font-bold uppercase leading-[1.2] text-[#073953]">
          Quiénes somos
        </p>
        <h1 class="mt-24 text-[42px] font-bold leading-[1.02] text-white md:text-[70px]">
          Una red solidaria para cuidar la vida con ciencia, dignidad y humanidad.
        </h1>
        <p class="mt-22 max-w-[690px] text-lg leading-[1.6] text-white/84 md:text-xl">
          MetaVida unifica esfuerzos para acompañar a personas que enfrentan enfermedades críticas, sus familias y comunidades.
        </p>
      </div>
    </div>
  </section>

  <section class="bg-white py-86">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-42 lg:grid-cols-[0.92fr_1.08fr] lg:items-center">
      <div class="relative">
        <img
          src="/images/metavida/mision.png"
          alt="Representación visual de la misión de MetaVida"
          class="aspect-square w-full rounded-[8px] object-cover shadow-[0_26px_70px_rgba(21,72,128,0.14)]"
        />
      </div>

      <article class="grid gap-22">
        <div>
          <p class="text-base font-bold uppercase leading-[1.2] text-[#18a878]">Nuestra misión</p>
          <h2 class="mt-8 text-[36px] font-bold uppercase leading-[1.02] text-[#0784a7] md:text-[58px]">
            Unificar y coordinar esfuerzos solidarios
          </h2>
          <div class="mt-14 h-5 w-52 bg-[#18c682]"></div>
        </div>

        <p class="text-lg leading-[1.72] text-[#27384a]">
          La misión de la Asociación Metavida es unificar y coordinar esfuerzos solidarios en beneficio de todos sus integrantes, priorizando el apoyo integral en el ámbito de la salud, así como en todas aquellas áreas que influyen directamente en el bienestar general y la calidad de vida del ser humano.
        </p>
        <p class="text-lg leading-[1.72] text-[#27384a]">
          Esta misión se basa en el principio de que cada persona, especialmente quienes enfrentan enfermedades críticas, merece acceso equitativo a atención médica de calidad, contención emocional, oportunidades de desarrollo personal y respaldo social. A través de un enfoque multidisciplinario e inclusivo, Metavida busca ser un puente entre quienes necesitan ayuda y quienes desean brindarla, generando un impacto positivo y duradero en la vida de las personas, sus familias y comunidades.
        </p>
      </article>
    </div>
  </section>

  <section class="bg-[#eef7ff] py-88">
    <div class="mx-auto w-[min(1180px,calc(100%-32px))]">
      <div class="mb-38 max-w-[760px]">
        <p class="text-base font-bold uppercase leading-[1.2] text-[#18a878]">Cómo lo hacemos</p>
        <h2 class="mt-8 text-[34px] font-bold leading-tight text-[#102a57] md:text-[52px]">
          Líneas de acción que convierten la solidaridad en apoyo concreto
        </h2>
      </div>

      <div class="grid gap-16 md:grid-cols-2">
        {#each missionGoals as goal}
          {@const Icon = goal.icon}
          <article class="grid min-h-136 grid-cols-[58px_1fr] gap-18 rounded-[8px] border border-[#dbeaf8] bg-white p-22 shadow-[0_14px_34px_rgba(21,72,128,0.08)]">
            <div class="grid h-50 w-50 place-items-center rounded-[8px] bg-[#e8fff5] text-[#128b65]">
              <Icon size={25} strokeWidth={1.8} />
            </div>
            <p class="text-base leading-[1.58] text-[#27384a]">{goal.text}</p>
          </article>
        {/each}
      </div>
    </div>
  </section>

  <section class="bg-white py-92">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-44 lg:grid-cols-[1.02fr_0.98fr] lg:items-center">
      <article class="order-2 grid gap-22 lg:order-1">
        <div>
          <p class="text-base font-bold uppercase leading-[1.2] text-[#18a878]">Nuestra visión</p>
          <h2 class="mt-8 text-[36px] font-bold uppercase leading-[1.02] text-[#0784a7] md:text-[58px]">
            Representar, unir y acompañar
          </h2>
          <div class="mt-14 h-5 w-52 bg-[#18c682]"></div>
        </div>
        <p class="text-xl leading-[1.72] text-[#27384a]">
          Ser una asociación representativa ante las autoridades municipales, regionales y otras instituciones del estado, que unifica esfuerzos para el apoyo integral en el área de la salud y el bienestar general de sus integrantes y la comunidad.
        </p>
      </article>

      <img
        src="/images/metavida/vision.png"
        alt="Representación visual de la visión de MetaVida"
        class="order-1 aspect-square w-full rounded-[8px] object-cover shadow-[0_26px_70px_rgba(21,72,128,0.14)] lg:order-2"
      />
    </div>
  </section>
</main>
