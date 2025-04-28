describe("Teams", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-teams": ["post"] });

    cy.addMock("/supervision-api/v1/reference-data?filter=teamType", "GET", {
      status: 200,
      body: {
        teamType: [
          {
            handle: "ALLOCATIONS",
            label: "Allocations",
          },
        ],
      },
    });

    cy.visit("/teams/add");
  });

  it("allows me to add a new team", () => {
    cy.get("#f-name").clear().type("New team");
    cy.contains("label[for=f-service-conditional]", "Supervision").click();
    cy.get("#f-supervision-type").select("Allocations");
    cy.get("#f-phone").clear().type("0123045067");

    cy.addMock("/supervision-api/v1/teams", "POST", {
      status: 201,
      body: { id: 123 },
    });

    cy.get("button[type=submit]").click();

    cy.url().should("include", "/teams/123");
  });
});
