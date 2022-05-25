import * as GOVUKFrontend from "govuk-frontend";
import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";

document.body.className = document.body.className
  ? document.body.className + " js-enabled"
  : "js-enabled";
GOVUKFrontend.initAll();
