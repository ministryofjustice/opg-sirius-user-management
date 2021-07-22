describe("Random Reviews", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/random-reviews");
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
