describe("Feedback", () => {
  beforeEach(() => {
    cy.setupPermissions();

    cy.addMock("/supervision-api/v1/feedback/supervision", "POST", {
      status: 200,
    });

    cy.visit("/feedback");
  });

  it("allows me to add feedback", () => {
    cy.get("#name").type("Toad McToady");
    cy.get("#email").type("toad@toadhall.com");
    cy.get("#case-number").type("12345");
    cy.get("#more-detail").type("Test feedback");
    cy.get("button[type=submit]").click();
  });
});
