import { NEllipsis, type SelectGroupOption, type SelectOption } from "naive-ui";
import type { RenderLabel } from "naive-ui/es/_internal/select-menu/src/interface";
import { h, type VNodeChild } from "vue";

export const renderLabel: RenderLabel = (option: SelectOption | SelectGroupOption): VNodeChild => {
  const rawLabel = option.label ?? option.value ?? '';
  const label = typeof rawLabel === 'string' ? rawLabel : String(option.value ?? '');

  return h(
    NEllipsis,
    {
      tooltip: true,
    },
    {
      default: () => label,
    },
  );
};