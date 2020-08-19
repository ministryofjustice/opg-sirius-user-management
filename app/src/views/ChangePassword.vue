<template>
  <div class="govuk-grid-row">
    <div class="govuk-grid-column-two-thirds">
      <h1 class="govuk-heading-xl">Change password</h1>

      <div v-if="error" class="govuk-error-summary" aria-labelledby="error-summary-title" role="alert" tabindex="-1" data-module="govuk-error-summary">
        <h2 class="govuk-error-summary__title" id="error-summary-title">There was a problem</h2>

        <div class="govuk-error-summary__body">
          <ul class="govuk-list govuk-error-summary__list">
            <li>{{ error }}</li>
          </ul>
        </div>
      </div>
      
      <form class="form" @submit.prevent="handleSubmit">
        <div class="govuk-form-group">
          <label class="govuk-label govuk-label--m" for="f-currentpassword">Current password</label>
          <input class="govuk-input" id="f-currentpassword" type="password" v-model="currentPassword" :disabled="loading">
        </div>

        <fieldset class="govuk-fieldset">
          <legend class="govuk-fieldset__legend govuk-fieldset__legend--m">New password</legend>

          <div class="govuk-form-group">
            <label class="govuk-label" for="f-password1">Create your new password</label>
            <input class="govuk-input" id="f-password1" type="password" v-model="newPassword" :disabled="loading">
          </div>

          <div class="govuk-form-group">
            <label class="govuk-label" for="f-password2">Confirm new password</label>
            <input class="govuk-input" id="f-password2" type="password" v-model="newPasswordConfirm" :disabled="loading">
          </div>
        </fieldset>

        <button class="govuk-button" data-module="govuk-button" :disabled="loading">Save changes</button>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
  import { defineComponent } from 'vue';

  export default defineComponent({
    name: 'ChangePassword',
    data() {
      return {
        currentPassword: '',
        newPassword: '',
        newPasswordConfirm: '',
        loading: false,
        error: '',
      };
    },
    methods: {
      async handleSubmit() {
        this.error = '';
        this.loading = true;
        
        const data = new URLSearchParams();
        data.append('existingPassword', this.currentPassword);
        data.append('password', this.newPassword);
        data.append('confirmPassword', this.newPasswordConfirm);

        try {
          const resp = await fetch(`${process.env.VUE_APP_SIRIUS_URL}/auth/change-password`, {
            mode: 'cors',
            method: 'POST',
            headers: new Headers({'Content-Type': 'application/x-www-form-urlencoded'}),
            body: data.toString(),
          });
          this.loading = false;
          
          if (resp.ok) {
            (this as any).$router.push('/my-details');
          } else {
            this.error = await resp.json().then(body => body.errors);
          }
        } catch (err) {
          this.loading = false;
          this.error = 'something unexpected happened, try again?';
        }
      },
    },
  });
</script>
