const { pactWith } = require("jest-pact");
const { Matchers } = require("@pact-foundation/pact");
const { default: changePassword } = require("./changePassword");
const { string } = Matchers;

pactWith(
  { consumer: "SiriusUserManagement", provider: "Sirius" },
  (provider) => {
    beforeEach(() => {
      globalThis.baseUrl = provider.mockService.baseUrl;
    });

    describe("Change Password", () => {
      describe("Valid password change", () => {
        beforeEach(() =>
          provider.addInteraction({
            state: "Current password is OldPass",
            uponReceiving: "a new password",
            withRequest: {
              method: "POST",
              path: "/auth/change-password",
              headers: {
                "Content-type":
                  "application/x-www-form-urlencoded;charset=UTF-8",
              },
            },
            willRespondWith: {
              status: 200,
              headers: { "Content-Type": "application/json" },
              body: {},
            },
          })
        );

        it("accepts a valid password change", async () => {
          const { status, error } = await changePassword(
            "OldPass",
            "NewPass",
            "NewPass"
          );
          expect(status).toEqual(200);
          expect(error).toEqual(null);
        });
      });

      describe("Wrong password provided", () => {
        beforeEach(() =>
          provider.addInteraction({
            state: "Current password is OldPass",
            uponReceiving: "an invalid existing password",
            withRequest: {
              method: "POST",
              path: "/auth/change-password",
              headers: {
                "Content-type":
                  "application/x-www-form-urlencoded;charset=UTF-8",
              },
              body:
                "existingPassword=BadPass&password=NewPass&confirmPassword=NewPass",
            },
            willRespondWith: {
              status: 400,
              headers: { "Content-Type": "application/json" },
              body: {
                errors: "Password supplied was incorrect or user is not active",
              },
            },
          })
        );

        it("accepts a valid password change", async () => {
          const { status, error } = await changePassword(
            "BadPass",
            "NewPass",
            "NewPass"
          );
          expect(status).toEqual(400);
          expect(error).toEqual(
            "Password supplied was incorrect or user is not active"
          );
        });
      });
    });
  }
);
