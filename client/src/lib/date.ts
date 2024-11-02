export const defaultDateFromUnixTimestamp = (timestamp: number): string =>
  new Date(timestamp * 1000).toLocaleString();

export const dateFromTimestampMs = (timestamp: number): string => {
  return new Date(timestamp).toLocaleString();
};

export const getFirstAndLastDateFromTimestamps = (
  timestamps: number[],
): number[] => {
  const startDate = Math.min(...timestamps);
  const endDate = Math.max(...timestamps);

  return [startDate, endDate];
};
