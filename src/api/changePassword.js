import request from "./request";

export default async function changePassword(
  existingPassword,
  password,
  confirmPassword
) {
  const { status, body } = await request(
    "/auth/change-password",
    "POST",
    new URLSearchParams({ existingPassword, password, confirmPassword })
  );

  return {
    status,
    error: (body && body.errors) || null,
  };
}
