<script lang="ts">
  import { POST } from '$libs/http';

  let name = $state('');
  let email = $state('');
  let message = $state('');
  let feedback = $state('');

  async function submitLead() {
    feedback = 'Guardando...';
    try {
      const res = await POST({
        route: '/api/leads',
        data: { name, email, message }
      });
      feedback = `Lead creado: ${res.id}`;
      name = '';
      email = '';
      message = '';
    } catch (err) {
      feedback = err instanceof Error ? err.message : String(err);
    }
  }
</script>

<section>
  <h1>Leads</h1>
  <form class="panel" onsubmit={(event) => { event.preventDefault(); submitLead(); }}>
    <label>
      Nombre
      <input bind:value={name} required />
    </label>
    <label>
      Email
      <input bind:value={email} type="email" />
    </label>
    <label>
      Mensaje
      <textarea bind:value={message} rows="5"></textarea>
    </label>
    <button>Crear lead</button>
    {#if feedback}
      <p>{feedback}</p>
    {/if}
  </form>
</section>

<style>
  h1 {
    margin: 0 0 20px;
  }

  form {
    max-width: 640px;
    display: grid;
    gap: 16px;
    padding: 22px;
  }

  label {
    display: grid;
    gap: 8px;
    color: #2a3532;
    font-weight: 700;
  }

  input,
  textarea {
    width: 100%;
    box-sizing: border-box;
    border: 1px solid #cfd8cf;
    border-radius: 8px;
    padding: 11px 12px;
    background: #fff;
  }

  button {
    width: fit-content;
    min-height: 42px;
    border: 0;
    border-radius: 8px;
    padding: 0 16px;
    background: #15493e;
    color: #fff;
    font-weight: 700;
  }

  p {
    margin: 0;
    color: #53615d;
  }
</style>

