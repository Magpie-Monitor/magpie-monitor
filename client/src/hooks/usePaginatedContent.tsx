import { useState } from 'react';

interface PaginatedContent<T> {
  addContent: (content: T[]) => void;
  totalContentCount: number;
  content: T[];
  contentPage: number;
  setTotalContentCount: (count: number) => void;
  isAllContentFetched: () => boolean;
}

const usePaginatedContent = <T,>(): PaginatedContent<T> => {
  const [contentPage, setContentPage] = useState(0);
  const [content, setContent] = useState<T[]>([]);
  const [totalContentCount, setTotalContentCount] = useState(-1);

  const addContent = (newContent: T[]) => {
    setContent((prev) => [...prev, ...newContent]);
    setContentPage((page) => page + 1);
  };

  const isAllContentFetched = () => {
    return content.length >= totalContentCount && totalContentCount >= 0;
  };

  return {
    addContent,
    totalContentCount,
    content,
    contentPage,
    setTotalContentCount,
    isAllContentFetched,
  };
};

export default usePaginatedContent;
