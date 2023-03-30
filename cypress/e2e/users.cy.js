describe("Users", () => {
  beforeEach(() => {
    cy.visit("/users");
  });

  it("allows me to search for not in a team", () => {
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
    const expected = [
      "Anton Mccoy",
      "my friendly team",
      "anton.mccoy@opgtest.com",
      "Active",
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
