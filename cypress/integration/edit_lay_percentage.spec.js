describe("Edit lay percentage", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/random-reviews/edit/lay-percentage");
    });

    it("shows the lay percentage value", () => {
        cy.get("#f-layPercentage").should("have.value", "20");
    });

    it("allows me to change the lay percentage value", () => {
        cy.get("#f-layPercentage").clear().type("30");
        cy.get("button[type=submit]").click();

    });
});
