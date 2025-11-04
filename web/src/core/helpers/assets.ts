import { illustrationsSet } from "@/core/helpers/config";
import { useThemeStore } from "@/stores/theme";

export const getIllustrationsPath = (illustrationName: string): string => {
  const extension = illustrationName.substring(
    illustrationName.lastIndexOf("."),
    illustrationName.length
  );
  const illustration =
    useThemeStore().mode == "dark"
      ? `${illustrationName.substring(
          0,
          illustrationName.lastIndexOf(".")
        )}-dark`
      : illustrationName.substring(0, illustrationName.lastIndexOf("."));
  return (
    new URL(`/media/illustrations/${illustrationsSet.value}/${illustration}${extension}`, import.meta.url).href
  );
};

export const getAssetPath = (path: string): string => {
  // 参考自： https://segmentfault.com/a/1190000041676437?utm_source=sf-similar-article 解决url问题
  // console.log(new URL(path, import.meta.url).href);
  if(path.indexOf("media") != -1){
    return new URL(path, import.meta.url).href;
  }
  return import.meta.env.BASE_URL + path;
};

