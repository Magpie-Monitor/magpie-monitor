import './PageTemplate.scss';
import React from 'react';

interface PageTemplateProps {
  header: React.ReactNode;
  children: React.ReactNode;
  scrollRef?: React.RefObject<HTMLDivElement>;
  id?: string;
}

const PageTemplate = ({
  header,
  children,
  scrollRef: ref,
  id,
}: PageTemplateProps) => {
  return (
    <div className={'page-template'} ref={ref} id={id}>
      <div className={'page-template__content'}>
        <div className={'page-template__header'}>{header}</div>
        <div className={'page-template__body'}>{children}</div>
      </div>
    </div>
  );
};

export default PageTemplate;
