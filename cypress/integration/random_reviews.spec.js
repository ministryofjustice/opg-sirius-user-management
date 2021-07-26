describe("Random Reviews", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/random-reviews");
    });

    it("shows all random reviews", () => {
        const expected = [
            ["Lay", "20 %"],
            ["Review cycle", "3 year(s)"]
        ];

        cy.get(".hook-layPercentageRow").each(($el, index) => {
            cy.wrap($el).within(() => {
                cy.get(".hook-layPercentageKey").contains(expected[index][0]);
                cy.get(".hook-layPercentageValue").contains(expected[index][1]);
            });
        });
    });

    it("the lay percentage change option is present", () => {
        cy.contains("#hook-layPercentageChange", "Change");
    });

    it("the review cycle change option is present", () => {
        cy.contains("#hook-layReviewCycleChange", "Change");
    });
});
