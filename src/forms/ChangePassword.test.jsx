/* eslint-env jest */
import React from "react";
import { render } from "@testing-library/react";
import ChangePassword from "./ChangePassword.jsx";

describe("the ChangePassword form is rendered", () => {
  test("it should render three inputs", () => {
    const { getByLabelText } = render(<ChangePassword />);
    [
      "Current password",
      "Create your new password",
      "Confirm new password",
    ].forEach((label) => {
      const input = getByLabelText(label);
      expect(input).toBeInTheDocument();
      expect(input).toHaveValue("");
    });
  });

  test("it should render a submit button", (done) => {
    const { getByText } = render(<ChangePassword />);
    const button = getByText("Save changes");
    expect(button).toBeInTheDocument();
    expect(button).toHaveClass("govuk-button");

    done();
  });
});
