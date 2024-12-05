export const FadeInTransition = {
  from: { opacity: 0, transform: 'translateY(20px)' },
  enter: { opacity: 1, transform: 'translateY(0)' },
  leave: { opacity: 0, transform: 'translateY(20px)' },
  config: { duration: 100 },
  trail: 100,
};
