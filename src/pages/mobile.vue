<script setup lang="ts">
import { Icon } from '@iconify/vue';
import Peer, { type DataConnection } from 'peerjs';
import { ref } from 'vue';
import { QrcodeStream } from 'vue-qrcode-reader';
import { parsePhoto } from '../composables/parse-photo';
import { settings } from '../composables/settings';

const peer = new Peer()
const peerId = ref<string>();

const showSettings = ref(false);

peer.on('open', (id) => {
  peerId.value = id;
});

let connection = ref<DataConnection>();
const connectedPeerId = ref<string>();

function onDetect(data: any) {
  showQrCodeScanner.value = false;
  const otherPeerId = data[0]['rawValue'] as string
  connectedPeerId.value = otherPeerId;
  connectToPeer();
}

function connectToPeer() {
  if (!connectedPeerId.value) return;
  connection.value = peer.connect(connectedPeerId.value);

  connection.value.on('open', () => {
    connection.value?.send('connected');
  });
}

const processingPhoto = ref(false);

function takePhoto() {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'image/*';
  // from rear camera
  input.capture = 'environment';

  input.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;
    processingPhoto.value = true;

    const reader = new FileReader();
    reader.onload = async (e) => {

      const image = new Image();
      image.src = e.target?.result as string;
      image.onload = async () => {
        try {
          const { result } = await parsePhoto((reader.result as string)?.split(',')[1]);
          connection.value?.send(result);
        } catch (e) {
          console.error(e);
        }

        processingPhoto.value = false;
      };
    };
    reader.readAsDataURL(file);
  };

  input.click();
}

const showQrCodeScanner = ref(false);


</script>

<template lang="pug">
.w-screen.h-screen.flex.justify-center.items-center.flex-col
  .fixed.top-0.right-0.m-4
    n-button(type="default" @click="showSettings = true")
      Icon(icon="mdi:cog")

  .flex.flex-col.items-center.space-y-8(v-if="!connection")
    Icon(icon="streamline:qr-code" width="100" height="100")
    n-button.flex.items-center(type="success" size="large" @click="showQrCodeScanner = true" :loading="!peerId")
      Icon(icon="mdi:qrcode-scan")
      .mx-2 Scan QR Code
    n-input(v-model:value="connectedPeerId" placeholder="Your Peer ID" @keyup.enter="connectToPeer()")

  .flex.flex-col.items-center.space-y-8(v-else)
    Icon(icon="mdi:camera" width="100" height="100")
    n-button(type="success" size="large" @click="takePhoto()" :loading="processingPhoto")
      Icon(icon="mdi:camera")
      .mx-2 Take Photo

n-modal(v-model:show="showQrCodeScanner" preset="dialog")
  .w-30.h-30
    qrcode-stream(@detect="onDetect")

n-modal(v-model:show="showSettings" preset="dialog" title="Settings")
  n-form
    n-form-item(label="Gemini API Key")
      n-input(v-model:value="settings.geminiApiKey" placeholder="Gemini API Key")
    n-form-item(label="Gemini Model")
      n-input(v-model:value="settings.geminiModel" placeholder="Gemini API Model")
</template>