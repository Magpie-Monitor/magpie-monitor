import { debounce } from 'lib/debounce';
import { useEffect } from 'react';

interface UseInfiniteScrollParams {
  scrollTargetRef: React.RefObject<HTMLDivElement>;
  debounceTreshhold?: number;
  scrollTreshhold?: number;
  handleScroll: () => void;
}

const useInfiniteScroll = ({
  scrollTargetRef,
  debounceTreshhold = 200,
  scrollTreshhold = 2,
  handleScroll,
}: UseInfiniteScrollParams) => {
  useEffect(() => {
    const onScroll = () => {
      const element = scrollTargetRef.current;
      if (!element) {
        return;
      }

      if (
        element.scrollTop + element.clientHeight >=
        element.scrollHeight / scrollTreshhold
      ) {
        handleScroll();
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
