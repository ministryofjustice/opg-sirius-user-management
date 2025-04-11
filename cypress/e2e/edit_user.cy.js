describe("Edit user", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-users": ["put"] });

    cy.addMock("/api/v1/roles", "GET", {
      status: 200,
      body: ["System Admin", "Finance", "Self-Allocation", "File Creation"],
    });

    cy.addMock("/api/v1/users/123", "GET", {
      status: 200,
      body: {
        firstname: "Hadley",
        surname: "Collins",
        email: "h.collins@opg.example",
        roles: ["OPG User", "Finance", "Self-Allocation"],
      },
    });

    cy.visit("/edit-user/123");
  });

  it("allows me to edit a user", () => {
    cy.get("#f-firstname").should("have.value", "Hadley");
    cy.get("#f-surname").should("have.value", "Collins");
    cy.get("#f-email").should("have.value", "h.collins@opg.example");

    cy.get("[name='organisation'][value='COP User']").should("not.be.checked");
    cy.get("[name='organisation'][value='OPG User']").should("be.checked");
    cy.get("[name='roles'][value='Finance']").should("be.checked");
    cy.get("[name='roles'][value='System Admin']").should("not.be.checked");

    cy.get("#f-firstname").type("Abe");
    cy.get("[name='roles'][value='System Admin']").check();

    cy.addMock("/api/v1/users/123", "PUT", {
      status: 200,
      body: {},
    });

    cy.get("button[type=submit]").click();

    cy.contains(".moj-alert", "You have successfully edited a user.");
  });
});
