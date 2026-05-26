<script lang="ts">
  import {
    ArrowRight,
    CalendarDays,
    HandHeart,
    HeartPulse,
    Landmark,
    Mail,
    Phone,
    Send,
    ShieldCheck,
    UsersRound
  } from 'lucide-svelte';
  import SiteHeader from './SiteHeader.svelte';
  import { posts, programs, supportActions, surveys, values } from './page.content';

  let contactName = $state('');
  let contactEmail = $state('');
  let contactPhone = $state('');
  let contactSubject = $state('');
  let contactMessage = $state('');

  function submitContactForm() {
    // Use mailto until the backend contact endpoint exists; logs make form behavior traceable in pre-alpha.
    console.debug('Submitting public contact form', {
      contactName,
      contactEmail,
      contactPhone,
      contactSubject,
      hasMessage: contactMessage.trim().length > 0
    });

    const emailBody = [
      `Nombre: ${contactName}`,
      `Email: ${contactEmail}`,
      `Telefono: ${contactPhone}`,
      '',
      contactMessage
    ].join('\n');

    const encodedSubject = encodeURIComponent(contactSubject || 'Contacto desde MetaVida');
    const encodedBody = encodeURIComponent(emailBody);
    window.location.href = `mailto:informes@metavida.life?subject=${encodedSubject}&body=${encodedBody}`;
  }
</script>

<svelte:head>
  <title>MetaVida</title>
  <meta
    name="description"
    content="MetaVida une esfuerzos por el bienestar, la salud y el acompañamiento de pacientes."
  />
</svelte:head>

