describe("Add team member", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-teams": ["PUT"] });

    cy.addMock("/api/v1/teams/65", "GET", {
      status: 200,
      body: {
        id: 65,
      },
    });

    cy.visit("/teams/add-member/65");
  });

  it("allows me to add a user to a team", () => {
    cy.get(".govuk-table").should("not.exist");

    cy.addMock("/api/v1/search/users?includeSuspended=1&query=admin", "GET", {
      status: 200,
      body: [
        {
          displayName: "system admin",
          email: "system.admin@opgtest.com",
          id: 47,
          surname: "admin",
          suspended: false,
        },
      ],
    });

    cy.get("#f-search").clear().type("admin");
    cy.get("button[type=submit]").click();

    cy.get(".govuk-table__row").should("have.length", 2);

    const expected = [
      "system admin",
      "system.admin@opgtest.com",
      "Add to team",
    ];

    cy.get(".govuk-table__body > .govuk-table__row")
      .children()
      .each(($el, index) => {
        cy.wrap($el).should("contain", expected[index]);
      });

    cy.addMock("/api/v1/teams/65", "PUT", {
      status: 200,
      body: {},
    });

    cy.contains("button", "Add to team").click();

    cy.contains(
      ".moj-alert",
      "You have successfully added system.admin@opgtest.com to the team."
    );
  });
});
