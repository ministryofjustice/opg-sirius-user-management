describe("Feedback", () => {
  beforeEach(() => {
    cy.visit("/supervision/feedback");
  });

  it("allows me to add a feedback", () => {
    cy.get("#name").type("Mr Toad");
    cy.get("#email").type("toad@toadmail.com");
    cy.get("#case-number").type("123456");
    cy.get("#more-detail").type("I have some thoughts to feedback");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "/supervision/feedback");

    cy.get("#govuk-notification-banner-title").should("be.visible");
  });
});
