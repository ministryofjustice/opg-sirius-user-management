describe("Random Reviews", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-random-review-settings": ["get", "post"] });

    cy.addMock("/api/v1/random-review-settings", "GET", {
      status: 200,
      body: {
        layPercentage: 20,
        paPercentage: 30,
        proPercentage: 0,
        reviewCycle: 3,
      },
    });

    cy.visit("/random-reviews");
  });

  it("shows all random review settings", () => {
    cy.get(".hook-layPercentageRow").each(($el) => {
      cy.wrap($el).within(() => {
        cy.contains(".hook-layPercentageKey", "Lay");
        cy.contains(".hook-layPercentageValue", "20 %");
      });
    });

    cy.get(".hook-paPercentageRow").each(($el) => {
      cy.wrap($el).within(() => {
        cy.contains(".hook-paPercentageKey", "PA");
        cy.contains(".hook-paPercentageValue", "0 %");
      });
    });

    cy.contains("#hook-paPercentageChange", "Change");

    cy.get(".hook-reviewCycleRow").each(($el) => {
      cy.wrap($el).within(() => {
        cy.contains(".hook-reviewCycleKey", "Review cycle");
        cy.contains(".hook-reviewCycleValue", "3 year(s)");
      });
    });

    cy.get(".hook-proPercentageRow").each(($el) => {
      cy.wrap($el).within(() => {
        cy.contains(".hook-proPercentageKey", "Pro");
        cy.contains(".hook-proPercentageValue", "0 %");
      });
    });

    cy.contains("#hook-proPercentageChange", "Change");

    cy.contains("#hook-reviewCycleChange", "Change");
  });

  describe("Edit lay percentage", () => {
    it("throws an error after inputting the incorrect value", () => {
      cy.get("#hook-layPercentageChange").contains("Change").click();
      cy.get("#f-layPercentage").clear().type("200");

      cy.addMock("/api/v1/random-review-settings", "POST", {
        status: 400,
        body: {
          detail: "Enter a percentage between 0 and 100 for lay cases",
          status: 400,
        },
      });

      cy.get("button[type=submit]").click();
      cy.contains(
        "#name-error",
        "Enter a percentage between 0 and 100 for lay cases"
      );
    });
  });
});
