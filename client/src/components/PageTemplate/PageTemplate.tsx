import './PageTemplate.scss';

interface PageTemplateProps {
  header: React.ReactNode;
  children: React.ReactNode;
}

const PageTemplate = ({ header, children }: PageTemplateProps) => {
  return (
    <div className={'page-template'}>
      <div className={'page-template__content'}>
        <div className={'page-template__header'}>{header}</div>
        <div className={'page-template__body'}>{children}</div>
      </div>
    </div>
  );
};

export default PageTemplate;
