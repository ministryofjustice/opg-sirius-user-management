describe("Add user", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-users": ["post", "put"] });

    cy.addMock("/supervision-api/v1/roles", "GET", {
      status: 200,
      body: ["System Admin"],
    });

    cy.visit("/users");
  });

  it("allows me to add a user", () => {
    cy.contains("a", "Add new user").click();

    cy.get("#f-email").clear().type("123456789");
    cy.get("#f-firstname").clear().type("123456789");
    cy.get("#f-surname").clear().type("123456789");

    cy.addMock("/supervision-api/v1/users", "POST", {
      status: 201,
    });

    cy.get("button[type=submit]").click();

    cy.contains(".moj-alert", "You have successfully added a new user.");
  });
});
