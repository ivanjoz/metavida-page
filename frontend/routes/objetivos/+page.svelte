<script lang="ts">
  import {
    Apple,
    ChartNoAxesCombined,
    Database,
    FlaskConical,
    Globe2,
    HandHeart,
    HeartPulse,
    Laptop,
    Megaphone,
    UsersRound
  } from 'lucide-svelte';
  import { base } from '$app/paths';
  import SiteHeader from '../SiteHeader.svelte';

  // Objectives are copied from the archived WordPress accordion and rendered as scan-friendly cards.
  const objectives = [
    {
      title: 'Investigación científica',
      text: 'Promover y financiar investigaciones científicas sobre cáncer metastásico, fomentando el desarrollo de nuevas terapias, biomarcadores y estrategias de tratamiento a través de alianzas con universidades, centros de investigación y laboratorios nacionales e internacionales.',
      icon: FlaskConical
    },
    {
      title: 'Educación y nutrición oncológica',
      text: 'Implementar campañas de sensibilización sobre el cáncer metastásico dirigidas a la población en general, pacientes, cuidadores y profesionales de la salud, incluyendo educación sobre la importancia de la alimentación en la prevención y tratamiento del cáncer. Se promoverá la capacitación en nutrición oncológica, brindando información sobre dietas especializadas, suplementos y estrategias de soporte nutricional para mejorar la calidad de vida de los pacientes.',
      icon: Apple
    },
    {
      title: 'Redes de soporte',
      text: 'Establecer grupos de apoyo y redes comunitarias para pacientes oncológicos y sus familias, brindando orientación psicológica, legal y social, así como talleres sobre alimentación saludable y cocina terapéutica adaptada a las necesidades de los pacientes.',
      icon: UsersRound
    },
    {
      title: 'Acceso a tratamientos',
      text: 'Gestionar la inclusión de nuevos medicamentos y terapias avanzadas en los sistemas de salud públicos y privados, además de garantizar el acceso a suplementos nutricionales y dietas especializadas para pacientes con cáncer metastásico que lo requieran.',
      icon: HeartPulse
    },
    {
      title: 'Telemedicina',
      text: 'Desarrollar plataformas digitales que permitan a los pacientes recibir orientación médica y nutricional remota, acceso a segundas opiniones de especialistas internacionales y seguimiento de tratamientos.',
      icon: Laptop
    },
    {
      title: 'Participación civil',
      text: 'Promover la creación de voluntariados y programas de participación comunitaria para apoyar a pacientes con cáncer metastásico, incluyendo la entrega de alimentos especializados, suplementos y asesoramiento en hábitos alimenticios saludables.',
      icon: HandHeart
    },
    {
      title: 'Cooperación internacional',
      text: 'Establecer alianzas estratégicas con entidades globales como la OMS, la UICC y otras organizaciones que trabajen en la prevención y tratamiento del cáncer, así como con instituciones dedicadas a la seguridad alimentaria y la nutrición clínica.',
      icon: Globe2
    },
    {
      title: 'Rehabilitación integral',
      text: 'Crear programas especializados en rehabilitación física, nutricional y emocional para pacientes con cáncer metastásico, incluyendo planes de alimentación adaptados a sus necesidades y tratamientos.',
      icon: ChartNoAxesCombined
    },
    {
      title: 'Datos y evidencia',
      text: 'Implementar registros nacionales e internacionales sobre la relación entre la alimentación y la evolución del cáncer metastásico, con el fin de mejorar políticas públicas y estrategias médicas basadas en evidencia.',
      icon: Database
    },
    {
      title: 'Equidad en salud',
      text: 'Garantizar que poblaciones vulnerables, comunidades rurales y personas con bajos recursos económicos tengan acceso a diagnósticos tempranos, tratamientos especializados y asistencia nutricional, incluyendo la entrega de alimentos adecuados para pacientes con cáncer metastásico.',
      icon: Megaphone
    }
  ];
</script>

<svelte:head>
  <title>Objetivos | MetaVida</title>
  <meta
    name="description"
    content="Objetivos específicos de MetaVida para investigación, apoyo a pacientes, nutrición oncológica, acceso a tratamientos y equidad en salud."
  />
</svelte:head>

<main class="min-h-screen bg-[#f4f9ff] text-[#102a57]">
  <SiteHeader />

  <section class="relative isolate overflow-hidden bg-[#0b2246] pt-82 text-white">
    <img src={`${base}/images/metavida/prevencion-890x660.jpg`} alt="" class="absolute inset-0 -z-20 h-full w-full object-cover opacity-42" />
    <div class="absolute inset-0 -z-10 bg-[linear-gradient(115deg,rgba(7,25,56,0.97),rgba(7,67,91,0.84)_60%,rgba(24,198,130,0.34))]"></div>

    <div class="mx-auto grid min-h-[520px] w-[min(1180px,calc(100%-32px))] content-end py-78">
      <div class="max-w-[900px]">
        <p class="max-w-max rounded-full bg-[#18c682] px-18 py-8 text-sm font-bold uppercase leading-[1.2] text-[#073953]">
          Objetivos
        </p>
        <h1 class="mt-24 text-[42px] font-bold leading-[1.02] text-white md:text-[70px]">
          Acciones concretas para mejorar la vida de pacientes y familias.
        </h1>
        <p class="mt-22 max-w-[720px] text-lg leading-[1.6] text-white/84 md:text-xl">
          La Asociación Metavida ha establecido objetivos específicos que articulan investigación, educación, soporte comunitario, acceso a tratamientos y equidad en salud.
        </p>
      </div>
    </div>
  </section>

  <section class="bg-white py-86">
    <div class="mx-auto grid w-[min(1180px,calc(100%-32px))] gap-36 lg:grid-cols-[0.78fr_1.22fr] lg:items-start">
      <div class="lg:sticky lg:top-116">
        <p class="text-base font-bold uppercase leading-[1.2] text-[#18a878]">Plan de trabajo</p>
        <h2 class="mt-8 text-[36px] font-bold uppercase leading-[1.02] text-[#0784a7] md:text-[58px]">
          Diez objetivos específicos
        </h2>
        <div class="mt-14 h-5 w-52 bg-[#18c682]"></div>
        <p class="mt-22 text-lg leading-[1.65] text-[#27384a]">
          Cada objetivo conecta una necesidad real del paciente oncológico con una línea de acción institucional medible.
        </p>
      </div>

      <div class="grid gap-16">
        {#each objectives as objective, index}
          {@const Icon = objective.icon}
          <article class="grid gap-18 rounded-[8px] border border-[#dbeaf8] bg-[#f8fcff] p-22 shadow-[0_12px_30px_rgba(21,72,128,0.07)] md:grid-cols-[72px_1fr]">
            <div class="grid h-58 w-58 place-items-center rounded-[8px] bg-[#e8fff5] text-[#128b65]">
              <Icon size={28} strokeWidth={1.8} />
            </div>
            <div>
              <p class="text-sm font-bold uppercase leading-[1.2] text-[#8ba0b8]">Objetivo {index + 1}</p>
              <h3 class="mt-4 text-2xl font-bold leading-tight text-[#075f80]">{objective.title}</h3>
              <p class="mt-10 text-base leading-[1.62] text-[#27384a]">{objective.text}</p>
            </div>
          </article>
        {/each}
      </div>
    </div>
  </section>
</main>
