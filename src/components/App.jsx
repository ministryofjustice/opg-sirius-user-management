import React from "react";

import Header from "./moj/Header";
import PrimaryNavigation from "./moj/PrimaryNavigation";
import ChangePassword from "../pages/ChangePassword";
import { Footer } from "govuk-react-jsx";

const App = () => (
  <>
    <a href="#main-content" className="govuk-skip-link">
      Skip to main content
    </a>
    <Header
      serviceLabel={{ text: "Sirius User Management", href: "/" }}
      items={[
        {
          text: "Supervision",
          href: "#",
        },
        {
          text: "LPA",
          href: "#",
        },
        {
          text: "Logout",
          href: "#",
        },
      ]}
    ></Header>
    <PrimaryNavigation
      items={[
        {
          text: "Users",
          href: "/users",
        },
        {
          text: "Teams",
          href: "/teams",
        },
        {
          text: "My details",
          href: "/my-details",
        },
        {
          text: "Change password",
          href: "/change-password",
          active: true,
        },
      ]}
    ></PrimaryNavigation>
    <div className="govuk-width-container ">
      <main
        className="govuk-main-wrapper govuk-main-wrapper--auto-spacing"
        id="main-content"
        role="main"
      >
        <ChangePassword />
      </main>
    </div>
    <Footer></Footer>
  </>
);

export default App;
