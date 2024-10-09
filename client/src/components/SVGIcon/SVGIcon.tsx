import './SVGIcon.scss';

interface SVGIconProps {
  iconName: string;
}

const SVGIcon = ({ iconName }: SVGIconProps) => {
  return (
    <div className={iconName}>
      <div className={'svg-icon'} />
    </div>
  );
};

export default SVGIcon;
