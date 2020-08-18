import React from "react";

const PrimaryNavigation = ({ items }) => (
  <div className="moj-primary-navigation">
    <div className="moj-primary-navigation__container">
      <div className="moj-primary-navigation__nav">
        <nav className="moj-primary-navigation" aria-label="Primary navigation">
          <ul className="moj-primary-navigation__list">
            {items.map((item, index) => (
              <li key={index} className="moj-primary-navigation__item">
                <a
                  className="moj-primary-navigation__link"
                  aria-current={item.active && "page"}
                  href={item.href}
                >
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

export default PrimaryNavigation;
