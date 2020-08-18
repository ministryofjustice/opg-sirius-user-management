import React, { useReducer } from "react";
import { Button, Fieldset, Input } from "govuk-react-jsx";

const initialState = {
  existingPassword: "",
  password: "",
  confirmPassword: "",
};

const reducer = (state, { field, value }) => ({
  ...state,
  [field]: value,
});

const ChangePassword = ({ onSubmit }) => {
  const [state, dispatch] = useReducer(reducer, initialState);

  const onChange = (e) => {
    dispatch({ field: e.target.name, value: e.target.value });
  };

  const submit = (e) => {
    e.preventDefault();
    onSubmit(state.existingPassword, state.password, state.confirmPassword);
  };

  return (
    <form onSubmit={submit}>
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
  );
};

export default ChangePassword;
