//frontend/src/i18n/get-dictionary.ts
import "server-only";

const dictionaries = {
  en: () => import("./dictionaries/en.json").then((module) => module.default),
  pt: () => import("./dictionaries/pt.json").then((module) => module.default),
  fr: () => import("./dictionaries/fr.json").then((module) => module.default),
};

export type Dictionary = Awaited<ReturnType<typeof dictionaries.en>>;

export const getDictionary = async (locale: string) => {
  const loadDictionary =
    dictionaries[locale as keyof typeof dictionaries] || dictionaries.en;
  return loadDictionary();
};
