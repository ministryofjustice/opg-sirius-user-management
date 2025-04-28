import { addMock, reset } from "../mocks/wiremock";

Cypress.Commands.add("addMock", async (url, method, response) => {
  await addMock(url, method, response);
});

Cypress.Commands.add("resetMocks", async () => {
  await reset();
});

/**
 * @param {{[key: string]: string[]}} permissions
 */
Cypress.Commands.add("setupPermissions", async (permissions = {}) => {
  await addMock("/supervision-api/v1/permissions", "GET", {
    status: 200,
    body: Object.entries(permissions).reduce(
      (set, [endpoint, methods]) => ({
        ...set,
        [endpoint]: {
          permissions: methods,
        },
      }),
      {}
    ),
  });

  await addMock("/supervision-api/v1/users/current", "GET", {
    status: 200,
    body: {},
  });
});
