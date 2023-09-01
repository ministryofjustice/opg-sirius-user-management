import "cypress-axe";
import { TerminalLog } from "../support/e2e";

afterEach(() => {
    cy.injectAxe();
    cy.configureAxe({
      rules: [
        {id: "region", selector: "*:not(.govuk-back-link)"},
        {id: "aria-allowed-attr", selector: "*:not(input[type='radio'][aria-expanded])"},
      ],
    })
    cy.checkA11y(null, null, TerminalLog);
});
