describe("Edit a team", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-teams": ["put", "post"] });

    cy.addMock("/api/v1/teams/837", "GET", {
      status: 200,
      body: {
        id: 837,
        displayName: "Finance Team",
        teamType: { handle: "FINANCE", label: "Finance" },
        phoneNumber: "01818118181",
        email: "finance.team@opgtest.com",
      },
    });

    cy.addMock("/api/v1/reference-data?filter=teamType", "GET", {
      status: 200,
      body: {
        teamType: [
          {
            handle: "ALLOCATIONS",
            label: "Allocations",
          },
          {
            handle: "FINANCE",
            label: "Finance",
          },
        ],
      },
    });

    cy.visit("/teams/edit/837");
  });

  it("shows the team details", () => {
    cy.get("#f-name").should("have.value", "Finance Team");
    cy.get("#f-service-conditional").should("be.checked");
    cy.get("#f-service-conditional-2").should("not.be.checked");
    cy.get("#f-type").should("have.value", "FINANCE");
    cy.get("#f-phoneNumber").should("have.value", "01818118181");
    cy.get("#f-email").should("have.value", "finance.team@opgtest.com");
  });

  it("allows me to change the team's details", () => {
    cy.get("#f-name").clear().type("Allocations team");
    cy.get("#f-type").select("ALLOCATIONS");
    cy.get("#f-phoneNumber").clear().type("03573953");
    cy.get("#f-email").clear().type("other.team@opgtest.com");

    cy.addMock("/api/v1/teams/837", "PUT", {
      status: 200,
      body: {},
    });

    cy.get("button[type=submit]").click();

    cy.contains(
      ".moj-banner",
      "You have successfully edited Allocations team."
    );
  });
});
