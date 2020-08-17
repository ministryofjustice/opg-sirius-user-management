import React from "react";

const AppNavigation = ({ items }) => (
  <div class="moj-primary-navigation">
    <div class="moj-primary-navigation__container">
      <div class="moj-primary-navigation__nav">
        <nav class="moj-primary-navigation" aria-label="Primary navigation">
          <ul class="moj-primary-navigation__list">
            {items.map((item) => (
              <li class="moj-primary-navigation__item">
                <a class="moj-primary-navigation__link" aria-current={item.active && 'page'} href={item.href}>
                  {item.text}
                </a>
              </li>
            ))}
          </ul>
        </nav>
      </div>
    </div>
  </div>
);

export default AppNavigation;
