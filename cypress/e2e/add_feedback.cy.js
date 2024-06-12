describe("Feedback", () => {
  beforeEach(() => {
    cy.visit("/supervision/feedback");
  });

  it("allows me to add a feedback", () => {
    cy.get("#f-name").clear().type("New team");
    cy.contains("label[for=f-service-conditional]", "Supervision").click();
    cy.get("#f-supervision-type").select("Allocations");
    cy.get("#f-phone").clear().type("0123045067");
    cy.get("button[type=submit]").click();

    cy.url().should("include", "/teams/123");
  });
});
