import React from "react";

const AppNavigation = () => (
  <div class="moj-primary-navigation">
    <div class="moj-primary-navigation__container">
      <div class="moj-primary-navigation__nav">
        <nav class="moj-primary-navigation" aria-label="Primary navigation">
          <ul class="moj-primary-navigation__list">
            <li class="moj-primary-navigation__item">
              <a class="moj-primary-navigation__link" href="/home">
                Users
              </a>
            </li>

            <li class="moj-primary-navigation__item">
              <a class="moj-primary-navigation__link" href="/teams">
                Teams
              </a>
            </li>

            <li class="moj-primary-navigation__item">
              <a class="moj-primary-navigation__link" href="/my-details">
                My details
              </a>
            </li>

            <li class="moj-primary-navigation__item">
              <a class="moj-primary-navigation__link" href="/change-password">
                Change password
              </a>
            </li>
          </ul>
        </nav>
      </div>
    </div>
  </div>
);

export default AppNavigation;
