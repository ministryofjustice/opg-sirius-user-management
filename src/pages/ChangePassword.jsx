import React, { useEffect, useMemo, useState } from "react";
import { ErrorSummary } from "govuk-react-jsx";
import ChangePasswordForm from "../forms/ChangePassword";
import Banner from "../components/moj/Banner";
import changePassword from "../api/changePassword";

const ChangePassword = () => {
  const [success, setSuccess] = useState(null);
  const [error, setError] = useState(null);
  const errorList = useMemo(() => [{ children: error }], [error]);

  const onSubmit = async (existingPassword, password, confirmPassword) => {
    const { status, error } = await changePassword(
      existingPassword,
      password,
      confirmPassword
    );
    setSuccess(status < 400);
    setError(error);
  };

  return (
    <div className="govuk-grid-row">
      <div className="govuk-grid-column-two-thirds">
        {error && <ErrorSummary errorList={errorList} />}
        {success && (
          <Banner type="success">Password changed successfully</Banner>
        )}

        <h1 className="govuk-heading-xl">Change password</h1>

        <ChangePasswordForm key={success} onSubmit={onSubmit} />
      </div>
    </div>
  );
};

export default ChangePassword;
