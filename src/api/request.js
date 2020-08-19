export default async function request(
  path,
  method = "GET",
  body = null,
  headers = {}
) {
  const baseUrl = globalThis.baseUrl || "http://localhost:8081";
  const response = await fetch(`${baseUrl}${path}`, {
    headers,
    body,
    method,
  });

  const json = await response.json();

  return { status: response.status, body: json };
}