<main class="min-h-screen bg-[#eef7ff] text-[#102a57]">
  <SiteHeader />

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
      <h1 class="max-w-[900px] text-[42px] font-semibold leading-[1.03] text-white [text-shadow:0_4px_24px_rgba(0,0,0,0.55)] md:text-[68px] lg:text-[84px]">
        Uniendo esfuerzos por <span class="text-[#8cc9ff]">el bienestar y la salud.</span>
      </h1>
      <div class="flex max-w-[820px] flex-col gap-24 md:flex-row md:items-center md:justify-between">
        <p class="max-w-[590px] text-lg leading-[1.55] text-white/90 [text-shadow:0_3px_16px_rgba(0,0,0,0.5)] md:text-xl">
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

  <section class="bg-[#edf4fb] py-92">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-42 lg:grid-cols-[0.9fr_1.1fr] lg:items-center">
      <div>
        <p class="text-[30px] font-semibold leading-tight text-[#16bf7e] md:text-[38px]">Acompañando con ciencia</p>
        <h2 class="mt-6 max-w-[520px] text-[44px] font-bold uppercase leading-[0.98] text-[#0784a7] md:text-[66px]">
          Esperanza y humanidad
        </h2>
        <div class="mt-16 h-5 w-54 bg-[#18c682]"></div>
      </div>

      <div class="grid gap-28">
        <article class="grid gap-18 border-b border-[#18c682] pb-28 md:grid-cols-[80px_1fr]">
          <div class="grid h-72 w-72 place-items-center text-[#0784a7]">
            <HandHeart size={62} strokeWidth={1.7} />
          </div>
          <p class="text-xl leading-[1.55] text-[#1d2a39]">
            <strong>MetaVida nace para acompañar a personas</strong> que enfrentan el cáncer y otras enfermedades complejas, junto a sus familias y cuidadores.
          </p>
        </article>
        <article class="grid gap-18 md:grid-cols-[80px_1fr]">
          <div class="grid h-72 w-72 place-items-center text-[#0784a7]">
            <UsersRound size={60} strokeWidth={1.7} />
          </div>
          <p class="text-xl leading-[1.55] text-[#1d2a39]">
            <strong>Creemos que nadie debería atravesar una enfermedad difícil en soledad.</strong> Por eso brindamos orientación, apoyo humano, educación e investigación para mejorar la calidad de vida de quienes más lo necesitan.
          </p>
        </article>
      </div>
    </div>
  </section>

  <section id="valores" class="bg-[linear-gradient(180deg,#008c96_0%,#082a54_100%)] py-86 text-white">
    <div class="mx-auto w-[min(1180px,calc(100%-32px))]">
      <div class="mb-50 text-center">
        <h2 class="text-[40px] font-bold uppercase leading-tight md:text-[64px]">Nuestros valores</h2>
        <div class="mx-auto mt-14 h-5 w-42 bg-[#18c682]"></div>
      </div>

      <div class="grid gap-30 md:grid-cols-3">
        {#each values as value}
          {@const Icon = value.icon}
          <article class="flex flex-col items-center gap-24 text-center">
            <div class="grid h-86 w-86 place-items-center rounded-full border-2 border-[#18c682] text-[#18c682]">
              <Icon size={38} strokeWidth={1.7} />
            </div>
            <div>
              <h3 class="text-xl font-bold uppercase leading-tight">{value.title}</h3>
              <p class="mt-8 text-base leading-[1.45] text-white/78">{value.text}</p>
            </div>
          </article>
        {/each}
      </div>
    </div>
  </section>

  <section class="bg-white py-96">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-42 lg:grid-cols-[0.86fr_1fr] lg:items-center">
      <div>
        <h2 class="max-w-[520px] text-[44px] font-bold uppercase leading-[1.02] text-[#0784a7] md:text-[62px]">
          <span class="text-[#18c682]">¿A quiénes</span> acompañamos?
        </h2>
        <div class="mt-14 h-5 w-42 bg-[#18c682]"></div>
        <div class="mt-34 grid gap-22 text-lg leading-[1.65] text-[#27384a]">
          <p>
            Acompañamos a pacientes oncológicos, personas con enfermedades crónicas complejas, sobrevivientes de cáncer, familiares, cuidadores y personas en situación de vulnerabilidad.
          </p>
          <p>
            También trabajamos con profesionales de la salud, universidades, investigadores, voluntarios, donantes e instituciones aliadas.
          </p>
        </div>
      </div>

      <img
        src="/images/metavida/encuesta_paciente.jpg"
        alt="Equipo de salud acompañando a pacientes"
        class="min-h-360 w-full rounded-[8px] object-cover shadow-[0_22px_58px_rgba(21,72,128,0.12)]"
      />
    </div>
  </section>

  <section id="programas" class="bg-white pb-96">
    <div class="mx-auto w-[min(1180px,calc(100%-32px))]">
      <div class="mb-44 text-center">
        <h2 class="text-[40px] font-bold uppercase leading-tight text-[#0784a7] md:text-[64px]">
          Nuestros <span class="text-[#18c682]">programas</span>
        </h2>
        <div class="mx-auto mt-14 h-5 w-42 bg-[#18c682]"></div>
      </div>

      <div class="grid gap-18 md:grid-cols-2 lg:grid-cols-3">
        {#each programs as program}
          <article class="relative min-h-286 overflow-hidden rounded-[8px] bg-[#f0f2f5] shadow-[0_14px_36px_rgba(21,72,128,0.08)]">
            {#if program.image}
              <img src={program.image} alt="" class="absolute inset-0 h-full w-full object-cover" />
              <div class="absolute inset-0 bg-[linear-gradient(180deg,rgba(6,23,48,0)_34%,rgba(4,39,72,0.9)_100%)]"></div>
            {:else}
              {@const Icon = program.icon}
              <div class="absolute inset-x-0 top-54 flex justify-center text-[#0a7192]">
                <Icon size={64} strokeWidth={1.5} />
              </div>
            {/if}
            <div class="absolute inset-x-0 bottom-0 p-22 {program.image ? 'text-white' : 'text-[#102a57]'}">
              <h3 class="text-xl font-bold uppercase leading-tight text-[#0784a7] {program.image ? '!text-white' : ''}">
                {program.title}
              </h3>
              <div class="mt-4 h-3 w-34 bg-[#18c682]"></div>
              <p class="mt-8 text-base leading-[1.45] {program.image ? 'text-white/88' : 'text-[#27384a]'}">
                {program.text}
              </p>
            </div>
          </article>
        {/each}
      </div>
    </div>
  </section>

  <section id="apoyo" class="bg-[#eef4fb] py-84">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-22 md:grid-cols-2">
      {#each supportActions as action}
        <article class="grid gap-22 rounded-[8px] bg-white p-12 shadow-[0_14px_36px_rgba(21,72,128,0.11)] md:grid-cols-[minmax(190px,0.9fr)_1fr] md:items-center">
          <img src={action.image} alt="" class="h-220 w-full rounded-[6px] object-cover md:h-full" />
          <div class="p-12 md:p-20">
            <h3 class="text-2xl font-bold uppercase leading-tight text-[#075f80]">{action.title}</h3>
            <div class="mt-4 h-3 w-34 bg-[#18c682]"></div>
            <p class="mt-10 text-base leading-[1.55] text-[#27384a]">{action.text}</p>
            <a
              href={action.href}
              class="mt-18 inline-flex h-40 items-center justify-center rounded-[6px] bg-[#075f80] px-16 text-sm font-bold !text-white transition hover:bg-[#18a878]"
            >
              Quiero saber más
            </a>
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

  <section class="bg-[#17375f] py-78 text-white">
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

  <section id="contactenos" class="bg-[#eef4fb] py-92">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-30 lg:grid-cols-[0.82fr_1.18fr] lg:items-start">
      <div class="grid gap-22">
        <div>
          <p class="max-w-max rounded-full bg-white px-22 py-9 text-sm font-bold uppercase leading-[1.2] text-[#2e7dd7] shadow-[0_10px_28px_rgba(21,72,128,0.08)]">
            Contactenos
          </p>
          <h2 class="mt-18 text-[36px] font-bold leading-tight text-[#102a57] md:text-[54px]">
            Escríbenos y conversemos
          </h2>
          <p class="mt-16 max-w-[520px] text-lg leading-[1.65] text-[#596a82]">
            Completa el formulario o usa los datos directos de contacto para coordinar apoyo, alianzas, voluntariado o donaciones.
          </p>
        </div>

        <div class="grid gap-14">
          <a
            href="mailto:informes@metavida.life"
            class="grid gap-14 rounded-[8px] border border-[#dbeaf8] bg-white p-20 shadow-[0_12px_30px_rgba(21,72,128,0.07)] md:grid-cols-[52px_1fr]"
          >
            <span class="grid h-48 w-48 place-items-center rounded-[8px] bg-[#e8fff5] text-[#128b65]">
              <Mail size={24} />
            </span>
            <span>
              <span class="block text-sm font-bold uppercase leading-[1.2] text-[#8ba0b8]">Escribenos</span>
              <span class="mt-4 block text-lg font-bold leading-tight text-[#075f80]">informes@metavida.life</span>
            </span>
          </a>

          <a
            href="tel:+51945364062"
            class="grid gap-14 rounded-[8px] border border-[#dbeaf8] bg-white p-20 shadow-[0_12px_30px_rgba(21,72,128,0.07)] md:grid-cols-[52px_1fr]"
          >
            <span class="grid h-48 w-48 place-items-center rounded-[8px] bg-[#e8fff5] text-[#128b65]">
              <Phone size={24} />
            </span>
            <span>
              <span class="block text-sm font-bold uppercase leading-[1.2] text-[#8ba0b8]">Llamanos</span>
              <span class="mt-4 block text-lg font-bold leading-tight text-[#075f80]">(+51) 945 364 062</span>
            </span>
          </a>
        </div>
      </div>

      <form
        class="grid gap-16 rounded-[8px] bg-white p-22 shadow-[0_22px_52px_rgba(21,72,128,0.11)] md:p-30"
        onsubmit={(event) => {
          event.preventDefault();
          submitContactForm();
        }}
      >
        <div class="grid gap-16 md:grid-cols-2">
          <label class="grid gap-7">
            <span class="text-sm font-bold uppercase leading-[1.2] text-[#596a82]">Tu Nombre *</span>
            <input
              class="h-48 rounded-[6px] border border-[#d3e1ef] bg-[#f8fcff] px-14 text-base text-[#102a57] outline-none transition focus:border-[#18a878] focus:bg-white"
              name="name"
              required
              autocomplete="name"
              bind:value={contactName}
            />
          </label>

          <label class="grid gap-7">
            <span class="text-sm font-bold uppercase leading-[1.2] text-[#596a82]">Tu Email *</span>
            <input
              class="h-48 rounded-[6px] border border-[#d3e1ef] bg-[#f8fcff] px-14 text-base text-[#102a57] outline-none transition focus:border-[#18a878] focus:bg-white"
              name="email"
              type="email"
              required
              autocomplete="email"
              bind:value={contactEmail}
            />
          </label>
        </div>

        <div class="grid gap-16 md:grid-cols-2">
          <label class="grid gap-7">
            <span class="text-sm font-bold uppercase leading-[1.2] text-[#596a82]">Tu Teléfono *</span>
            <input
              class="h-48 rounded-[6px] border border-[#d3e1ef] bg-[#f8fcff] px-14 text-base text-[#102a57] outline-none transition focus:border-[#18a878] focus:bg-white"
              name="phone"
              type="tel"
              required
              autocomplete="tel"
              bind:value={contactPhone}
            />
          </label>

          <label class="grid gap-7">
            <span class="text-sm font-bold uppercase leading-[1.2] text-[#596a82]">Asunto</span>
            <input
              class="h-48 rounded-[6px] border border-[#d3e1ef] bg-[#f8fcff] px-14 text-base text-[#102a57] outline-none transition focus:border-[#18a878] focus:bg-white"
              name="subject"
              bind:value={contactSubject}
            />
          </label>
        </div>

        <label class="grid gap-7">
          <span class="text-sm font-bold uppercase leading-[1.2] text-[#596a82]">Mensaje *</span>
          <textarea
            class="min-h-150 resize-y rounded-[6px] border border-[#d3e1ef] bg-[#f8fcff] p-14 text-base text-[#102a57] outline-none transition focus:border-[#18a878] focus:bg-white"
            name="message"
            required
            bind:value={contactMessage}
          ></textarea>
        </label>

        <button
          type="submit"
          class="inline-flex h-52 w-max items-center justify-center gap-10 rounded-[6px] bg-[#075f80] px-22 text-base font-bold text-white transition hover:bg-[#18a878]"
        >
          Enviar
          <Send size={18} />
        </button>
      </form>
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
