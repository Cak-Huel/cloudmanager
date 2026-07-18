<template>
  <div id="login" :class="{ recaptcha: recaptcha }">
    <form @submit="submit">
      <img :src="logoURL" alt="Cloud Manager" />
      <h1>{{ name }}</h1>
      <p v-if="reason != null" class="logout-message">
        {{ t(`login.logout_reasons.${reason}`) }}
      </p>
      <div v-if="error !== ''" class="wrong">{{ error }}</div>

      <md-outlined-text-field
        autofocus
        style="width: 100%; margin-bottom: 1em;"
        type="text"
        autocapitalize="off"
        :value="username"
        @input="username = ($event.target as HTMLInputElement).value"
        :label="t('login.username')"
      ></md-outlined-text-field>

      <md-outlined-text-field
        style="width: 100%; margin-bottom: 1em;"
        :type="showPassword ? 'text' : 'password'"
        :value="password"
        @input="password = ($event.target as HTMLInputElement).value"
        :label="t('login.password')"
      >
        <md-icon-button
          slot="trailing-icon"
          type="button"
          @click.prevent="showPassword = !showPassword"
          style="margin-right: 4px;"
        >
          <i class="material-icons">{{ showPassword ? 'visibility_off' : 'visibility' }}</i>
        </md-icon-button>
      </md-outlined-text-field>

      <md-outlined-text-field
        style="width: 100%; margin-bottom: 1.5em;"
        v-if="createMode"
        :type="showConfirmPassword ? 'text' : 'password'"
        :value="passwordConfirm"
        @input="passwordConfirm = ($event.target as HTMLInputElement).value"
        :label="t('login.passwordConfirm')"
      >
        <md-icon-button
          slot="trailing-icon"
          type="button"
          @click.prevent="showConfirmPassword = !showConfirmPassword"
          style="margin-right: 4px;"
        >
          <i class="material-icons">{{ showConfirmPassword ? 'visibility_off' : 'visibility' }}</i>
        </md-icon-button>
      </md-outlined-text-field>

      <div v-if="recaptcha" id="recaptcha"></div>

      <md-filled-button
        style="width: 100%; margin-top: 0.5em; margin-bottom: 1em;"
        type="submit"
      >
        {{ createMode ? t('login.signup') : t('login.submit') }}
      </md-filled-button>

      <p @click="toggleMode" v-if="signup">
        {{ createMode ? t("login.loginInstead") : t("login.createAnAccount") }}
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { StatusError } from "@/api/utils";
import * as auth from "@/utils/auth";
import {
  name,
  logoURL,
  recaptcha,
  recaptchaKey,
  signup,
} from "@/utils/constants";
import { inject, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";

// Define refs
const createMode = ref<boolean>(false);
const error = ref<string>("");
const username = ref<string>("");
const password = ref<string>("");
const passwordConfirm = ref<string>("");
const showPassword = ref<boolean>(false);
const showConfirmPassword = ref<boolean>(false);

const route = useRoute();
const router = useRouter();
const { t } = useI18n({});
// Define functions
const toggleMode = () => (createMode.value = !createMode.value);

const $showError = inject<IToastError>("$showError")!;

const reason = route.query["logout-reason"] ?? null;

const submit = async (event: Event) => {
  event.preventDefault();
  event.stopPropagation();

  const redirect = (route.query.redirect || "/files/") as string;

  let captcha = "";
  if (recaptcha) {
    captcha = window.grecaptcha.getResponse();

    if (captcha === "") {
      error.value = t("login.wrongCredentials");
      return;
    }
  }

  if (createMode.value) {
    if (password.value !== passwordConfirm.value) {
      error.value = t("login.passwordsDontMatch");
      return;
    }
  }

  try {
    if (createMode.value) {
      await auth.signup(username.value, password.value);
    }

    await auth.login(username.value, password.value, captcha);
    router.push({ path: redirect });
  } catch (e: any) {
    // console.error(e);
    if (e instanceof StatusError) {
      if (e.status === 409) {
        error.value = t("login.usernameTaken");
      } else if (e.status === 403) {
        error.value = t("login.wrongCredentials");
      } else if (e.status === 400) {
        const match = e.message.match(/minimum length is (\d+)/);
        if (match) {
          error.value = t("login.passwordTooShort", { min: match[1] });
        } else {
          error.value = e.message;
        }
      } else {
        $showError(e);
      }
    }
  }
};

// Run hooks
onMounted(() => {
  if (!recaptcha) return;

  window.grecaptcha.ready(function () {
    window.grecaptcha.render("recaptcha", {
      sitekey: recaptchaKey,
    });
  });
});
</script>
