describe("Edit lay review cycle", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/random-reviews/edit/lay-review-cycle");
    });

    it("allows me to change the lay review cycle value", () => {
        cy.get("#f-layReviewCycle").clear().type("15");
        cy.get("button[type=submit]").click();
        cy.get('#name-error').contains("Enter a review cycle between 1 and 10 for lay cases")
    });
});
