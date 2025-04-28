describe("Edit my details", () => {
  beforeEach(() => {
    cy.setupPermissions({ "v1-users-updatetelephonenumber": ["put"] });

    cy.addMock("/supervision-api/v1/users/current", "GET", {
      status: 200,
      body: {
        id: 949,
        phoneNumber: "03004560300",
      },
    });

    cy.visit("/my-details/edit");
  });

  it("shows my phone number", () => {
    cy.get("#f-phonenumber").should("have.value", "03004560300");
  });

  it("allows me to change my phone number", () => {
    cy.get("#f-phonenumber").clear().type("123456789");

    cy.addMock("/supervision-api/v1/users/949/updateTelephoneNumber", "PUT", {
      status: 200,
    });

    cy.get("button[type=submit]").click();

    cy.contains(".moj-alert", "You have successfully edited your details.");
  });
});
