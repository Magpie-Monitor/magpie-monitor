export const defaultDateFromUnixTimestamp = (timestamp: number): string =>
  new Date(timestamp * 1000).toLocaleString();

export const dateTimeFromTimestampMs = (timestamp: number): string => {
  return new Date(timestamp).toLocaleString();
};

export const dateFromTimestampMs = (timestamp: number): string => {
  const date = new Date(timestamp);

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();

  return `${day}.${month}.${year}`;
};

export const dateTimeWithoutSecondsFromTimestampMs = (timestamp: number): string => {
  const date = new Date(timestamp);

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();

  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  return `${day}.${month}.${year} ${hours}:${minutes}`;
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
