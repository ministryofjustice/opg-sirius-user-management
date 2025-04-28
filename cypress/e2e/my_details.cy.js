describe("My details", () => {
  it("shows my details", () => {
    cy.setupPermissions({ "v1-users-updatetelephonenumber": ["put"] });

    cy.addMock("/supervision-api/v1/users/current", "GET", {
      status: 200,
      body: {
        firstname: "system",
        surname: "admin",
        email: "system.admin@opgtest.com",
        phoneNumber: "03004560300",
        roles: ["OPG User", "Finance", "System Admin"],
        teams: [
          {
            displayName: "Administrative Team",
          },
        ],
      },
    });

    cy.visit("/my-details");

    const expected = [
      ["Name", "system admin"],
      ["Email", "system.admin@opgtest.com"],
      ["Phone number", "03004560300"],
      ["Organisation", "OPG User"],
      ["Team", "Administrative Team"],
      ["Roles", "Finance, System Admin"],
    ];

    cy.get(".govuk-summary-list__row").each(($el, index) => {
      cy.wrap($el).within(() => {
        cy.get(".govuk-summary-list__key").should(
          "have.text",
          expected[index][0]
        );
        cy.get(".govuk-summary-list__value").should(
          "have.text",
          expected[index][1]
        );
      });
    });

    cy.contains(".govuk-link", "Change phone number");
  });

  it("doesn't allow me to edit my phone number without permission", () => {
    cy.setupPermissions({});

    cy.visit("/my-details");

    cy.contains(".govuk-link", "Change phone number").should("not.exist");
  });
});
