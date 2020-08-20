import { changePassword } from "@/services/change-password";
import ChangePassword from "@/views/ChangePassword.vue";
import { shallowMount } from "@vue/test-utils";
import flushPromises from "flush-promises";

jest.mock("@/services/change-password");

let changePasswordMock: jest.Mock = changePassword as any;

beforeEach(() => {
  changePasswordMock.mockClear();
});

describe("ChangePassword.vue", () => {
  describe("when sending a valid form", () => {
    beforeEach(() => {
      changePasswordMock.mockResolvedValue({
        ok: true
      });
    });

    it("calls Sirius with the password request", async () => {
      const $router = { push: jest.fn() };

      const wrapper = shallowMount(ChangePassword, {});
      (wrapper.vm as any).$router = $router;

      wrapper.find("#f-currentpassword").setValue("this");
      wrapper.find("#f-password1").setValue("that");
      wrapper.find("#f-password2").setValue("ok");

      await wrapper.find("form").trigger("submit.prevent");
      await flushPromises();

      expect($router.push).toBeCalledWith("/my-details");
    });
  });

  describe("when sending an invalid form", () => {
    beforeEach(() => {
      changePasswordMock.mockResolvedValue({
        ok: false,
        error: "problems"
      });
    });

    it("calls Sirius with the password request", async () => {
      const wrapper = shallowMount(ChangePassword, {});

      wrapper.find("#f-currentpassword").setValue("this");

      await wrapper.find("form").trigger("submit.prevent");
      await flushPromises();

      expect(wrapper.find(".govuk-error-summary__list li").text()).toEqual(
        "problems"
      );
    });
  });
});
