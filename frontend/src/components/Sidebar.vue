<template>
  <div v-show="active" @click="closeHovers" class="overlay"></div>
  <nav :class="{ active }">
    <template v-if="isLoggedIn">
      <div v-if="user.perm.create" class="new-menu-wrapper">
        <button
          @click="toggleNewMenu"
          class="action new-btn"
          id="new-button"
          :aria-label="$t('buttons.new')"
          :title="$t('buttons.new')"
        >
          <i class="material-icons">add</i>
          <span>{{ $t("buttons.new") }}</span>
        </button>

        <div v-if="newMenuOpen" class="new-menu-dropdown">
          <button @click="triggerNewDir" class="dropdown-item">
            <i class="material-icons">create_new_folder</i>
            <span>{{ $t("sidebar.newFolder") }}</span>
          </button>
          <button @click="triggerNewFile" class="dropdown-item">
            <i class="material-icons">note_add</i>
            <span>{{ $t("sidebar.newFile") }}</span>
          </button>
          <button @click="triggerUploadFile" class="dropdown-item">
            <i class="material-icons">upload_file</i>
            <span>{{ $t("buttons.upload") }} {{ $t("buttons.file") }}</span>
          </button>
          <button @click="triggerUploadFolder" class="dropdown-item">
            <i class="material-icons">drive_folder_upload</i>
            <span>{{ $t("buttons.upload") }} {{ $t("buttons.folder") }}</span>
          </button>
        </div>
      </div>

      <button
        @click="toAccountSettings"
        class="action"
        :class="{ active: $route.path.startsWith('/settings') && !$route.path.startsWith('/settings/global') }"
        :aria-selected="($route.path.startsWith('/settings') && !$route.path.startsWith('/settings/global')) ? 'true' : 'false'"
      >
        <i class="material-icons">account_circle</i>
        <span>{{ user.username }}</span>
      </button>
      <button
        class="action"
        :class="{ active: isFiles }"
        :aria-selected="isFiles ? 'true' : 'false'"
        @click="toRoot"
        :aria-label="$t('sidebar.myFiles')"
        :title="$t('sidebar.myFiles')"
      >
        <i class="material-icons">folder</i>
        <span>{{ $t("sidebar.myFiles") }}</span>
      </button>

      <div v-if="user.perm.admin">
        <button
          class="action"
          :class="{ active: $route.path.startsWith('/settings/global') }"
          :aria-selected="$route.path.startsWith('/settings/global') ? 'true' : 'false'"
          @click="toGlobalSettings"
          :aria-label="$t('sidebar.settings')"
          :title="$t('sidebar.settings')"
        >
          <i class="material-icons">settings</i>
          <span>{{ $t("sidebar.settings") }}</span>
        </button>
      </div>
      <button
        v-if="canLogout"
        @click="logout"
        class="action"
        id="logout"
        :aria-label="$t('sidebar.logout')"
        :title="$t('sidebar.logout')"
      >
        <i class="material-icons">logout</i>
        <span>{{ $t("sidebar.logout") }}</span>
      </button>
    </template>
    <template v-else>
      <router-link
        v-if="!hideLoginButton"
        class="action"
        to="/login"
        :aria-label="$t('sidebar.login')"
        :title="$t('sidebar.login')"
      >
        <i class="material-icons">login</i>
        <span>{{ $t("sidebar.login") }}</span>
      </router-link>

      <router-link
        v-if="signup"
        class="action"
        to="/login"
        :aria-label="$t('sidebar.signup')"
        :title="$t('sidebar.signup')"
      >
        <i class="material-icons">person_add</i>
        <span>{{ $t("sidebar.signup") }}</span>
      </router-link>
    </template>

    <div
      class="credits"
      v-if="isFiles && !disableUsedPercentage"
    >
      <progress-bar :val="usage.usedPercentage" size="small"></progress-bar>
      <br />
      {{ $t("sidebar.diskUsed", { used: usage.used, total: usage.total }) }}
    </div>

    <p class="credits">
      <span>
        <span v-if="disableExternal">File Browser</span>
        <a
          v-else
          rel="noopener noreferrer"
          target="_blank"
          href="https://github.com/filebrowser/filebrowser"
          >File Browser</a
        >
        <span> {{ " " }} {{ version }}</span>
      </span>
      <span>
        <a @click="help">{{ $t("sidebar.help") }}</a>
      </span>
    </p>
  </nav>
</template>

<script>
import { reactive, ref } from "vue";
import { mapActions, mapState } from "pinia";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import * as auth from "@/utils/auth";
import * as upload from "@/utils/upload";
import {
  version,
  signup,
  hideLoginButton,
  disableExternal,
  disableUsedPercentage,
  noAuth,
  logoutPage,
  loginPage,
} from "@/utils/constants";
import { files as api } from "@/api";
import ProgressBar from "@/components/ProgressBar.vue";
import prettyBytes from "pretty-bytes";

