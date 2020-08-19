<template>
  <a href="#main-content" class="govuk-skip-link">Skip to main content</a>
  
  <div class="moj-primary-navigation">
    <div class="moj-primary-navigation__container">
      <div class="moj-primary-navigation__nav">
        <nav class="moj-primary-navigation" aria-label="Primary navigation">
          <ul class="moj-primary-navigation__list">
            <li class="moj-primary-navigation__item">
              <router-link class="moj-primary-navigation__link" to="/home">Users</router-link>
            </li>
            
            <li class="moj-primary-navigation__item">
              <router-link class="moj-primary-navigation__link" to="/teams">Teams</router-link>
            </li>
            
            <li class="moj-primary-navigation__item">
              <router-link class="moj-primary-navigation__link" to="/my-details">My details</router-link>
            </li>
            
            <li class="moj-primary-navigation__item">
              <router-link class="moj-primary-navigation__link" to="/change-password">Change password</router-link>
            </li>
          </ul>
        </nav>
      </div>
    </div>
  </div>

  <div class="govuk-width-container ">
    <main class="govuk-main-wrapper govuk-main-wrapper--auto-spacing" id="main-content" role="main">
      <router-view />
    </main>
  </div>
</template>
