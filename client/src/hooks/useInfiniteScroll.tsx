import { debounce } from 'lib/debounce';
import { useEffect, useState } from 'react';

interface UseInfiniteScrollParams {
  scrollTargetRef: React.RefObject<HTMLDivElement>;
  debounceTreshhold?: number;
  scrollTreshhold?: number;
  handleScroll: () => Promise<void> | void;
}

const useInfiniteScroll = ({
  scrollTargetRef,
  debounceTreshhold = 200,
  scrollTreshhold = 2,
  handleScroll,
}: UseInfiniteScrollParams) => {
  const [isHandling, setIsHandling] = useState(false);
  useEffect(() => {
    const onScroll = async () => {
      const element = scrollTargetRef.current;
      if (!element) {
        return;
      }

      if (isHandling === true) {
        return;
      }

      if (
        element.scrollTop + element.clientHeight >=
        element.scrollHeight / scrollTreshhold
      ) {
        setIsHandling(true);
        await handleScroll();
        setIsHandling(false);
      }
    };

    const debouncedScroll = debounce(onScroll, debounceTreshhold);
    const target = scrollTargetRef.current;

    if (target == null) {
      return;
    }

    target.addEventListener('scroll', debouncedScroll);

    return () => target.removeEventListener('scroll', debouncedScroll);
  }, [debounceTreshhold, handleScroll, scrollTargetRef, scrollTreshhold]);
};

export default useInfiniteScroll;
