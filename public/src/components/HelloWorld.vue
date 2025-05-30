<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { useAppStore } from '../stores/app'

  const store = useAppStore()
  const scriptSigCode = ref('')
  const scriptPubKeyCode = ref('')
  const scriptSigBin = ref('')
  const scriptPubKeyBin = ref('')
  const signedData = ref('315c73131650bbeec49493f3c8bd8adfbb0621a8fb83a29844973a43cbb6c063')

  const lineNumbersSigCode = ref('1')
  const lineNumbersPubKeyCode = ref('1')
  const textareaSigCode = ref<HTMLTextAreaElement | null>(null)
  const textareaPubKeyCode = ref<HTMLTextAreaElement | null>(null)
  const preSigCode = ref<HTMLElement | null>(null)
  const prePubKeyCode = ref<HTMLElement | null>(null)

  function updateLineNumbers (code: string, lineNumbers: Ref<string>) {
    const lines = code.split('\n').length
    lineNumbers.value = Array.from({ length: lines }, (_, i) => i + 1).join('\n')
  }

  function syncScrollPubKey () {
    // Синхронизируем прокрутку lineNumbers с textarea
    if (prePubKeyCode.value && textareaPubKeyCode.value) {
      prePubKeyCode.value.scrollTop = textareaPubKeyCode.value.scrollTop
    }
  }

  function syncScrollSig () {
    // Синхронизируем прокрутку lineNumbers с textarea
    if (preSigCode.value && textareaSigCode.value) {
      preSigCode.value.scrollTop = textareaSigCode.value.scrollTop
    }
  }

  const runScript = async () => {
    if (!scriptSigBin.value && !scriptPubKeyBin.value) {
      store.error = 'Code is empty'
    }
    try {
      await store.runScript({
        scriptSig: scriptSigBin.value,
        scriptPubKey: scriptPubKeyBin.value,
        signedData: signedData.value,
      })
      if (store.scriptResult) {
        //
      }
    } catch {
      console.error(store.error)
    }
  }

  const compile = async () => {
    if (!scriptSigCode.value && !scriptPubKeyCode.value) {
      store.error = 'Code is empty'
    }
    try {
      await store.compileScript({
        scriptSig: scriptSigCode.value,
        scriptPubKey: scriptPubKeyCode.value,
      })
      if (store.compileResult) {
        scriptSigBin.value = store.compileResult.scriptSig
        scriptPubKeyBin.value = store.compileResult.scriptPubKey
      }
    } catch {
      console.error(store.error)
    }
  }

  const parseScript = async () => {
    if (!scriptSigBin.value && !scriptPubKeyBin.value) {
      store.error = 'Code is empty'
    }
    try {
      await store.parseScript({
        scriptSig: scriptSigBin.value,
        scriptPubKey: scriptPubKeyBin.value,
      })
      if (store.parseResult) {
        scriptSigCode.value = store.parseResult.scriptSig
        scriptPubKeyCode.value = store.parseResult.scriptPubKey
      }
    } catch {
      console.error(store.error)
    }
  }

  watch(scriptSigCode, () => updateLineNumbers(scriptSigCode.value, lineNumbersSigCode))
  watch(scriptPubKeyCode, () => updateLineNumbers(scriptPubKeyCode.value, lineNumbersPubKeyCode))
</script>

<template>
  <v-row class="fill-height" style="min-height: 0;">
    <v-col class="fill-height flex-grow-0" style="min-height: 0;">
      <OpList />
    </v-col>
    <v-col class="fill-height flex-grow-1">
      <div
        class="p-2 mb-2 flex-grow-1"
        style="position: relative; display: flex; width: 100%; min-height: 0; height: 20vh"
      >
        <pre
          ref="preSigCode"
          class="hide-scrollbar"
          style="margin: 0; padding: 0.5em 0; min-width: 2em; text-align: right;
               color: #888; user-select: none; font-family: monospace; overflow: auto;
               height: 100%; scrollbar-width: none; -ms-overflow-style: none;
               border-right: 1px solid #888 !important; background: #23272e !important;"
        >{{ lineNumbersSigCode }}</pre>
        <textarea
          ref="textareaSigCode"
          v-model="scriptSigCode"
          class="border-sm pa-2"
          spellcheck="false"
          style="resize: none; width: 100%; display: block; height: 100%; min-height: 0; color: #888; font-family: monospace; border-left: none;"
          wrap="off"
          @scroll="syncScrollSig"
        />
      </div>
      <div
        class="p-2 mb-2 flex-grow-1"
        style="position: relative; display: flex; width: 100%; min-height: 0; height: 20vh"
      >
        <pre
          ref="prePubKeyCode"
          class="hide-scrollbar"
          style="margin: 0; padding: 0.5em 0; min-width: 2em; text-align: right;
               color: #888; user-select: none; font-family: monospace; overflow: auto;
               height: 100%; scrollbar-width: none; -ms-overflow-style: none;
               border-right: 1px solid #888 !important; background: #23272e !important;"
        >{{ lineNumbersPubKeyCode }}</pre>
        <textarea
          ref="textareaPubKeyCode"
          v-model="scriptPubKeyCode"
          class="border-sm pa-2"
          spellcheck="false"
          style="resize: none; width: 100%; display: block; height: 100%; min-height: 0; color: #888; font-family: monospace; border-left: none;"
          wrap="off"
          @scroll="syncScrollPubKey"
        />
      </div>
      <v-btn rounded="lg" size="x-large" @click="compile">Compile</v-btn>
      <div
        class="mb-2 flex-grow-1"
        style="width: 100%; min-height: 0; height: 30vh"
      >
        <textarea
          v-model="scriptSigBin"
          class="border-sm my-2 pa-2"
          spellcheck="false"
          style="resize: none; width: 100%; display: block; height: 49%; min-height: 0; color: #888; font-family: monospace;"
        />
        <textarea
          v-model="scriptPubKeyBin"
          class="border-sm my-2 pa-2"
          spellcheck="false"
          style="resize: none; width: 100%; display: block; height: 49%; min-height: 0; color: #888; font-family: monospace;"
        />
      </div>
      <div>
        <v-btn rounded="lg" size="x-large" @click="runScript">Run</v-btn>
        <v-btn class="ml-2" rounded="lg" size="x-large" @click="parseScript">Parse</v-btn>
        <v-alert v-if="store.error" class="mt-2" type="error">
          {{ store.error }}
        </v-alert>
      </div>
    </v-col>
  </v-row>

</template>

<style scoped>
.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
.hide-scrollbar {
  background: #23272e !important;
  border-right: 1px solid #888 !important;
}
textarea.border-sm,
textarea.my-2,
textarea.pa-2 {
  border-color: #888 !important;
  color: #888;
}
textarea.border-sm:focus,
textarea.my-2:focus,
textarea.pa-2:focus {
  border-color: #333 !important;
  box-shadow: none !important;
}
</style>
