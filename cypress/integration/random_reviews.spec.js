describe("Random Reviews", () => {
    beforeEach(() => {
        cy.visit("/random-reviews", {
            headers: {
                Cookie: "XSRF-TOKEN=abcde; Other=other",
                "OPG-Bypass-Membrane": "1",
                "X-XSRF-TOKEN": "abcde",
            },
        });
    });

    it("shows all random reviews", () => {
        const expected = [
            ["Lay", 20],
            ["Review cycle", 3]
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
    });

    it("allows me to change the lay percentage", () => {
        cy.contains("#layPercentageChange", "Change");
    });

    it("allows me to change the review cycle", () => {
        cy.contains("#layReviewCycleChange", "Change");
    });
});
