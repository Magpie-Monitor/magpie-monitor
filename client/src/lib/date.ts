export const defaultDateFromUnixTimestamp = (timestamp: number): string =>
  new Date(timestamp * 1000).toLocaleString();
