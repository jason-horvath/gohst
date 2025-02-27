/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_MANIFEST_PATH: string;
  readonly VITE_PORT: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
