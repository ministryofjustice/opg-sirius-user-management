describe("Edit lay percentage", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/random-reviews/edit/lay-percentage");
    });

    it("allows me to change the lay percentage value and throws an error after inputting the incorrect value", () => {
        cy.get("#f-layPercentage").clear().type("200");
        cy.get("button[type=submit]").click();
        cy.get('#name-error').contains("Enter a percentage between 0 and 100 for lay cases")
    });
});
