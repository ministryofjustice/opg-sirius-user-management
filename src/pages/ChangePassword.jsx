import React, { useReducer } from "react";
import { Button, Fieldset, Input } from "govuk-react-jsx";

const initialState = {
  existingPassword: '',
  password: '',
  confirmPassword: ''
}

const reducer = (state, {field, value}) => ({
  ...state,
  [field]: value
})

const ChangePassword = () => {
  const [state, dispatch] = useReducer(reducer, initialState)

  const onSubmit = async (e) => {
    e.preventDefault();

    const body = new URLSearchParams(state);

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

  const onChange = (e) => {
    dispatch({ field: e.target.name, value: e.target.value })
  }

  return (
    <div class="govuk-grid-row">
      <div class="govuk-grid-column-two-thirds">
        <h1 class="govuk-heading-xl">Change password</h1>

        <form class="form" action="/my-details" method="post" onSubmit={onSubmit}>
          <Input
            id="existingPassword"
            name="existingPassword"
            type="password"
            label={{
              children: "Current password",
              className: "govuk-label--m",
            }}
            value={state.existingPassword}
            onChange={onChange}
          ></Input>

          <Fieldset
            legend={{
              children: "New password",
              className: "govuk-fieldset__legend--m",
            }}
          >
            <Input
              id="password"
              name="password"
              type="password"
              label={{ children: "Create your new password" }}
              value={state.password}
              onChange={onChange}
            ></Input>
            <Input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              label={{ children: "Confirm new password" }}
              value={state.confirmPassword}
              onChange={onChange}
            ></Input>
          </Fieldset>

          <Button type="submit">Save changes</Button>
        </form>
      </div>
    </div>
  );
};

export default ChangePassword;
