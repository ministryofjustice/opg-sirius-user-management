import React from "react";

import AppHeader from "./AppHeader";
import AppNavigation from "./AppNavigation";
import ChangePassword from "../pages/ChangePassword";
import { Footer } from "govuk-react-jsx";

const App = () => (
  <>
    <a href="#main-content" className="govuk-skip-link">
      Skip to main content
    </a>
    <AppHeader></AppHeader>
    <AppNavigation></AppNavigation>
    <div class="govuk-width-container ">
      <main
        class="govuk-main-wrapper govuk-main-wrapper--auto-spacing"
        id="main-content"
        role="main"
      >
        <ChangePassword />
      </main>
    </div>
    <Footer></Footer>
  </>
);

export default App
