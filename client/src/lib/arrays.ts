export const groupBy = <T>(
  list: T[],
  getKey: (arg: T) => string,
): Record<string, T[]> => {
  const groupedEntries: Record<string, T[]> = {};

  for (const element of list) {
    const entry = groupedEntries[getKey(element)];
    if (Array.isArray(entry)) {
      entry.push(element);
    } else {
      groupedEntries[getKey(element)] = [element];
    }
  }

  return groupedEntries;
};
