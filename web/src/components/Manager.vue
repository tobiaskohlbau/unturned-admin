<template>
  <s-card width="600px" height="400px">
    <s-card-title>Manager</s-card-title>
    <s-card-content>
      <s-scroller>
        <i
          class="material-icons button"
          @click="openFolder('..')"
          :class="[path != '.' ? 'button_enabled' : 'button_disabled']"
          >arrow_back</i
        >
        <div class="line" v-for="folder in listing.folders" :key="folder">
          <div @click="openFolder(folder)">
            <i class="material-icons">folder</i>
            {{ folder }}
          </div>
        </div>
        <div class="line" v-for="file in listing.files" :key="file.name">
          <div @click="openFile(file)">
            <i class="material-icons">text_snippet</i>
            {{ file.name }}
          </div>
          <i class="material-icons button_enabled" @click="deleteFile(file)"
            >delete</i
          >
        </div>
      </s-scroller>
    </s-card-content>
    <s-card-actions>
      <s-button @click.stop="scheduleUpdate">Update</s-button>
      <s-button @click.stop="cancelUpdate" v-if="updateInProgress"
        >Cancle Update</s-button
      >
      <s-button @click.stop="backup">Backup</s-button>
    </s-card-actions>
    <s-popup :show="editorPopup">
      <s-card width="700px" height="600px">
        <s-card-title>Editor</s-card-title>
        <s-card-content>
          <editor :filepath="filepath" ref="editor"></editor>
        </s-card-content>
        <s-card-actions>
          <s-button @click.stop="editor.save()">Save</s-button>
          <s-button @click.stop="editorPopup = false">Close</s-button>
        </s-card-actions>
      </s-card>
    </s-popup>
    <s-popup :show="imagePopup">
      <s-card width="700px" height="600px">
        <s-card-title>Imageviewer</s-card-title>
        <s-card-content>
          <img :src="'/api/files?path=' + filepath" />
        </s-card-content>
        <s-card-actions>
          <s-button @click.stop="imagePopup = false">Close</s-button>
        </s-card-actions>
      </s-card>
    </s-popup>
  </s-card>
</template>

<style scoped>
.button_enabled {
  color: #000000;
  cursor: pointer;
}

.button_disabled {
  color: #bebebe;
  cursor: not-allowed;
}

.line {
  display: flex;
  align-content: center;
}

.line > :first-child {
  color: #000000;
  cursor: pointer;
  display: flex;
  align-items: center;
  flex-grow: 1;
}
</style>

<script lang="ts">
import { computed, defineComponent, onMounted, Ref, ref } from "vue";

import { http, httpDelete, HttpResponse } from "../utils";
import { Notify } from "./SNotification";
import { File } from "../models";
import Editor from "./Editor.vue";
import { useRouter } from "vue-router";

interface Response {
  files: File[];
  folders: string[];
}

export default defineComponent({
  name: "Manager",
  components: {
    editor: Editor,
  },
  setup({}, { attrs, emit }) {
    const path = ref(".");
    const editorPopup = ref(false);
    const imagePopup = ref(false);
    const listing: Ref<Response> = ref({
      files: new Array<File>(),
      folders: new Array<string>(),
    });
    const filepath: Ref<string> = ref("");
    const updateInProgress = ref(false);
    const editor = ref(null);
    const router = useRouter();

    const updatePath = async () => {
      let response: HttpResponse<Response>;
      try {
        response = await http<Response>(`/api/files?path=${path.value}`);
        listing.value = response.parsedBody as Response;
      } catch (err) {
        console.log(err);
        Notify("Failed to fetch files!");
        router.push({ name: "Login", params: { error: "true" }});
      }
    };

    onMounted(async () => {
      updatePath();
    });

    const openFolder = (folder: string) => {
      if (folder == "..") {
        const idx = path.value.lastIndexOf("/");
        path.value = path.value.substring(0, idx);
      } else {
        path.value = `${path.value}/${folder}`;
      }
      updatePath();
    };

    const deleteFile = async (file: File) => {
      let response: HttpResponse<string[]>;
      try {
        response = await httpDelete<string[]>(
          `/api/files?path=${path.value}/${file.name}`
        );
        const index = listing.value.files.indexOf(file);
        if (index > -1) {
          listing.value.files.splice(index, 1);
        }
      } catch (err) {
        Notify(`Failed to delete file! ${err}`);
        router.push('/login');
      }
    };

    const openFile = (file: File) => {
      filepath.value = file.path + "/" + file.name;

      if (file.content_type == "image/jpeg") {
        imagePopup.value = true;
        return;
      }

      editorPopup.value = true;
    };

    const scheduleUpdate = async () => {
      let response: HttpResponse<string[]>;
      try {
        response = await http<string[]>("/api/update");
        Notify(
          "Server update requested. Wait a few minutes for update to complete!"
        );
        updateInProgress.value = true;
      } catch (err) {
        Notify(`Failed to update server! ${err}`);
        router.push('/login');
      }
    };

    const cancelUpdate = async () => {
      let response: HttpResponse<string[]>;
      try {
        response = await httpDelete<string[]>("/api/update");
        Notify("Update canceled");
        updateInProgress.value = false;
      } catch (err) {
        Notify(`Failed to cancel update! ${err}`);
        router.push('/login');
      }
    };

    const backup = async () => {
      let response: HttpResponse<string[]>;
      try {
        response = await http<string[]>("/api/backup");
        Notify("Backup completed!");
      } catch (err) {
        Notify(`Failed to backup server! ${err}`);
        router.push('/login');
      }
    };

    return {
      listing,
      deleteFile,
      openFile,
      updateInProgress,
      scheduleUpdate,
      cancelUpdate,
      backup,
      openFolder,
      editorPopup,
      imagePopup,
      filepath,
      path,
      editor,
    };
  },
});
</script>
