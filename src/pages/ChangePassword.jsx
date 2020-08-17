import React from "react";
import { Button, Fieldset, Input } from "govuk-react-jsx";

const ChangePassword = () => (
  <div class="govuk-grid-row">
    <div class="govuk-grid-column-two-thirds">
      <h1 class="govuk-heading-xl">Change password</h1>

      <form
        class="form"
        action="/my-details"
        method="post"
        data-bitwarden-watching="1"
      >
        <Input
          id="currentPassword"
          label={{
            children: "Current password",
            className: "govuk-label--m",
          }}
        ></Input>

        <Fieldset
          legend={{
            children: "New password",
            className: "govuk-fieldset__legend--m",
          }}
        >
          <Input
            id="password"
            label={{ children: "Create your new password" }}
          ></Input>
          <Input
            id="confirmPassword"
            label={{ children: "Confirm new password" }}
          ></Input>
        </Fieldset>

        <Button>Save changes</Button>
      </form>
    </div>
  </div>
);

export default ChangePassword;
