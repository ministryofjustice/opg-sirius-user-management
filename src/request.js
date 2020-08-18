export default async function request(
  path,
  method = "GET",
  body = null,
  headers = {}
) {
  const response = await fetch(`http://localhost:8081${path}`, {
    headers,
    body,
    method,
  });

  const json = await response.json();

  return { status: response.status, body: json };
}
