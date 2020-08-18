import React, { useEffect, useReducer, useMemo, useState } from "react";
import { Button, Fieldset, Input } from "govuk-react-jsx";
import Banner from "../components/moj/Banner";
import { ErrorSummary } from "govuk-react-jsx/govuk/components/error-summary";
import changePassword from "../api/changePassword";

const initialState = {
  existingPassword: "",
  password: "",
  confirmPassword: "",
};

const reducer = (state, { field, value }) => ({
  ...state,
  [field]: value,
});

const ChangePassword = () => {
  const [state, dispatch] = useReducer(reducer, initialState);
  const [success, setSuccess] = useState(null);
  const [error, setError] = useState(null);
  const errorList = useMemo(() => [{ children: error }], [error]);

  const onSubmit = async (e) => {
    e.preventDefault();

    const { status, error } = await changePassword(
      state.existingPassword,
      state.password,
      state.confirmPassword
    );
    setSuccess(status < 400);
    setError(error);
  };

  const onChange = (e) => {
    dispatch({ field: e.target.name, value: e.target.value });
  };

  useEffect(() => {
    if (success) {
      Object.keys(state).forEach((field) => {
        dispatch({ field, value: "" });
      });
    }
  }, [success]);

  return (
    <div className="govuk-grid-row">
      <div className="govuk-grid-column-two-thirds">
        {error && <ErrorSummary errorList={errorList} />}
        {success && (
          <Banner type="success">Password changed successfully</Banner>
        )}

        <h1 className="govuk-heading-xl">Change password</h1>

        <form
          className="form"
          action="/my-details"
          method="post"
          onSubmit={onSubmit}
        >
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
