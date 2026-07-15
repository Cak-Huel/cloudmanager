<template>
  <header>
    <img v-if="showLogo" :src="logoURL" />
    <Action
      v-if="showMenu"
      class="menu-button"
      icon="menu"
      :label="t('buttons.toggleSidebar')"
      @action="layoutStore.showHover('sidebar')"
    />

    <slot />

    <div
      id="dropdown"
      :class="{ active: layoutStore.currentPromptName === 'more' }"
    >
      <slot name="actions" />
    </div>

    <Action
      :icon="currentTheme === 'dark' ? 'light_mode' : 'dark_mode'"
      :label="currentTheme === 'dark' ? 'Light Mode' : 'Dark Mode'"
      @action="toggleThemeAction"
    />

    <Action
      v-if="ifActionsSlot"
      id="more"
      icon="more_vert"
      :label="t('buttons.more')"
      @action="layoutStore.showHover('more')"
    />

    <div
      class="overlay"
      v-show="layoutStore.currentPromptName == 'more'"
      @click="layoutStore.closeHovers"
    />
  </header>
</template>

<script setup lang="ts">
import { useLayoutStore } from "@/stores/layout";

import { logoURL } from "@/utils/constants";

import Action from "@/components/header/Action.vue";
import { computed, useSlots, ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { getTheme, toggleTheme } from "@/utils/theme";

defineProps<{
  showLogo?: boolean;
  showMenu?: boolean;
}>();

const layoutStore = useLayoutStore();
const slots = useSlots();

const { t } = useI18n();

const ifActionsSlot = computed(() => (slots.actions ? true : false));

const currentTheme = ref(getTheme());

const toggleThemeAction = () => {
  toggleTheme();
  currentTheme.value = getTheme();
};

onMounted(() => {
  currentTheme.value = getTheme();
});
</script>

<style></style>