const USAGE_DEFAULT = { used: "0 B", total: "0 B", usedPercentage: 0 };

export default {
  name: "sidebar",
  setup() {
    const usage = reactive(USAGE_DEFAULT);
    const newMenuOpen = ref(false);
    return { usage, usageAbortController: new AbortController(), newMenuOpen };
  },
  components: {
    ProgressBar,
  },
  inject: ["$showError"],
  computed: {
    ...mapState(useAuthStore, ["user", "isLoggedIn"]),
    ...mapState(useFileStore, ["isFiles", "reload"]),
    ...mapState(useLayoutStore, ["currentPromptName"]),
    active() {
      return this.currentPromptName === "sidebar";
    },
    signup: () => signup,
    hideLoginButton: () => hideLoginButton,
    version: () => version,
    disableExternal: () => disableExternal,
    disableUsedPercentage: () => disableUsedPercentage,
    canLogout: () => !noAuth && (loginPage || logoutPage !== "/login"),
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers", "showHover"]),
    abortOngoingFetchUsage() {
      this.usageAbortController.abort();
    },
    async fetchUsage() {
      const path = this.$route.path.endsWith("/")
        ? this.$route.path
        : this.$route.path + "/";
      let usageStats = USAGE_DEFAULT;
      if (this.disableUsedPercentage) {
        return Object.assign(this.usage, usageStats);
      }
      try {
        this.abortOngoingFetchUsage();
        this.usageAbortController = new AbortController();
        const usage = await api.usage(path, this.usageAbortController.signal);
        usageStats = {
          used: prettyBytes(usage.used, { binary: true }),
          total: prettyBytes(usage.total, { binary: true }),
          usedPercentage: Math.round((usage.used / usage.total) * 100),
        };
      } finally {
        return Object.assign(this.usage, usageStats);
      }
    },
    toRoot() {
      this.$router.push({ path: "/files" });
      this.closeHovers();
    },
    toAccountSettings() {
      this.$router.push({ path: "/settings/profile" });
      this.closeHovers();
    },
    toGlobalSettings() {
      this.$router.push({ path: "/settings/global" });
      this.closeHovers();
    },
    help() {
      this.showHover("help");
    },
    logout: auth.logout,
    toggleNewMenu(event) {
      event.stopPropagation();
      this.newMenuOpen = !this.newMenuOpen;
    },
    closeNewMenu() {
      this.newMenuOpen = false;
    },
    triggerNewDir() {
      this.closeNewMenu();
      this.showHover('newDir');
    },
    triggerNewFile() {
      this.closeNewMenu();
      this.showHover('newFile');
    },
    openUpload(isFolder) {
      const input = document.createElement("input");
      input.type = "file";
      input.multiple = true;
      input.webkitdirectory = isFolder;
      input.onchange = this.uploadInput;
      input.click();
    },
    async uploadInput(event) {
      const files = event.currentTarget.files;
      if (!files || files.length === 0) return;

      const folder_upload = !!files[0].webkitRelativePath;
      const uploadFiles = [];
      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        const fullPath = folder_upload ? file.webkitRelativePath : undefined;
        uploadFiles.push({
          file,
          name: file.name,
          size: file.size,
          isDir: false,
          fullPath,
        });
      }

      const path = this.$route.path.endsWith("/")
        ? this.$route.path
        : this.$route.path + "/";

      const uploadPath = path.startsWith("/files") ? path : "/files/";

      const conflict = await upload.checkConflict(uploadFiles, uploadPath);

      if (conflict.length > 0) {
        this.showHover({
          prompt: "resolve-conflict",
          props: {
            conflict: conflict,
            isUploadAction: true,
          },
          confirm: (ev, result) => {
            ev.preventDefault();
            this.closeHovers();
            for (let i = result.length - 1; i >= 0; i--) {
              const item = result[i];
              if (item.checked.length == 2) {
                continue;
              } else if (item.checked.length == 1 && item.checked[0] == "origin") {
                uploadFiles[item.index].overwrite = true;
              } else {
                uploadFiles.splice(item.index, 1);
              }
            }
            if (uploadFiles.length > 0) {
              upload.handleFiles(uploadFiles, uploadPath);
            }
          },
        });
        return;
      }

      upload.handleFiles(uploadFiles, uploadPath);
    },
    triggerUploadFile() {
      this.closeNewMenu();
      this.openUpload(false);
    },
    triggerUploadFolder() {
      this.closeNewMenu();
      this.openUpload(true);
    },
  },
  watch: {
    $route: {
      handler(to) {
        if (to.path.includes("/files")) {
          this.fetchUsage();
        }
      },
      immediate: true,
    },
  },
  mounted() {
    document.addEventListener("click", this.closeNewMenu);
  },
  unmounted() {
    this.abortOngoingFetchUsage();
    document.removeEventListener("click", this.closeNewMenu);
  },
};
</script>
