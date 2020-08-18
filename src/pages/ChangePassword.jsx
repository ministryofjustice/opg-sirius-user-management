import React, { useEffect, useMemo, useState } from "react";
import { ErrorSummary } from "govuk-react-jsx";
import ChangePasswordForm from '../forms/ChangePassword'
import Banner from "../components/moj/Banner";
import request from "../request";

const ChangePassword = () => {
  const [success, setSuccess] = useState(null);
  const [error, setError] = useState(null);
  const errorList = useMemo(() => [{ children: error }], [error]);

  const onSubmit = async (existingPassword, password, confirmPassword) => {
    const { status, body } = await request(
      "/auth/change-password",
      "POST",
      new URLSearchParams({ existingPassword, password, confirmPassword })
    );
    setSuccess(status < 400);
    setError((body && body.errors) || null);
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
