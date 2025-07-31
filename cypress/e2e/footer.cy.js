describe("Footer", () => {
    beforeEach(() => {
        cy.setupPermissions({ "v1-users-updatetelephonenumber": ["put"] });

        cy.addMock("/api/v1/users/current", "GET", {
            status: 200,
            body: {
                firstname: "system",
                surname: "admin",
                email: "system.admin@opgtest.com",
                phoneNumber: "03004560300",
                roles: ["OPG User", "Finance", "System Admin"],
                teams: [
                    {
                        displayName: "Administrative Team",
                    },
                ],
            },
        });

        cy.visit("/my-details");
    });

    it("should show the accessibility link", () => {
        cy.get('[data-cy="accessibilityStatement"]').should("contain", "Accessibility statement");
    });
});
