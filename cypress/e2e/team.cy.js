describe("Team", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-teams": ["put"] });

    cy.addMock("/supervision-api/v1/teams/14", "GET", {
      status: 200,
      body: {
        id: 748,
        displayName: "Finance Team",
        members: [
          {
            displayName: "John Ruecker",
            email: "j.ruecker1@opg.example",
          },
        ],
      },
    });

    cy.visit("/teams/14");
  });

  it("shows team name", () => {
    cy.contains("h1", "Finance Team");
  });

  it("shows team members", () => {
    cy.get(".govuk-table__row").should("have.length", 2);

    const expected = ["Select", "John", "j.ruecker1@opg.example"];

    cy.get(".govuk-table__body > .govuk-table__row")
      .children()
      .each(($el, index) => {
        cy.wrap($el).should("contain", expected[index]);
      });
  });

  it("allows me to edit the team", () => {
    cy.contains(".govuk-button", "Edit team");
  });

  it("allows me to add a team member", () => {
    cy.contains(".govuk-button", "Add user to team");
  });

  it("allows me to remove team members", () => {
    cy.contains(".govuk-button", "Remove selected from team");

    cy.get(
      ".govuk-table__body > .govuk-table__row input[type=checkbox]"
    ).should("have.length", 1);
  });
});
