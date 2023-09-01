
import "cypress-axe";
import { TerminalLog } from "../support/e2e";
import navTabs from "../fixtures/navigation.json"

describe("Accessibility", () => {
    navTabs.forEach(([page, url]) =>
        it(`should render ${page} page accessibly`, () => {
            cy.visit(url);
            cy.injectAxe();
            cy.checkA11y(null, null, TerminalLog);
        })
    )
});