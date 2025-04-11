describe("Users", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-users": ["put"] });

    cy.visit("/users");
  });

  it("allows me to search for not in a team", () => {
    cy.addMock("/api/v1/search/users?includeSuspended=1&query=admin", "GET", {
      status: 200,
      body: [
        {
          displayName: "system admin",
          email: "system.admin@opgtest.com",
          teams: [],
        },
      ],
    });

    const expected = [
      "system admin",
      "",
      "system.admin@opgtest.com",
      "Active",
      "Edit",
    ];
    search("admin", expected);
  });

  it("allows me to search for a user in a team", () => {
    cy.addMock("/api/v1/search/users?includeSuspended=1&query=anton", "GET", {
      status: 200,
      body: [
        {
          displayName: "Anton Mccoy",
          email: "anton.mccoy@opgtest.com",
          teams: [
            {
              displayName: "Visits Team",
            },
          ],
          suspended: true,
        },
      ],
    });

    const expected = [
      "Anton Mccoy",
      "Visits Team",
      "anton.mccoy@opgtest.com",
      "Suspended",
      "Edit",
    ];
    search("anton", expected);
  });

  function search(searchTerm, expected) {
    cy.get(".govuk-table").should("not.exist");

    cy.get("#f-search").clear().type(searchTerm);
    cy.get("button[type=submit]").click();

    cy.get(".govuk-table__row").should("have.length", 2);

    cy.get(".govuk-table__body > .govuk-table__row")
      .children()
      .each(($el, index) => {
        cy.wrap($el).should("contain", expected[index]);
      });
  }
});
