export function formatCount(value: number) {
  if (value >= 10000) {
    return `${(value / 10000).toFixed(1)}w`;
  }

  return String(value);
}

export function formatDate(value?: string) {
  if (!value) {
    return '未发布';
  }

  return value.slice(0, 10);
}

export function resolveMediaUrl(url?: string) {
  if (!url) {
    return '';
  }

  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }

  return url;
}
