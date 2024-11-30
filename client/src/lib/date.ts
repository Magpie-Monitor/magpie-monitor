export const defaultDateFromUnixTimestamp = (timestamp: number): string =>
  new Date(timestamp * 1000).toLocaleString();

export const dateTimeFromTimestampMs = (timestamp: number): string => {
  return new Date(timestamp).toLocaleString();
};

export const dateFromTimestampMs = (timestamp: number): string => {
  return new Date(timestamp).toLocaleDateString();
};

export const dateOnlyFromTimestampMs = (timestamp: number): string => {
  return new Date(timestamp).toLocaleDateString();
};

export const getFirstAndLastDateFromTimestamps = (
  timestamps: number[],
): number[] => {
  const startDate = Math.min(...timestamps);
  const endDate = Math.max(...timestamps);

  return [startDate, endDate];
};

export const getDateFromTimestamps = (startDate: number): string => {
  return new Date(startDate).toISOString().split('T')[0];
};
