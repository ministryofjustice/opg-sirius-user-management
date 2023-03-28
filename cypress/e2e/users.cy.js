describe("Users", () => {
  beforeEach(() => {
    cy.visit("/users");
  });

  it("allows me to search for admin user", () => {
    const expected = [
      "system admin",
      "",
      "system.admin@opgtest.com",
      "Active",
      "Edit",
    ];
    search("admin", expected);
  });

  it("allows me to search for non-admin user", () => {
    const expected = [
      "CardPayment User",
      "Card Payment Team",
      "card.payment.user@opgtest.com",
      "Active",
      "Edit",
    ];
    search("CardPayment", expected);
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
