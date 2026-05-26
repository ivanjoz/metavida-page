<script lang="ts">
  import { onMount } from 'svelte';
  import Input from '$components/form/Input.svelte';
  import Button from '$components/buttons/Button.svelte';
  import { Notify } from '$libs/helpers';
  import { checkIsLogin, Env } from '$core/env';

  let form = $state({ email: '', password: '' });
  let isLoading = $state(false);

  onMount(() => {
    if (checkIsLogin() === 2) {
      window.location.href = '/admin';
    }
  });

  const sendLogin = async () => {
    if (form.email.length < 4 || form.password.length < 4) {
      Notify.failure('Debe proporcionar un email y una contraseña válidos');
      return;
    }

    isLoading = true;
    Notify.info('Enviando credenciales...');

    try {
      const res = await fetch(Env.makeRoute('p-user-login'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: form.email, password: form.password })
      });

      const body = await res.json();

      if (!res.ok) {
        Notify.failure(body.error || body.message || 'Credenciales inválidas');
        return;
      }

      Env.setSession(body.UserToken, body.TokenExpTime, body.UserInfo || '');
      window.location.href = '/admin';
    } catch {
      Notify.failure('Error de conexión. Intente nuevamente.');
    } finally {
      isLoading = false;
    }
  };
</script>

<div class="login-wrap">
  <div class="login-box">
    <div class="login-tt">Iniciar Sesión</div>

    <a class="login-brand" href="/">
      <img src="/images/metavida/logo.svg" alt="Metavida" />
    </a>

    <Input
      css="mb-12 w-full"
      label="Email"
      saveOn={form}
      save="email"
      type="text"
      required={true}
    />

    <Input
      css="mb-12 w-full"
      label="Contraseña"
      saveOn={form}
      save="password"
      type="password"
      required={true}
    />

    <div class="flex justify-center mt-16">
      <Button
        name="Ingresar"
        color="green"
        disabled={isLoading}
        onClick={() => sendLogin()}
      />
    </div>
  </div>
</div>

<style>
  .login-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: #eef7ff;
  }

  .login-box {
    position: relative;
    width: 28rem;
    background: white;
    border-radius: 14px;
    border-top: 4px solid #15493e;
    padding: 2rem 2rem 2.4rem;
    box-shadow: rgba(17, 17, 26, 0.1) 0px 4px 16px,
                rgba(17, 17, 26, 0.08) 0px 8px 24px;
  }

  .login-tt {
    position: absolute;
    top: -2.2rem;
    left: 2rem;
    height: 2.2rem;
    padding: 0 12px;
    background: #15493e;
    color: white;
    border-radius: 10px 10px 0 0;
    display: flex;
    align-items: center;
    font-size: 15px;
  }

  .login-brand {
    display: flex;
    justify-content: center;
    margin-bottom: 1.6rem;
  }

  .login-brand img {
    height: 52px;
    object-fit: contain;
  }

  @media (max-width: 520px) {
    .login-box {
      width: 90vw;
      padding: 1.2rem 1.2rem 1.6rem;
    }
  }
</style>
