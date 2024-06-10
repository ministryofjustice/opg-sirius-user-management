export default class CloseTab {
  constructor(element) {
    this.closeLinks = element.querySelectorAll(".close-tab");

    this.closeLinks.forEach((element) => {
       this._closeTab = this._closeTab.bind(this);
       element.addEventListener("click", this._closeTab);
    });
  }

    _closeTab(event) {
        window.location.href = window.close();
    }
};
