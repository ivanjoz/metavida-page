<script lang="ts">
  import { GET } from '$libs/http';

  let status = $state('Cargando...');

  GET('/health')
    .then((res) => {
      status = res.status || 'ok';
    })
    .catch((err) => {
      status = err.message;
    });
</script>

<section>
  <h1>Panel administrativo</h1>
  <div class="grid">
    <article class="panel">
      <strong>Backend</strong>
      <span>{status}</span>
    </article>
    <article class="panel">
      <strong>Leads</strong>
      <span>Modulo inicial</span>
    </article>
    <article class="panel">
      <strong>Base de datos</strong>
      <span>Cloudflare D1</span>
    </article>
  </div>
</section>

<style>
  h1 {
    margin: 0 0 24px;
    color: #17211f;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 16px;
  }

  article {
    display: grid;
    gap: 8px;
    padding: 20px;
  }

  strong {
    font-size: 18px;
  }

  span {
    color: #596662;
  }

  @media (max-width: 760px) {
    .grid {
      grid-template-columns: 1fr;
    }
  }
</style>

