import {
  FlaskConical,
  GraduationCap,
  HandHeart,
  MessageCircle,
  Microscope,
  Ribbon,
  Sparkles,
  UsersRound,
  Utensils
} from 'lucide-svelte';

// Keep homepage copy and assets centralized so the Svelte file only renders sections.
export const navItems = [
  { label: 'Inicio', href: '/#inicio' },
  { label: 'Quiénes somos', href: '/quienes_somos' },
  { label: 'Objetivos', href: '/objetivos' },
  { label: 'Contactenos', href: '/#contactenos' },
  { label: 'Mi cuenta', href: '/client' }
];

export const values = [
  {
    title: 'Empatía',
    text: 'Escuchamos y acompañamos con sensibilidad.',
    icon: HandHeart
  },
  {
    title: 'Ciencia',
    text: 'Promovemos información confiable y basada en evidencia.',
    icon: Microscope
  },
  {
    title: 'Dignidad',
    text: 'Reconocemos el valor de cada persona y su historia.',
    icon: Sparkles
  },
  {
    title: 'Solidaridad',
    text: 'Construimos redes reales de apoyo.',
    icon: HandHeart
  },
  {
    title: 'Esperanza',
    text: 'Creemos en la fuerza de acompañar y seguir adelante.',
    icon: Ribbon
  }
];

export const programs = [
  {
    title: 'Metavida Contigo',
    text: 'Acompañamiento y orientación para pacientes y familias durante el proceso de enfermedad.',
    image: '/images/metavida/encuesta_paciente.jpg'
  },
  {
    title: 'Metavida Escucha',
    text: 'Soporte emocional y espacios de contención para quienes necesitan ser escuchados.',
    icon: MessageCircle
  },
  {
    title: 'Metavida Nutre',
    text: 'Orientación nutricional para fortalecer el bienestar del paciente.',
    icon: Utensils
  },
  {
    title: 'Metavida Aprende',
    text: 'Educación, talleres y recursos para pacientes, cuidadores y comunidad.',
    icon: GraduationCap
  },
  {
    title: 'Metavida Investiga',
    text: 'Investigación científica para generar evidencia y mejores soluciones.',
    icon: FlaskConical
  },
  {
    title: 'Red Solidaria Metavida',
    text: 'Voluntariado, donaciones y apoyo humanitario para quienes más lo necesitan.',
    icon: UsersRound
  }
];

export const supportActions = [
  {
    title: 'Necesito apoyo',
    text: 'Recibe orientación y acompañamiento.',
    image: '/images/metavida/encuesta_toxi1.jpg',
    href: '#encuestas'
  },
  {
    title: 'Quiero donar',
    text: 'Ayúdanos a llegar a más personas.',
    image: '/images/metavida/prevencion-890x660.jpg',
    href: 'https://www.metavida.life/contacto/'
  },
  {
    title: 'Quiero ser voluntario',
    text: 'Súmate a nuestra red solidaria.',
    image: '/images/metavida/uero111-890x660.jpg',
    href: 'https://www.metavida.life/contacto/'
  },
  {
    title: 'Conoce nuestros programas',
    text: 'Descubre cómo podemos acompañarte.',
    image: '/images/metavida/dem111-890x660.jpg',
    href: '#programas'
  }
];

export const surveys = [
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

export const posts = [
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
