describe("Delete a team", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-teams": ["delete"] });

    cy.addMock("/api/v1/teams/65", "GET", {
      status: 200,
      body: {
        displayName: "Deletion Test Team",
      },
    });

    cy.visit("/teams/delete/65");
  });

  it("shows the team details", () => {
    cy.get(".govuk-body").should(
      "contain",
      "Are you sure you want to delete the team Deletion Test Team?"
    );
  });

  it("allows me to delete the team", () => {
    cy.addMock("/api/v1/teams/65", "DELETE", {
      status: 204,
    });

    cy.contains("button", "Delete team").click();
    cy.url().should("include", "/teams");
  });
});
