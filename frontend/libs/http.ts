import { Env } from '$core/env';

type RequestProps = {
  route: string;
  data?: unknown;
  headers?: Record<string, string>;
};

async function request(method: string, props: RequestProps) {
  const res = await fetch(Env.makeRoute(props.route), {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...(props.headers || {})
    },
    body: props.data === undefined ? undefined : JSON.stringify(props.data)
  });

  const contentType = res.headers.get('content-type') || '';
  const body = contentType.includes('application/json') ? await res.json() : await res.text();

  if (!res.ok) {
    const message = typeof body === 'string' ? body : body.error || body.message || 'Request failed';
    throw new Error(message);
  }
  return body;
}

export const GET = (route: string) => request('GET', { route });
export const POST = (props: RequestProps) => request('POST', props);
export const PUT = (props: RequestProps) => request('PUT', props);

