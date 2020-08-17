import React from "react";
import ReactDOM from "react-dom";
import "./index.scss";
import { Button } from "govuk-react-jsx";

const App = () => (
  <div class="govuk-body">
    <h1>Change password</h1>
    <Button className="govuk-button--secondary">Continue</Button>
  </div>
);

ReactDOM.render(<App />, document.querySelector("#root"));
