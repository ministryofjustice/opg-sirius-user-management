import { changePassword } from "@/services/change-password";
import flushPromises from "flush-promises";
import { pactWith } from "jest-pact";

pactWith({ consumer: "SiriusUserManagement", provider: "Sirius" }, provider => {
  beforeEach(() => {
    process.env.VUE_APP_SIRIUS_URL = provider.mockService.baseUrl;
  });

  describe("changePassword", () => {
    describe("when sending a valid form", () => {
      beforeEach(() => {
        provider.addInteraction({
          state: "Current password is this",
          uponReceiving: "a new password",
          withRequest: {
            method: "POST",
            path: "/auth/change-password",
            headers: {
              "Content-type": "application/x-www-form-urlencoded;charset=UTF-8"
            },
            body: "existingPassword=this&password=that&confirmPassword=ok"
          },
          willRespondWith: {
            status: 200,
            headers: { "Content-Type": "application/json" },
            body: {}
          }
        });
      });

      it("calls Sirius with the password request", async () => {
        const { ok } = await changePassword("this", "that", "ok");
        await flushPromises();

        expect(ok).toBe(true);
      });
    });

    describe("when sending an invalid form", () => {
      beforeEach(() => {
        provider.addInteraction({
          state: "Current password is this",
          uponReceiving: "a bad request",
          withRequest: {
            method: "POST",
            path: "/auth/change-password",
            headers: {
              "Content-type": "application/x-www-form-urlencoded;charset=UTF-8"
            },
            body: "existingPassword=this&password=&confirmPassword="
          },
          willRespondWith: {
            status: 400,
            headers: { "Content-Type": "application/json" },
            body: {
              errors: "you what"
            }
          }
        });
      });

      it("calls Sirius with the password request", async () => {
        const { ok, error } = await changePassword("this", "", "");
        await flushPromises();

        expect(ok).toBe(false);
        expect(error).toBe("you what");
      });
    });
  });
});
