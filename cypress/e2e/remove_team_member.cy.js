describe("Team", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-teams": ["put", "post"] });

    cy.addMock("/api/v1/teams/748", "GET", {
      status: 200,
      body: {
        id: 748,
        displayName: "Finance Team",
        members: [
          {
            displayName: "John Ruecker",
          },
        ],
      },
    });

    cy.visit("/teams/748");
  });

  it("allows me to remove a member", () => {
    cy.get("label[for=f-select-user-0]").click();
    cy.get("button[type=submit]").click();

    cy.url().should("include", "/teams/remove-member/748");
    cy.get(".govuk-body").should(
      "contain",
      "Are you sure you want to remove John Ruecker from the Finance Team team?"
    );

    cy.addMock("/api/v1/teams/748", "PUT", {
      status: 200,
      body: {},
    });

    cy.get("button[type=submit]").click();
    cy.url().should("include", "/teams/748");
  });
});
