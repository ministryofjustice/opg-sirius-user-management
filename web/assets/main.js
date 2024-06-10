import * as GOVUKFrontend from "govuk-frontend";
import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import GoBack from "./go-back";
import CloseTab from "./close-tab";

const goBack = document.querySelectorAll('[data-module="moj-go-back"]');
goBack.forEach(function (goBack) {
  new GoBack(goBack);
});

const closeTab = document.querySelectorAll('[data-module="moj-close-tab"]');
closeTab.forEach(function (closeTab) {
  new CloseTab(closeTab);
});

GOVUKFrontend.initAll();
