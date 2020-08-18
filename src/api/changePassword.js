export default async function changePassword(
  existingPassword,
  password,
  confirmPassword
) {
  const response = await fetch("http://localhost:8081/auth/change-password", {
    method: "POST",
    body: new URLSearchParams({ existingPassword, password, confirmPassword }),
  });

  const json = await response.json();

  return { status: response.status, error: json.errors || null };
}
