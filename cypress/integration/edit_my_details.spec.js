describe("Edit my details", () => {
    beforeEach(() => {
        Cypress.Cookies.debug(true);

        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/my-details/edit");

        cy.server();
    });

    it("shows my phone number", () => {
        cy.get("#f-phonenumber").should("have.value", "03004560300");
    });

    it("allows me to change my phone number", () => {
        cy.get("#f-phonenumber").clear().type("123456789");

        cy.route("POST", "/my-details/edit").as("formSuccess");
        cy.get("button[type=submit]").click();

        cy.contains(
            ".moj-banner",
            "SuccessYou have successfully edited your details."
        );
    });
});
