import * as GOVUKFrontend from "govuk-frontend";
import * as MOJFrontend from "@ministryofjustice/frontend";
import CloseTab from "./close-tab";

const closeTab = document.querySelectorAll('[data-module="moj-close-tab"]');
closeTab.forEach(function (closeTab) {
  new CloseTab(closeTab);
});

GOVUKFrontend.initAll();
