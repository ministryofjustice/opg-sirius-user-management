describe("Delete user", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-users": ["put", "delete"] });

    cy.addMock("/supervision-api/v1/users/123", "GET", {
      status: 200,
      body: {
        firstname: "system",
        surname: "admin",
      },
    });

    cy.visit("/delete-user/123");
  });

  it("allows me to delete a user", () => {
    cy.get(".govuk-body").should(
      "contain",
      "Are you sure you want to delete system admin?"
    );

    cy.addMock("/supervision-api/v1/users/123", "DELETE", {
      status: 200,
    });

    cy.get("button[type=submit]").contains("Delete user").click();

    cy.get('a[href*="/users"]').contains("Continue").click();
    cy.url().should("include", "/users");
  });
});
