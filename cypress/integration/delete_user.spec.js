describe("Delete user", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/delete-user/123");
    });

    it("allows me to delete a user", () => {
        cy.get("button[type=submit]").click();
        cy.url().should("include", "/delete-user/123")
        cy.get('a[href*="/users"]').contains('Continue').click()
        cy.url().should("include", "/users");
    });
});
