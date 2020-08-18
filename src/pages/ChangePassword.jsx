import React, { useState } from "react";
import { Button, Fieldset, Input } from "govuk-react-jsx";

const ChangePassword = () => {
  const [existingPassword, setExistingPassword] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const submit = async (e) => {
    e.preventDefault();

    const body = new URLSearchParams({
      existingPassword,
      password,
      confirmPassword
    });

    const response = await fetch("http://localhost:8081/auth/change-password", {
      method: "POST",
      body,
    });

    const json = await response.json();

    if (response.status === 200) {
      // success
    } else {
      // failure
    }
  };

  return (
    <div class="govuk-grid-row">
      <div class="govuk-grid-column-two-thirds">
        <h1 class="govuk-heading-xl">Change password</h1>

        <form class="form" action="/my-details" method="post" onSubmit={submit}>
          <Input
            id="existingPassword"
            type="password"
            label={{
              children: "Current password",
              className: "govuk-label--m",
            }}
            value={existingPassword}
            onChange={(e) => setExistingPassword(e.target.value)}
          ></Input>

          <Fieldset
            legend={{
              children: "New password",
              className: "govuk-fieldset__legend--m",
            }}
          >
            <Input
              id="password"
              type="password"
              label={{ children: "Create your new password" }}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            ></Input>
            <Input
              id="confirmPassword"
              type="password"
              label={{ children: "Confirm new password" }}
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
            ></Input>
          </Fieldset>

          <Button type="submit">Save changes</Button>
        </form>
      </div>
    </div>
  );
};

export default ChangePassword;
