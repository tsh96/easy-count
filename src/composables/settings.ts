import { useStorage } from "@vueuse/core";

export const settings = useStorage('settings', {
  geminiApiKey: '',
  geminiModel: 'gemini-1.5-flash',
});