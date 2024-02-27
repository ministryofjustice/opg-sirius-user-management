describe("Team", () => {
  beforeEach(() => {
    cy.visit("/teams/65");
  });

  it("allows me to remove a member", () => {
    cy.get("label[for=f-select-user-0]").click();
    cy.get("button[type=submit]").click();

    cy.url().should("include", "/teams/remove-member/65");
    cy.get(".govuk-body").should(
      "contain",
      "Are you sure you want to remove John from the Cool Team team?"
    );

    cy.get("button[type=submit]").click();
    cy.url().should("include", "/teams/65");
  });
});
