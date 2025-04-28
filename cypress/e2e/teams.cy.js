describe("Teams", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-teams": ["put"] });

    cy.addMock("/supervision-api/v1/teams", "GET", {
      status: 200,
      body: [
        {
          id: 748,
          displayName: "Finance Team",
          teamType: {
            handle: "FINANCE",
            label: "Finance",
          },
          members: [
            {
              displayName: "John Ruecker",
            },
          ],
        },
        {
          id: 15,
          displayName: "File Creation Team",
          members: [
            {
              displayName: "Arvel Buckridge",
            },
            {
              displayName: "Thea Wyman",
            },
          ],
        },
      ],
    });

    cy.visit("/teams");
  });

  it("lists all teams", () => {
    cy.get(".govuk-table__row").should("have.length", 3);

    const teams = [
      ["Finance Team", "Supervision — Finance", "1"],
      ["File Creation Team", "LPA", "2"],
    ];

    teams.forEach((team, teamIndex) => {
      team.forEach((cell, cellIndex) => {
        cy.get(
          `.govuk-table__body > .govuk-table__row:nth-child(${teamIndex + 1}) > *:nth-child(${cellIndex + 1})`
        ).should("contain", cell);
      });
    });
  });

  it("allows me to search for a team", () => {
    cy.get("#f-search").clear().type("Finance");
    cy.get("button[type=submit]").click();

    cy.get(".govuk-table__body > .govuk-table__row").should("have.length", 1);

    cy.get("#f-search").clear().type("no such team");
    cy.get("button[type=submit]").click();

    cy.get(".govuk-table__body > .govuk-table__row").should("have.length", 0);
  });

  it("allows me to add a new team", () => {
    cy.contains(".govuk-button", "Add new team");
  });
});
