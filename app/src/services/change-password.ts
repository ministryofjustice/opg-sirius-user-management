export async function changePassword(
  currentPassword: string,
  newPassword: string,
  newPasswordConfirm: string
) {
  const data = new URLSearchParams();
  data.append("existingPassword", currentPassword);
  data.append("password", newPassword);
  data.append("confirmPassword", newPasswordConfirm);

  try {
    const resp = await fetch(
      `${process.env.VUE_APP_SIRIUS_URL}/auth/change-password`,
      {
        mode: "cors",
        method: "POST",
        headers: new Headers({
          "Content-Type": "application/x-www-form-urlencoded"
        }),
        body: data.toString()
      }
    );

    if (resp.ok) {
      return { ok: true };
    } else {
      return { ok: false, error: await resp.json().then(body => body.errors) };
    }
  } catch (err) {
    return { ok: false, error: "something unexpected happened, try again?" };
  }
}
