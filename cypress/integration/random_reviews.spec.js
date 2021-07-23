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
                cy.get(".hook-layPercentageKey").should(
                    "have.text",
                    expected[index][0]
                );
                cy.get(".hook-layPercentageValue").should(
                    "have.text",
                    expected[index][1]
                );
            });
        });
    });

    it("allows me to change the lay percentage", () => {
        cy.contains("#hook-layPercentageChange", "Change");
    });

    it("allows me to change the review cycle", () => {
        cy.contains("#hook-layReviewCycleChange", "Change");
    });
});
